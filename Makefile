.PHONY: help build run test clean generate-sdk install-deps

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build the MCP server
build: ## Build the MCP server binary
	@echo "Building MCP server..."
	go build -o bin/linkwarden-mcp-server ./cmd/server

# Run the MCP server
run: ## Run the MCP server
	@echo "Running MCP server..."
	go run ./cmd/server

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean: ## Clean build artifacts and generated files
	@echo "Cleaning..."
	rm -rf bin/
	rm -f pkg/linkwarden/*.go

# Generate SDK from OpenAPI spec
generate-sdk: ## Generate Go SDK from OpenAPI specification
	@echo "Generating Linkwarden SDK..."
	./scripts/generate-sdk.sh

# Install development dependencies
install-deps: ## Install required development dependencies
	@echo "Installing oapi-codegen..."
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo "Installing other dependencies..."
	go mod download

# Initialize project (run after cloning)
init: install-deps generate-sdk ## Initialize the project for development
	@echo "Project initialized successfully!"

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint: ## Lint Go code (requires golangci-lint)
	@echo "Linting code..."
	golangci-lint run

# Update dependencies
update-deps: ## Update Go module dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Development workflow
dev: clean generate-sdk build ## Clean, generate SDK, and build for development
