SHELL=/bin/bash

CONFIG=./sample_configs/config.json

build:
	go build -o init ./cmd/sinit/

run: build
	./init -config ${CONFIG}
