BIN := $(HOME)/go/bin/maelstrom-multi-node-broadcast

build:
	go build -o $(BIN) .

test: build
	cd ../maelstrom && ./maelstrom test -w broadcast \
		--bin $(BIN) \
		--node-count 5 \
		--time-limit 20 \
		--rate 10
