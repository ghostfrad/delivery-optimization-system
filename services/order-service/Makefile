.PHONY: help build run test deps clean

help:
	@echo "Available commands:"
	@echo "  make build        Build the service"
	@echo "  make run          Run the service (uses mocks)"
	@echo "  make test         Run unit tests"
	@echo "  make deps         Download dependencies"
	@echo "  make clean        Clean build artifacts"

build:
	go build -o bin/order-service cmd/main.go

run:
	go run cmd/main.go

test:
	go test ./internal/service/... -v

test-coverage:
	go test ./internal/service/... -coverprofile=coverage.out
	go tool cover -html=coverage.out

deps:
	go mod download
	go mod tidy

clean:
	rm -rf bin/