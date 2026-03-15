# Challenge 3c: Fault-Tolerant Broadcast

https://fly.io/dist-sys/3c/

## What this does

Same as 3b (multi-node broadcast with gossip), but now Maelstrom introduces **network partitions** — nodes temporarily can't talk to each other. Messages must still eventually reach every node once the partition heals.

## Key concept

**Retry with timeout**: When forwarding a broadcast to a neighbor, wait for `broadcast_ok`. If no reply arrives within 2 seconds (node may be partitioned), retry in a loop until it succeeds. Once the partition heals, the retry goes through and the message propagates.

This is done in a **goroutine** per neighbor so the handler doesn't block waiting for slow/dead nodes.

## What changed from 3b

| 3b | 3c |
|----|-----|
| `n.RPC()` with empty no-op handler | `n.RPC()` inside a retry loop with 2s timeout |
| Assumes message always arrives | Retries until `broadcast_ok` is received |

## Retry loop pattern

```go
go func() {
    for {
        done := make(chan struct{})
        n.RPC(node, map[string]any{
            "type":    "broadcast",
            "message": messageValue,
        }, func(reply maelstrom.Message) error {
            close(done) // got broadcast_ok
            return nil
        })

        select {
        case <-done:
            return // success
        case <-time.After(2 * time.Second):
            // no reply → retry
        }
    }
}()
```

## How to run

```bash
# Build
go build -o ~/go/bin/maelstrom-fault-tolerance-broadcast .

# Test (run from the maelstrom/ directory)
cd ../maelstrom
./maelstrom test -w broadcast --bin ~/go/bin/maelstrom-fault-tolerance-broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition
```

## Expected result

```
Everything looks good! ヽ('ー`)ノ
```
