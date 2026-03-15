package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
	"slices"
)

func main() {
	n := maelstrom.NewNode()

	var messages []int

	var neighbors []string

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		messageRaw := body["message"]
		messageValue := int(messageRaw.(float64))

		// check if this messageValue already received
		// if not, save to local memory and start broadcast to neighbors node
		isPresent := slices.Contains(messages, messageValue)
		body = map[string]any{"type": "broadcast_ok"}

		if !isPresent {
			messages = append(messages, messageValue)
			for _, node := range neighbors {
				n.RPC(node, map[string]any{
					"type":    "broadcast",
					"message": messageValue,
				}, func(reply maelstrom.Message) error {
					return nil
				})
			}
		}

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

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
