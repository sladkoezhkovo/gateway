protoc:
	protoc -I ./proto ./proto/*.proto  --go_out=. --go-grpc_out=.

build:
	go build -v ./cmd/gateway


DEFAULT_GOAL := build