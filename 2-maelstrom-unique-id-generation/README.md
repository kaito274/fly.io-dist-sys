# Challenge 2: Unique ID Generation

https://fly.io/dist-sys/2/

## What this does

A 3-node cluster receives `generate` requests and must return a globally unique ID each time — across all nodes, even during network partitions.

## Key concept

In a distributed system with network partitions, nodes can't always talk to each other — so any ID scheme that requires coordination (e.g. a shared counter) will fail.

**UUID v4** solves this with pure randomness: a 128-bit random value where the collision probability is astronomically low (~1 in 10³⁸). Each node generates IDs independently with no coordination needed.

## Message flow

```
Client -> Node:  { "type": "generate" }
Node -> Client:  { "type": "generate_ok", "id": "550e8400-e29b-41d4-a716-446655440000" }
```

## How to run

```bash
# Build
go build -o ~/go/bin/maelstrom-unique-ids .

# Test (run from the maelstrom/ directory)
cd ../maelstrom
./maelstrom test -w unique-ids --bin ~/go/bin/maelstrom-unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
```

Notable flags:

- `--node-count 3` — 3 nodes must all generate unique IDs independently
- `--rate 1000` — 1000 requests/second
- `--nemesis partition` — Maelstrom will randomly cut network links between nodes during the test
- `--availability total` — every request must succeed (no timeouts allowed)

## Expected result

```
Everything looks good! ヽ('ー`)ノ
```
