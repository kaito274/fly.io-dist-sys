BIN := $(HOME)/go/bin/maelstrom-unique-ids

build:
	go build -o $(BIN) .

test: build
	cd ../maelstrom && ./maelstrom test -w unique-ids \
		--bin $(BIN) \
		--time-limit 30 \
		--rate 1000 \
		--node-count 3 \
		--availability total \
		--nemesis partition
