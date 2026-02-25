# Challenge 1: Echo

https://fly.io/dist-sys/1/

## What this does

A single node receives an `echo` message and replies with `echo_ok`, returning the same body.

This is the "hello world" of Maelstrom — just to verify the node can communicate over stdin/stdout.

## Key concept

Maelstrom talks to your node process via **stdin/stdout as a JSON pipe**:

```
Maelstrom ---(stdin)---> your binary ---(stdout)---> Maelstrom
```

- Every message is a JSON object on its own line
- **Never write anything other than valid JSON to stdout** — it will corrupt the protocol
- Use `log.Printf(...)` for debug output (goes to stderr, safe)

## How to run

```bash
# Build
go build -o ~/go/bin/maelstrom-echo .

# Test (run from the maelstrom/ directory)
cd ../maelstrom
./maelstrom test -w echo --bin ~/go/bin/maelstrom-echo --node-count 1 --time-limit 10
```

## Expected result

```
Everything looks good! ヽ('ー`)ノ
```
