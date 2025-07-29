.PHONY: build run test clean dev docker-build docker-run

# Variables
BINARY_NAME=server
BUILD_DIR=./bin
MAIN_PATH=./cmd/server
DOCKER_IMAGE=nft-platform-backend

# Build the application
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	@echo "Starting application..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run in development mode
dev:
	@echo "Starting in development mode..."
	go run $(MAIN_PATH)/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):latest .

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):latest

# Setup development environment
setup-dev:
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go mod tidy

# Database migrations (when implemented)
migrate-up:
	@echo "Running database migrations..."
	# Add migration commands here

migrate-down:
	@echo "Reverting database migrations..."
	# Add migration rollback commands here

# Start PostgreSQL with Docker (for development)
db-start:
	@echo "Starting PostgreSQL container..."
	docker run --name nft-platform-db -e POSTGRES_DB=nft_platform -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:13

# Stop PostgreSQL container
db-stop:
	@echo "Stopping PostgreSQL container..."
	docker stop nft-platform-db && docker rm nft-platform-db

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run in development mode"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Build and run Docker container"
	@echo "  setup-dev     - Setup development environment"
	@echo "  db-start      - Start PostgreSQL container"
	@echo "  db-stop       - Stop PostgreSQL container"
	@echo "  help          - Show this help"
