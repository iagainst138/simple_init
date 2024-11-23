SHELL=/bin/bash

CONFIG=./sample_configs/config.yaml

build:
	go build -o init ./cmd/sinit/

run: build
	./init -config ${CONFIG}
