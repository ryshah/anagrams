APP_NAME=anagrams
SERVER_DIR=cmd/server
CLIENT_DIR=cmd/client

CONFIG=config.yaml

.PHONY: help build run-server run-client test fmt lint clean deps load-test

help:
	@echo "Available commands:"
	@echo "  make build        - Build server and client binaries"
	@echo "  make run-server   - Run REST API server"
	@echo "  make run-client   - Run concurrent client"
	@echo "  make test         - Run unit tests"
	@echo "  make deps         - Download dependencies"
	@echo "  make clean        - Remove build artifacts"

deps:
	go mod tidy
	go mod download

build:
	go build -o bin/server ./$(SERVER_DIR)
	go build -o bin/client ./$(CLIENT_DIR)

run-server:
	go run ./$(SERVER_DIR) --config $(CONFIG)

run-client:
	go run ./$(CLIENT_DIR) --config $(CONFIG)

test:
	go test -v ./...

clean:
	rm -rf bin
