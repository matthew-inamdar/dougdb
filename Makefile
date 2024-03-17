.DEFAULT_GOAL := help

.PHONY: help
help: ## print help (this message)
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
    	| sed -n 's/^\(.*\): \(.*\)## \(.*\)/\1;\3/p' \
    	| column -t  -s ';'

.PHONY: build
build: ## build dougdb
	go build -o bin/dougdb ./cmd/dougdb

.PHONY: test
test: ## run tests
	@echo "Running unit tests..."
	@go test -v ./...
	@echo "Running race condition tests..."
	@go test -v -race ./...

.PHONY: protos
protos: ## compile protobuf
	protoc --go_out=./gen --go_opt=paths=source_relative \
		--go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
        --proto_path=protos raft/raft.proto
