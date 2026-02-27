# Challenge 3a: Single-Node Broadcast

https://fly.io/dist-sys/3a/

## What this does

A single node receives `broadcast` messages containing integer values, stores them in memory, and returns all stored values when a `read` is requested.

## Key concept

In-memory storage is fine here — Maelstrom never kills node processes during a test, so data persists for the duration of the run.

## Message flow

### broadcast

```
Client → Node:  { "type": "broadcast", "message": 1000 }
Node   → Client: { "type": "broadcast_ok" }
```

### read

```
Client → Node:  { "type": "read" }
Node   → Client: { "type": "read_ok", "messages": [1, 8, 72, 25] }
```

### topology

```
Client → Node:  { "type": "topology", "topology": { "n1": ["n2", "n3"] } }
Node   → Client: { "type": "topology_ok" }
```

> Topology can be ignored for 3a (single node). Use `n.NodeIDs()` if you need to build your own later.

## Go library methods

| Method                       | Use case                                                          |
| ---------------------------- | ----------------------------------------------------------------- |
| `n.Reply(msg, body)`         | Respond to an incoming request (adds `in_reply_to` automatically) |
| `n.Send(dest, body)`         | Fire-and-forget message to another node (no response expected)    |
| `n.RPC(dest, body, handler)` | Send to another node and handle its response                      |

For 3a, only `Reply()` is needed. `Send()` and `RPC()` become important in 3b (multi-node).

## How to run

```bash
# Build
go build -o ~/go/bin/maelstrom-broadcast .

# Test (run from the maelstrom/ directory)
cd ../maelstrom
./maelstrom test -w broadcast --bin ~/go/bin/maelstrom-broadcast --node-count 1 --time-limit 20 --rate 10
```

## Expected result

```
Everything looks good! ヽ('ー`)ノ
```
