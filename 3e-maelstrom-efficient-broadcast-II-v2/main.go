package main

import (
	"encoding/json"
	"log"
	"sort"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	// Use a set instead of []int for faster membership checks.
	messages := make(map[int]struct{})
	var mu sync.Mutex

	const SYNC_INTERVAL = 500 * time.Millisecond

	// ----------------------------
	// Helpers
	// ----------------------------

	getAllMessages := func() []int {
		mu.Lock()
		defer mu.Unlock()

		result := make([]int, 0, len(messages))
		for m := range messages {
			result = append(result, m)
		}
		sort.Ints(result) // helpful for deterministic reads/debugging
		return result
	}

	addMessage := func(val int) bool {
		mu.Lock()
		defer mu.Unlock()

		if _, exists := messages[val]; exists {
			return false
		}
		messages[val] = struct{}{}
		return true
	}

	addMessages := func(vals []int) {
		mu.Lock()
		defer mu.Unlock()

		for _, v := range vals {
			messages[v] = struct{}{}
		}
	}

	// ----------------------------
	// Client broadcast
	// ----------------------------
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		val := int(body["message"].(float64))
		addMessage(val)

		return n.Reply(msg, map[string]any{
			"type": "broadcast_ok",
		})
	})

	// ----------------------------
	// Read current messages
	// ----------------------------
	n.Handle("read", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]any{
			"type":     "read_ok",
			"messages": getAllMessages(),
		})
	})

	// ----------------------------
	// Topology
	// Keep handler for protocol compatibility.
	// This anti-entropy version doesn't rely on the provided topology;
	// it can reconcile directly with all peers.
	// ----------------------------
	n.Handle("topology", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]any{
			"type": "topology_ok",
		})
	})

	// ----------------------------
	// sync
	// A peer sends us the set of message IDs it currently has.
	// We compute which messages that peer is missing,
	// then send those back in a separate sync_ok message.
	// ----------------------------
	n.Handle("sync", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		rawHave, ok := body["have"].([]any)
		if !ok {
			rawHave = []any{}
		}

		peerHave := make(map[int]struct{}, len(rawHave))
		for _, v := range rawHave {
			peerHave[int(v.(float64))] = struct{}{}
		}

		// Compute what the sender is missing.
		missing := make([]int, 0)

		mu.Lock()
		for m := range messages {
			if _, ok := peerHave[m]; !ok {
				missing = append(missing, m)
			}
		}
		mu.Unlock()

		// Send only the missing messages back to the sender.
		// This is the "repair" step of anti-entropy.
		if len(missing) > 0 {
			if err := n.Send(msg.Src, map[string]any{
				"type":     "sync_ok",
				"messages": missing,
			}); err != nil {
				return err
			}
		}

		return nil
	})

	// ----------------------------
	// sync_ok
	// Peer sends us the messages we were missing.
	// ----------------------------
	n.Handle("sync_ok", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		rawMsgs, ok := body["messages"].([]any)
		if !ok {
			return nil
		}

		vals := make([]int, 0, len(rawMsgs))
		for _, v := range rawMsgs {
			vals = append(vals, int(v.(float64)))
		}

		addMessages(vals)
		return nil
	})

	// ----------------------------
	// Periodic anti-entropy loop
	// Periodically send our summary ("have") to all peers.
	// They compare against their own state and send us back
	// whatever we are missing.
	// ----------------------------
	go func() {
		ticker := time.NewTicker(SYNC_INTERVAL)
		defer ticker.Stop()

		for range ticker.C {
			have := getAllMessages()

			for _, peer := range n.NodeIDs() {
				if peer == n.ID() {
					continue
				}

				if err := n.Send(peer, map[string]any{
					"type": "sync",
					"have": have,
				}); err != nil {
					log.Printf("failed to send sync to %s: %v", peer, err)
				}
			}
		}
	}()

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}