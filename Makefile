all: server client

.PHONY: server
server: proto
	go build ./cmd/server/

.PHONY: client
client: proto
	go build ./cmd/client/

clean:
	rm ./server ./client

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./pkg/proto/service.proto
