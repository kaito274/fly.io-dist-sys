# Challenge 3b: Multi-Node Broadcast

https://fly.io/dist-sys/3b/

## What this does

5 nodes receive `broadcast` messages and must propagate them to every other node via gossip, so that any node can return the full set of messages on a `read` request.

## Key concepts

**Gossip / flooding**: When a node receives a new message, it forwards it to all its neighbors. Each node deduplicates by checking if it already has the value — if so, it stops forwarding.

**Topology**: Maelstrom sends a `topology` message once at startup telling each node who its direct neighbors are. This is saved and used when forwarding broadcasts.

**`n.RPC()` vs `n.Send()`**: When forwarding to neighbor nodes, use `RPC()` not `Send()`. Neighbor nodes always reply with `broadcast_ok`, and without a registered reply handler, the node crashes.

## Message flow

### broadcast (client → node, new message)

```
Client → n1:  { "type": "broadcast", "message": 42 }
n1     → n2:  { "type": "broadcast", "message": 42 }  ← gossip to neighbors
n1     → n0:  { "type": "broadcast", "message": 42 }  ← gossip to neighbors
n1     → Client: { "type": "broadcast_ok" }
```

### broadcast (node → node, already seen)

```
n2 → n1: { "type": "broadcast", "message": 42 }
n1 sees 42 is already stored → does NOT forward again (dedup)
n1 → n2: { "type": "broadcast_ok" }
```

### read

```
Client → Node:  { "type": "read" }
Node   → Client: { "type": "read_ok", "messages": [1, 8, 42] }
```

### topology

```
Maelstrom → Node: { "type": "topology", "topology": { "n0": ["n1","n3"], ... } }
Node      → Maelstrom: { "type": "topology_ok" }
```

## Go library methods used

| Method                       | Use case                                                          |
| ---------------------------- | ----------------------------------------------------------------- |
| `n.Reply(msg, body)`         | Respond to an incoming request (adds `in_reply_to` automatically) |
| `n.RPC(dest, body, handler)` | Forward broadcast to neighbor and handle `broadcast_ok` reply     |
| `n.ID()`                     | Get this node's own ID (e.g. `"n2"`)                              |

## How to run

```bash
# Build
go build -o ~/go/bin/maelstrom-multi-node-broadcast .

# Test (run from the maelstrom/ directory)
cd ../maelstrom
./maelstrom test -w broadcast --bin ~/go/bin/maelstrom-multi-node-broadcast --node-count 5 --time-limit 20 --rate 10
```

## Expected result

```
Everything looks good! ヽ('ー`)ノ
```
