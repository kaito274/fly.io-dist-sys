BIN := $(HOME)/go/bin/maelstrom-echo

build:
	go build -o $(BIN) .

test: build
	cd ../maelstrom && ./maelstrom test -w echo \
		--bin $(BIN) \
		--node-count 1 \
		--time-limit 10
