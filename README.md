# linkwarden-mcp-server

A Model Context Protocol (MCP) server that provides AI assistants with programmatic access to Linkwarden instances. Linkwarden is a self-hosted bookmark management service, and this MCP server enables AI agents to interact with Linkwarden's bookmark collections, links, and search functionality.

## Features

- **Collection Management**: Create, read, and delete collections with full API support
- **Link Management**: Create, read, archive, and delete links with comprehensive functionality
- **Advanced Search**: Search links with powerful filtering and pagination
- **Public Collection Access**: Access public collections and their metadata
- **Toolset Selectivity**: Enable only the tools you need
- **Read-Only Mode**: Optional safety mode for production environments
- **Flexible Configuration**: Support for command-line flags, environment variables, and config files

## Installation

### Prerequisites

- Go 1.23 or later
- Access to a Linkwarden instance (self-hosted or cloud)
- Linkwarden API token with appropriate permissions

### Build from Source

```bash
git clone https://github.com/irfansofyana/linkwarden-mcp-server.git
cd linkwarden-mcp-server
make build
```

### Using Go

```bash
go install github.com/irfansofyana/linkwarden-mcp-server@latest
```

## Configuration

The server can be configured using command-line flags, environment variables, or a configuration file.

### Required Configuration

- `--base-url`: Your Linkwarden instance URL (or `LINKWARDEN_BASE_URL` environment variable)
- `--token`: Your Linkwarden API token (or `LINKWARDEN_TOKEN` environment variable)

### Optional Configuration

- `--toolsets`: Comma-separated list of toolsets to enable (default: all)
- `--read-only`: Enable read-only mode (disables write operations)
- `--log-file`: Path to log file

### Examples

#### Command Line Flags

```bash
./linkwarden-mcp-server \
  --base-url https://your-linkwarden-instance.com \
  --token your-api-token-here \
  --toolsets search,collection,link
```

#### Environment Variables

```bash
export LINKWARDEN_BASE_URL=https://your-linkwarden-instance.com
export LINKWARDEN_TOKEN=your-api-token-here
export TOOLSETS=search,collection,link
./linkwarden-mcp-server
```


## Available Toolsets

### Collection Toolset

**Read Operations:**
- `get_all_collections`: Retrieve all collections
- `get_collection_by_id`: Get specific collection details
- `get_public_collections_links`: Get links from public collections
- `get_public_collections_tags`: Get tags from public collections
- `get_public_collection_by_id`: Get public collection by ID

**Write Operations:**
- `create_collection`: Create new collections
- `delete_collection_by_id`: Delete existing collections

### Link Toolset

**Read Operations:**
- `get_all_links`: Retrieve all links with filtering and pagination
- `get_link_by_id`: Get specific link details

**Write Operations:**
- `create_link`: Create new links with metadata and tags
- `delete_link_by_id`: Delete existing links
- `delete_links`: Delete multiple links by IDs
- `archive_link`: Archive links by ID

### Search Toolset

- `search_links`: Search links with various filters including:
  - Search query string
  - Sorting and pagination
  - Collection ID filtering
  - Tag ID filtering

## Features

### Search Capabilities
- Full-text search across links
- Search by name, URL, description, text content, and tags
- Filter by collection ID and tag ID
- Sort results and pagination support
- Search within public collections

### Link Management
- List all links with comprehensive filtering options
- Get link details by ID
- Create new links with rich metadata (name, URL, description, tags, collection)
- Delete individual or multiple links
- Archive links for preservation
- Support for link organization with collections and tags

### Collection Management
- List all collections
- Get collection details by ID
- Create new collections with metadata (name, description, color, icon)
- Delete collections by ID
- Support for nested collections (parent-child relationships)

### Public Collection Access
- Access public collections without authentication
- Retrieve links from public collections
- Get tags from public collections
- Search within public collections
- Advanced filtering options for public content

## Usage Examples

### Enable All Tools
```bash
./linkwarden-mcp-server \
  --base-url https://your-linkwarden-instance.com \
  --token your-api-token-here
```

### Read-Only Mode
```bash
./linkwarden-mcp-server \
  --base-url https://your-linkwarden-instance.com \
  --token your-api-token-here \
  --read-only
```

### Enable Specific Toolsets
```bash
./linkwarden-mcp-server \
  --base-url https://your-linkwarden-instance.com \
  --token your-api-token-here \
  --toolsets search,link
```

## MCP Client Integration

### Claude Desktop

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "linkwarden": {
      "command": "/path/to/linkwarden-mcp-server",
      "args": ["stdio"],
      "env": {
        "LINKWARDEN_BASE_URL": "https://your-linkwarden-instance.com",
        "LINKWARDEN_TOKEN": "your-api-token-here"
      }
    }
  }
}
```

### Other MCP Clients

The server uses stdio transport, so it can be integrated with any MCP client that supports stdio communication.

## Development

### Prerequisites

- Go 1.23 or later
- Make
- Docker (for certain development tasks)

### Setup

```bash
git clone https://github.com/irfansofyana/linkwarden-mcp-server.git
cd linkwarden-mcp-server
make deps
```

### Building

```bash
make build          # Build the binary
make clean          # Clean build artifacts
make install        # Install to GOPATH
```

### Testing

```bash
make test           # Run all tests
make test-unit      # Run unit tests only
make test-integration # Run integration tests
```

### Code Quality

```bash
make lint           # Run linter
make fmt            # Format code
```

### SDK Generation

The Linkwarden client SDK is automatically generated from the OpenAPI specification:

```bash
make generate-sdk   # Generate SDK from OpenAPI spec
```

## Architecture

### Core Components

- **Server**: MCP server implementation with stdio transport
- **Toolsets**: Modular system for organizing functionality
- **Validation**: Comprehensive parameter validation and error handling
- **Client**: Auto-generated Linkwarden API client

### Toolset System

The server uses a modular toolset system that allows:
- Selective enabling of functionality
- Read-only mode for safety
- Organized tool grouping
- Easy extension with new toolsets

### Transport

Currently supports stdio transport for communication with MCP clients.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

## Support

For issues and questions:
- GitHub Issues: [Report bugs or request features](https://github.com/irfansofyana/linkwarden-mcp-server/issues)
- Documentation: [Project documentation](https://github.com/irfansofyana/linkwarden-mcp-server/wiki)

## Changelog

### Latest Changes
- Added comprehensive link management toolset
- Implemented link creation, deletion, and archiving functionality
- Enhanced collection management with full CRUD operations
- Implemented public collection access tools
- Enhanced search functionality with advanced filtering
- Added read-only mode for production safety
- Improved parameter validation and error handling