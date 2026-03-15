package main

import (
	"encoding/json"
	"log"
	"slices"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	var messages []int

	var pendingMessages []int
	var mu sync.Mutex

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		messageRaw := body["message"]
		messageValue := int(messageRaw.(float64))

		mu.Lock()
		isPresent := slices.Contains(messages, messageValue)

		if !isPresent {
			messages = append(messages, messageValue)
			pendingMessages = append(pendingMessages, messageValue)
		}
		mu.Unlock()

		body = map[string]any{"type": "broadcast_ok"}
		return n.Reply(msg, body)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "read_ok"
		body["messages"] = messages

		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "topology_ok"
		delete(body, "topology")

		return n.Reply(msg, body)
	})

	n.Handle("gossip", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		messageBatchRaw, ok := body["message-batch"].([]any)
		if !ok || messageBatchRaw == nil {
			return nil
		}

		mu.Lock()
		for _, v := range messageBatchRaw {
			val := int(v.(float64))
			if !slices.Contains(messages, val) {
				messages = append(messages, val)
				pendingMessages = append(pendingMessages, val)
			}
		}
		mu.Unlock()
		return nil
	})

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		for range ticker.C {
			mu.Lock()
			batch := pendingMessages
			pendingMessages = nil
			mu.Unlock()
			if len(batch) == 0 {
				continue
			}
			for _, peer := range n.NodeIDs() {
				if peer == n.ID() {
					continue // skip yourself
				}
				n.Send(peer, map[string]any{
					"type":          "gossip",
					"message-batch": batch,
				})
			}

		}
	}()

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
