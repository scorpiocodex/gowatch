.PHONY: all build install test test-race test-cover clean fmt lint run help

# Variables
BINARY_NAME=gowatch
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=./cmd/gowatch
GO=go
GOFLAGS=-v

# Default target
all: clean fmt lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GO) build $(GOFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Built: $(BINARY_PATH)"

# Install to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(GOFLAGS) $(MAIN_PATH)
	@echo "Installed to: $$(go env GOPATH)/bin/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	$(GO) test -race -v ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"
	$(GO) tool cover -func=coverage.out | grep total

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Cleaned"

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	gofmt -w .

# Run linters
lint:
	@echo "Running linters..."
	$(GO) vet ./...
	@echo "Checking formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "The following files need formatting:"; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "Lint checks passed"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_PATH)

# Initialize project
init:
	@echo "Initializing project..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Project initialized"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 $(GO) build -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GO) build -o bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GO) build -o bin/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GO) build -o bin/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Built for all platforms in bin/"

# Help target
help:
	@echo "GoWatch Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all         - Clean, format, lint, test, and build"
	@echo "  build       - Build the binary"
	@echo "  install     - Install to GOPATH/bin"
	@echo "  test        - Run tests"
	@echo "  test-race   - Run tests with race detector"
	@echo "  test-cover  - Run tests with coverage report"
	@echo "  clean       - Remove build artifacts"
	@echo "  fmt         - Format code"
	@echo "  lint        - Run linters"
	@echo "  run         - Build and run the application"
	@echo "  init        - Initialize project dependencies"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  help        - Show this help message"