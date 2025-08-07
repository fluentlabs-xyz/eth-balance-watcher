.PHONY: build run test clean docker docker-up docker-down deps fmt vet

# Build the application
build:
	go build -o bin/eth-balance-watcher .

# Run the application
run:
	go run .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Build Docker image
docker:
	docker build -t eth-balance-watcher .

# Start services with docker-compose
docker-up:
	docker-compose up -d

# Stop services with docker-compose
docker-down:
	docker-compose down

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Development setup
dev-setup: deps fmt vet

# Full check before commit
check: fmt vet test

# Show help
help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  docker     - Build Docker image"
	@echo "  docker-up  - Start with docker-compose"
	@echo "  docker-down- Stop docker-compose services"
	@echo "  deps       - Download dependencies"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  lint       - Run linter"
	@echo "  check      - Full pre-commit check"