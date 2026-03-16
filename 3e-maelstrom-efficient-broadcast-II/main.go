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

	var neighbors []string
	var pendingMessages []int
	var mu sync.Mutex
	const BATCH_SIZE_LIMIT = 60
	const WAIT_TIME = 500
	trigger := make(chan struct{}, 1)

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
			if len(pendingMessages) >= BATCH_SIZE_LIMIT {
				select {
				case trigger <- struct{}{}: // signal to flush
				default:
				}
			}
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

		// save neighbors
		topology := body["topology"].(map[string]any)
		rawNeighbors := topology[n.ID()].([]any)
		for _, v := range rawNeighbors {
			neighbors = append(neighbors, v.(string))
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
				if len(pendingMessages) >= BATCH_SIZE_LIMIT {
					select {
					case trigger <- struct{}{}: // signal to flush
					default:
					}
				}
			}
		}
		mu.Unlock()
		return nil
	})

	flush := func() {
		mu.Lock()
		batch := pendingMessages
		pendingMessages = nil
		mu.Unlock()
		if len(batch) == 0 {
			return
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

	go func() {
		ticker := time.NewTicker(WAIT_TIME * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				flush()
			case <-trigger:
				ticker.Reset(WAIT_TIME * time.Millisecond)
				flush()
			}
		}
	}()

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
