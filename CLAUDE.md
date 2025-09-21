# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based MCP (Model Context Protocol) server that provides AI assistants with programmatic access to Linkwarden instances. The server implements a modular toolset system for managing bookmark collections, links, searching, and accessing public collections.

## Development Commands

### Essential Commands
```bash
# Build the project
make build

# Run tests
make test

# Run linter (requires golangci-lint)
make lint

# Format code
make fmt

# Generate SDK from OpenAPI specification
make generate-sdk

# Install dependencies
make install-deps

# Development workflow
make dev  # Clean, generate SDK, and build
```

### Testing
```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration
```

### Build Targets
```bash
# Build binary
make build

# Clean build artifacts
make clean

# Build and install to GOPATH
make install
```

## Architecture

### Core Components

1. **MCP Server** (`cmd/linkwarden-mcp-server/main.go`)
   - Entry point with Cobra CLI
   - Configuration management using Viper
   - Stdio transport implementation
   - Signal handling and graceful shutdown

2. **Toolset System** (`pkg/toolsets/`)
   - Modular organization of functionality
   - Read/write separation for safety
   - Selective toolset enabling
   - Global read-only mode support

3. **Tool Implementation** (`pkg/linkwardenmcp/`)
   - Individual MCP tools with validation
   - Type-safe parameter handling
   - Linkwarden API integration
   - Comprehensive error handling

4. **API Client** (`pkg/linkwarden/`)
   - Auto-generated from OpenAPI spec
   - Type-safe HTTP client
   - Authentication handling

### Key Architectural Patterns

- **Toolset Registration**: Tools are organized into toolsets (`search`, `collection`, `link`) that can be selectively enabled
- **Parameter Validation**: Comprehensive validation system in `tools_param.go` with type safety
- **Read/Write Separation**: Each toolset supports read-only mode for production safety
- **Configuration Flexibility**: Supports command-line flags, environment variables, and config files
- **Observability**: Structured logging with context-aware logging

## Configuration

### Required Configuration
- `--base-url` / `LINKWARDEN_BASE_URL`: Linkwarden instance URL
- `--token` / `LINKWARDEN_TOKEN`: API token

### Optional Configuration
- `--toolsets` / `TOOLSETS`: Comma-separated toolsets to enable
- `--read-only` / `READ_ONLY`: Enable read-only mode
- `--log-file` / `LOG_FILE`: Path to log file

### Configuration Priority
1. Default values
2. Environment variables
3. Command line flags (highest priority)

## Toolsets

### Collection Toolset
**Read Tools:**
- `get_all_collections`: List all collections
- `get_collection_by_id`: Get collection by ID
- `get_public_collections_links`: Get links from public collections
- `get_public_collections_tags`: Get tags from public collections
- `get_public_collection_by_id`: Get public collection by ID

**Write Tools:**
- `create_collection`: Create new collection
- `delete_collection_by_id`: Delete collection

### Link Toolset
**Read Tools:**
- `get_all_links`: Retrieve all links with filtering and pagination
- `get_link_by_id`: Get specific link details

**Write Tools:**
- `create_link`: Create new links with metadata and tags
- `delete_link_by_id`: Delete existing links
- `delete_links`: Delete multiple links by IDs
- `archive_link`: Archive links by ID

### Search Toolset
- `search_links`: Search links with filtering and pagination

## Adding New Tools

1. **Create Tool Implementation** in `pkg/linkwardenmcp/`
2. **Use Validation System** from `tools_param.go`
3. **Register in Toolset** in `tools.go`
4. **Update Documentation** in `docs/tools-reference.md`

### Example Tool Pattern
```go
func NewTool(obs *observability.Observability, client *linkwarden.ClientWithResponses) mcpgo.Tool {
    params := []mcpgo.ToolParameter{
        mcpgo.WithString("name", mcpgo.Description("Parameter description")),
    }

    handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
        args := make(map[string]interface{})
        validator := NewValidator(&req)
        validator.ValidateAndAddRequiredString(args, "name")

        if result, err := validator.HandleErrorsIfAny(); result != nil {
            return result, err
        }

        // API call logic
        resp, err := client.SomeApiCallWithResponse(ctx, args["name"].(string))
        if err != nil {
            return mcpgo.NewToolResultError("API call failed: " + err.Error()), nil
        }

        return mcpgo.NewToolResultJSON(resp.JSON200), nil
    }

    return mcpgo.NewTool("tool_name", "Tool description", params, handler)
}
```

## SDK Generation

The Linkwarden API client is auto-generated from OpenAPI specification:
```bash
# Generate SDK
make generate-sdk

# The generation uses oapi-codegen with configuration in:
# - linkwarden.codegen.yaml
# - api/linkwarden.openapi.yaml
```

## Development Setup

1. **Prerequisites**: Go 1.23+, Make, Linkwarden instance
2. **Setup**: `make install-deps`
3. **Configuration**: Copy `.env.example` and configure your instance
4. **Development**: `make dev`

## Testing Strategy

- **Unit Tests**: Test individual tool functions with mocks
- **Integration Tests**: Test against real Linkwarden instance (requires `TEST_BASE_URL`)
- **Parameter Validation**: Comprehensive validation testing
- **Error Handling**: Test error scenarios and edge cases

## Code Quality

- **Formatting**: Use `make fmt` for consistent formatting
- **Linting**: Use `make lint` (requires golangci-lint)
- **Testing**: Maintain high test coverage
- **Documentation**: Update documentation for new features

## Deployment

The server uses stdio transport for MCP communication. Common deployment patterns:
- **Claude Desktop**: Configure in `claude_desktop_config.json`
- **Development**: Run with `make run`
- **Production**: Use read-only mode and proper logging

## Security Considerations

- **Token Security**: Never commit API tokens, use environment variables
- **Read-Only Mode**: Enable in production for safety
- **Input Validation**: All parameters are validated before API calls
- **Error Handling**: Sensitive information is not exposed in error messages