# Challenge 3d: Efficient Broadcast

https://fly.io/dist-sys/3d/

## What this does

Same gossip broadcast as 3b/3c, but optimized to meet strict **latency and message-count targets** with 25 nodes and 100ms artificial network latency.

## Key concepts

**Batching**: Instead of sending each message immediately, accumulate new messages and send them all at once. This drastically reduces `msgs-per-op`.

**All-to-all gossip**: Instead of following the Maelstrom-provided grid topology (which takes 8 hops to span 25 nodes), each node gossips directly to all other nodes via `n.NodeIDs()`. This ensures every message arrives in exactly 1 hop, keeping latency ≤ 400ms even with 100ms network delay.

**Dual-trigger flush**: A batch is sent when either condition is met — whichever comes first:
- ⏱️ **Timeout** (`WAIT_TIME = 300ms`) — flush whatever is pending
- 📦 **Batch size** (`BATCH_SIZE_LIMIT = 30`) — flush immediately if the batch is full

This is implemented with a `select` listening on two channels: a `time.Ticker` and a buffered `trigger` channel.

**Mutex protection**: The `messages` and `pendingMessages` slices are shared between the message handler goroutines and the ticker goroutine. A `sync.Mutex` protects all read-check-write operations atomically.

## What changed from 3c

| 3c | 3d |
|----|-----|
| Immediate RPC per message with retry | Batched gossip via ticker + batch-size trigger |
| Gossips to topology neighbors | Gossips to all nodes via `n.NodeIDs()` |
| No mutex needed | Mutex required (ticker + handlers concurrent) |

## Performance targets

| Metric | Target | Achieved |
|---|---|---|
| Median latency | ≤ 400ms | ~219ms  |
| Max latency | ≤ 600ms | ~391  |
| Messages-per-op | ≤ 30 | ~24-25  |

## How to run

```bash
# Build + test + show performance summary
make test

# Or just build
make build
```

After the test, `make test` automatically prints:
```
=== Performance Summary ===
Valid:       :valid? true
Msgs/op:     msgs-per-op 24.260963
Latency p50: 213ms
Latency p95: 362ms
Latency p99: 384ms
Latency max: 391ms
```

## Expected result

```
Everything looks good! ヽ('ー`)ノ
```
