# Challenge 3e: Efficient Broadcast II

https://fly.io/dist-sys/3e/

## What this does

Same as 3d (batched all-to-all gossip), but with **tighter message-count targets** and **more relaxed latency targets**. The key is tuning the gossip ticker interval to hit the new sweet spot.

## What changed from 3d

| Metric | 3d target | 3e target |
|---|---|---|
| Messages-per-op | ≤ 30 | **≤ 20** (tighter) |
| Median latency | ≤ 400ms | ≤ 1000ms (relaxed) |
| Max latency | ≤ 600ms | ≤ 2000ms (relaxed) |

**The only code change:** ticker interval `300ms` → `500ms`.

A longer interval means more messages batch together before sending → fewer msgs-per-op. The latency increases slightly but stays well within the relaxed 1s/2s ceiling.

## The trade-off

| Ticker interval | Msgs/op | Latency |
|---|---|---|
| 🔼 Longer (500ms) | 📉 Fewer | 📈 Higher |
| 🔽 Shorter (300ms) | 📈 More | 📉 Lower |

3d and 3e are the same system — just dialing the ticker to satisfy different requirement curves.

## Performance targets

| Metric | Target | Achieved |
|---|---|---|
| Messages-per-op | ≤ 20 | ~16.5 |
| Median latency | ≤ 1000ms | ~321ms |
| Max latency | ≤ 2000ms | ~564ms |

## How to run

```bash
make test
```

After the test, `make test` automatically prints:
```
=== Performance Summary ===
Valid:       :valid? true
Msgs/op:     msgs-per-op 16.51
Latency p50: 321ms
Latency p95: 493ms
Latency p99: 500ms
Latency max: 564ms
```

## Expected result

```
Everything looks good! ヽ('ー`)ノ
```
