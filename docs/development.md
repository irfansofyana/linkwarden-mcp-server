# Development Guide

This guide covers how to set up a development environment, understand the codebase, and contribute to the linkwarden-mcp-server project.

## Prerequisites

- Go 1.23 or later
- Git
- Make
- Docker (optional, for certain development tasks)
- Linkwarden instance for testing

## Development Setup

### 1. Clone the Repository

```bash
git clone https://github.com/irfansofyana/linkwarden-mcp-server.git
cd linkwarden-mcp-server
```

### 2. Install Dependencies

```bash
make deps
```

### 3. Set Up Environment

Create a development configuration file:

```bash
cp config.example.yaml config.dev.yaml
```

Edit `config.dev.yaml` with your development Linkwarden instance details:

```yaml
base-url: http://localhost:3000
token: your-dev-token-here
toolsets: search,collection
readonly: false
log-file: ./dev.log
```

### 4. Build and Run

```bash
# Build the project
make build

# Run with development configuration
./linkwarden-mcp-server --config config.dev.yaml
```

## Project Structure

```
linkwarden-mcp-server/
├── cmd/linkwarden-mcp-server/    # Main application entry point
├── pkg/                          # Core packages
│   ├── linkwarden/              # Linkwarden API client (generated)
│   ├── linkwardenmcp/           # MCP tool implementations
│   │   ├── collection.go        # Collection management tools
│   │   ├── search.go           # Search functionality tools
│   │   ├── tools.go            # Toolset registration
│   │   └── tools_param.go      # Parameter validation helpers
│   ├── mcpgo/                  # MCP Go framework wrapper
│   ├── observability/          # Logging and monitoring
│   └── toolsets/               # Toolset management system
├── docs/                        # Documentation
├── examples/                    # Usage examples
├── scripts/                     # Development scripts
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
└── go.sum                       # Go module checksums
```

## Architecture Overview

### Core Components

1. **MCP Server** (`cmd/linkwarden-mcp-server/main.go`)
   - Initializes the MCP server
   - Handles configuration and startup
   - Manages toolset registration

2. **Toolsets** (`pkg/toolsets/`)
   - Modular system for organizing functionality
   - Supports read/write separation
   - Allows selective toolset enabling

3. **Tool Implementation** (`pkg/linkwardenmcp/`)
   - Individual tool implementations
   - Parameter validation and error handling
   - Linkwarden API integration

4. **API Client** (`pkg/linkwarden/`)
   - Auto-generated from OpenAPI specification
   - Type-safe client for Linkwarden API
   - Request/response structures

### Data Flow

```
MCP Client → MCP Server → Toolset → Tool → Linkwarden API
```

## Adding New Tools

### 1. Create Tool Implementation

Add your tool in the appropriate file in `pkg/linkwardenmcp/`:

```go
// In collection.go, search.go, or create a new file
func NewTool(
    obs *observability.Observability,
    client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
    params := []mcpgo.ToolParameter{
        mcpgo.WithString(
            "name",
            mcpgo.Description("The name parameter"),
        ),
    }

    handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
        // Validate parameters
        args := make(map[string]interface{})
        validator := NewValidator(&req)
        validator.ValidateAndAddRequiredString(args, "name")

        if result, err := validator.HandleErrorsIfAny(); result != nil {
            return result, err
        }

        // Call Linkwarden API
        resp, err := client.SomeApiCallWithResponse(ctx, args["name"].(string))
        if err != nil {
            return mcpgo.NewToolResultError("API call failed: " + err.Error()), nil
        }

        // Return result
        if resp.JSON200 != nil {
            return mcpgo.NewToolResultJSON(resp.JSON200)
        }

        return mcpgo.NewToolResultError("API returned error: " + resp.Status()), nil
    }

    return mcpgo.NewTool(
        "tool_name",
        "Tool description",
        params,
        handler,
    )
}
```

### 2. Register Tool in Toolset

Add your tool to the appropriate toolset in `pkg/linkwardenmcp/tools.go`:

```go
// Read tools
toolset.AddReadTools(
    NewTool(obs, client),
)

// Write tools
toolset.AddWriteTools(
    NewTool(obs, client),
)
```

### 3. Update Documentation

