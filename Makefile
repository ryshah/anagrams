APP_NAME=anagrams
SERVER_DIR=cmd/server
CLIENT_DIR=cmd/client

CONFIG=config.yaml

.PHONY: help build run-server run-client test fmt lint clean deps load-test

help:
	@echo "Available commands:"
	@echo "  make build        - Build server and client binaries (including docker)"
	@echo "  make run-server   - Run REST API server"
	@echo "  make run-client   - Run concurrent client"
	@echo "  make test         - Run unit tests"
	@echo "  make deps         - Download dependencies"
	@echo "  make clean        - Remove build artifacts (including docker)"
	@echo "  make docker-build - Build the Docker server"
	@echo "  make docker-run   - Run server in Docker"
	@echo "  make docker-clean - Stop and removes the running container"

deps:
	go mod tidy
	go mod download

build: docker-build
	go build -o bin/server ./$(SERVER_DIR)
	go build -o bin/client ./$(CLIENT_DIR)

run-server:
	go run ./$(SERVER_DIR) --config $(CONFIG)

run-client:
	go run ./$(CLIENT_DIR) --config $(CONFIG)

test:
	go test -v ./...

clean: docker-clean
	rm -rf bin

docker-build:
	docker build -t anagrams .

docker-run:
    # Adding network=host for Auth middleware to allow requests from same machine
	docker run -p 127.0.0.1:8080:8080  --network=host --name anagram anagrams

docker-clean:
	docker stop anagram || true
	docker rm -f anagram
