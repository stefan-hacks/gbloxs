.PHONY: build run clean install test help

# Build the application
build:
	@echo "Building gbloxs..."
	@go build -o gbloxs main.go
	@echo "Build complete! Run ./gbloxs to start"

# Run the application
run:
	@go run main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f gbloxs
	@echo "Clean complete"

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Run tests (if any)
test:
	@go test ./...

# Format code
fmt:
	@go fmt ./...

# Lint code
lint:
	@golangci-lint run || echo "golangci-lint not installed, skipping..."

# Show help
help:
	@echo "Available targets:"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Install dependencies"
	@echo "  test     - Run tests"
	@echo "  fmt      - Format code"
	@echo "  lint     - Lint code"
	@echo "  help     - Show this help message"