Add your tool to:
- `docs/tools-reference.md`
- `README.md` (if it's a major feature)
- Any relevant examples

## Parameter Validation

The project provides a comprehensive validation system in `pkg/linkwardenmcp/tools_param.go`:

### Validation Patterns

```go
// Required parameter
validator.ValidateAndAddRequiredString(args, "name")

// Optional parameter
validator.ValidateAndAddOptionalString(args, "description")

// Number parameter
validator.ValidateAndAddRequiredInt(args, "id")

// Boolean parameter
validator.ValidateAndAddOptionalBool(args, "pinned")

// Handle validation errors
if result, err := validator.HandleErrorsIfAny(); result != nil {
    return result, err
}
```

### Type Safety

All parameters are type-safe with automatic JSON marshaling/unmarshaling:

```go
// Extract value with type safety
value, err := extractValueGeneric[string](request, "name", true)
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Run tests with coverage
make test-coverage
```

### Writing Tests

#### Unit Tests

```go
func TestMyTool(t *testing.T) {
    obs := observability.NewObservability("test")
    client := &mockClient{}

    tool := MyTool(obs, client)

    // Test successful call
    req := mcpgo.CallToolRequest{
        Arguments: map[string]interface{}{
            "name": "test",
        },
    }

    result, err := tool.Handler(context.Background(), req)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

#### Integration Tests

```go
func TestMyToolIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Skip if no test configuration
    if os.Getenv("TEST_BASE_URL") == "" {
        t.Skip("TEST_BASE_URL not set")
    }

    obs := observability.NewObservability("test")
    client := createTestClient()

    tool := MyTool(obs, client)

    // Test with real API
    req := mcpgo.CallToolRequest{
        Arguments: map[string]interface{}{
            "id": 1,
        },
    }

    result, err := tool.Handler(context.Background(), req)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

## Code Quality

### Formatting and Linting

```bash
# Format code
make fmt

# Run linter
make lint

# Fix linting issues
make lint-fix
```

### Code Style Guidelines

1. **Go Conventions**: Follow standard Go formatting and naming conventions
2. **Error Handling**: Use explicit error handling with meaningful messages
3. **Documentation**: Add godoc comments for all exported functions
4. **Testing**: Maintain high test coverage for all tools

## SDK Generation

The Linkwarden client SDK is automatically generated from the OpenAPI specification:

```bash
# Generate SDK from OpenAPI spec
make generate-sdk

# Generate with custom options
make generate-sdk CUSTOM_OPTIONS="--additional-property=typeName"
```

### SDK Updates

When the Linkwarden API changes:

1. Update the OpenAPI specification
2. Regenerate the SDK
3. Update tool implementations to use new API features
4. Add tests for new functionality

## Build and Release

### Building for Development

```bash
# Build development binary
make build

# Build with debug symbols
make build-dev

# Build for current platform
make build-local
```

### Building for Release

```bash
# Build for multiple platforms
make build-all

# Create release packages
make package

# Build Docker image
make docker-build
```

## Debugging

### Logging

Enable debug logging:

```bash
./linkwarden-mcp-server \
  --config config.dev.yaml \
  --log-file ./debug.log
```

### Common Debugging Scenarios

#### API Issues

```bash
# Test API connectivity
curl -H "Authorization: Bearer your-token" \
     https://your-linkwarden-instance.com/api/v1/collections
```

#### Configuration Issues

```bash
# Validate configuration
./linkwarden-mcp-server --validate-config --config config.dev.yaml
```

#### MCP Protocol Issues

Use the MCP inspector or test client to debug protocol issues:

```bash
# Run with MCP inspector
npx @modelcontextprotocol/inspector stdio /path/to/linkwarden-mcp-server \
  --base-url https://your-linkwarden-instance.com \
  --token your-token
```

## Contributing

### Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `make test`
6. Update documentation
7. Submit a pull request

### Commit Message Format

Use conventional commit messages:

```
feat: add new collection management tool
fix: resolve parameter validation issue
docs: update README with new features
test: add integration tests for search functionality
```

### Code Review Process

1. All code changes require review
2. Automated tests must pass
3. Code must meet linting standards
4. Documentation must be updated
5. Breaking changes require major version bump

## Performance Considerations

### API Optimization

- Use pagination for large result sets
- Cache frequently accessed data
- Implement rate limiting for API calls
- Use connection pooling for HTTP clients

### Memory Management

- Reuse validator instances
- Pool JSON encoding/decoding resources
- Limit concurrent operations
- Monitor memory usage in production

## Deployment

### Local Development

```bash
# Run with hot reload (requires air)
make dev

# Run with Docker
make docker-run
```

### Production Deployment

```bash
# Build optimized binary
make build-prod

# Create systemd service
sudo systemctl enable linkwarden-mcp-server
sudo systemctl start linkwarden-mcp-server
```

### Monitoring

- Monitor API response times
- Track error rates
- Log tool usage patterns
- Set up alerts for failures

## Troubleshooting

### Common Issues

#### Build Failures

```bash
# Clear Go module cache
go clean -modcache
make deps
make build
```

#### Test Failures

```bash
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestMyTool ./pkg/linkwardenmcp
```

#### Runtime Issues

```bash
# Check logs
tail -f ./dev.log

# Run with debug mode
./linkwarden-mcp-server --debug --config config.dev.yaml
```

## Getting Help

- **GitHub Issues**: Report bugs or request features
- **Discussions**: Ask questions and share ideas
- **Documentation**: Check the latest docs
- **Examples**: Look at the examples directory

## Resources

- [Go Documentation](https://golang.org/doc/)
- [MCP Protocol Specification](https://modelcontextprotocol.io/)
- [Linkwarden API Documentation](https://docs.linkwarden.app/)
- [Cobra CLI Documentation](https://github.com/spf13/cobra)