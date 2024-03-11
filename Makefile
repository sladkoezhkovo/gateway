build:
	go build -v ./cmd/gateway

debug:
	make build
	.\gateway.exe -config configs/local.yml

protoc:
	protoc -I ./proto ./proto/*.proto  --go_out=. --go-grpc_out=.