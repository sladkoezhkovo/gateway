build:
	go build -v ./cmd/gateway

protoc:
	protoc -I ./proto ./proto/*.proto  --go_out=. --go-grpc_out=.