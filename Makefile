.PHONY: build test clean run-tests fmt vet lint install

# Build variables
BINARY_NAME=devpod-provider-upcloud
BUILD_DIR=bin
GO=go
GOFLAGS=-v

# Default target
all: build

# Build the provider binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Run tests
test:
	@echo "Running unit tests..."
	$(GO) test -v ./...

# Run BDD tests with Godog
bdd:
	@echo "Running BDD tests..."
	$(GO) test -v -run TestFeatures

# Run all tests (unit + BDD)
test-all: test bdd

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed" && exit 1)
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Install the provider locally for testing
install: build
	@echo "Installing provider locally..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) ~/.local/bin/ 2>/dev/null || \
		(mkdir -p ~/.local/bin && cp $(BUILD_DIR)/$(BINARY_NAME) ~/.local/bin/)
	@echo "Provider installed to ~/.local/bin/$(BINARY_NAME)"

# Run the provider with arguments
run: build
	@$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Development mode - rebuild and run on file changes (requires entr)
watch:
	@which entr > /dev/null || (echo "entr not installed" && exit 1)
	find . -name '*.go' | entr -r make build

# Generate test coverage
coverage:
	@echo "Generating test coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

# Check for security issues (requires gosec)
security:
	@which gosec > /dev/null || (echo "gosec not installed" && exit 1)
	gosec ./...

# Help target
help:
	@echo "Available targets:"
	@echo "  build      - Build the provider binary"
	@echo "  test       - Run unit tests"
	@echo "  bdd        - Run BDD tests with Godog"
	@echo "  test-all   - Run all tests (unit + BDD)"
	@echo "  fmt        - Format Go code"
	@echo "  vet        - Run go vet"
	@echo "  lint       - Run linter (requires golangci-lint)"
	@echo "  clean      - Remove build artifacts"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  install    - Install provider locally"
	@echo "  run        - Build and run with ARGS"
	@echo "  watch      - Auto-rebuild on file changes (requires entr)"
	@echo "  coverage   - Generate test coverage report"
	@echo "  security   - Check for security issues (requires gosec)"
	@echo "  help       - Show this help message"