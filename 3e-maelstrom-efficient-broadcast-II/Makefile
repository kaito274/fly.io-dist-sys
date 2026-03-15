BIN := $(HOME)/go/bin/maelstrom-efficient-broadcast-II
MAELSTROM := ../maelstrom/maelstrom

build:
	go build -o $(BIN) .

test: build
	cd ../maelstrom && ./maelstrom test -w broadcast \
		--bin $(BIN) \
		--node-count 25 \
		--time-limit 20 \
		--rate 100 \
		--latency 100 \
		2>&1 | tee /tmp/maelstrom-3d.log

	@echo ""
	@echo "=== Performance Summary ==="
	@printf "Valid:       %s\n"  "$$(grep -oE ':valid\? (true|false)' /tmp/maelstrom-3d.log | tail -1)"
	@printf "Msgs/op:     %s\n"  "$$(grep -oE 'msgs-per-op [0-9.]+' /tmp/maelstrom-3d.log | head -1)"
	@printf "Latency p50: %sms\n" "$$(grep -oE '0\.5 [0-9]+' /tmp/maelstrom-3d.log | tail -5 | head -1 | awk '{print $$2}')"
	@printf "Latency p95: %sms\n" "$$(grep -oE '0\.95 [0-9]+' /tmp/maelstrom-3d.log | tail -5 | head -1 | awk '{print $$2}')"
	@printf "Latency p99: %sms\n" "$$(grep -oE '0\.99 [0-9]+' /tmp/maelstrom-3d.log | tail -5 | head -1 | awk '{print $$2}')"
	@printf "Latency max: %sms\n" "$$(grep -oE '[[:space:]]1 [0-9]+\},' /tmp/maelstrom-3d.log | tail -1 | grep -oE '[0-9]+' | tail -1)"
