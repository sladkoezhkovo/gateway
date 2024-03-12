build:
	go build -v ./cmd/gateway

debug:
	make build
	.\gateway.exe -config configs/local.yml

protoc:
	protoc ./proto/admin.proto 		--go_out=.	--go-grpc_out=.
	protoc ./proto/auth.proto 		--go_out=.	--go-grpc_out=.
	protoc ./proto/invoice.proto	--go_out=.  --go-grpc_out=.
