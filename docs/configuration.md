# Configuration Guide

The linkwarden-mcp-server supports multiple configuration methods to suit different deployment scenarios and preferences.

## Configuration Methods

### 1. Command Line Flags

```bash
./linkwarden-mcp-server \
  --base-url https://your-linkwarden-instance.com \
  --token your-api-token-here \
  --toolsets search,collection,link \
  --read-only \
  --config /path/to/config.yaml \
  --log-file /var/log/linkwarden-mcp-server.log
```

### 2. Environment Variables

All configuration options can be set using environment variables. For the Linkwarden connection settings, use the specific environment variables `LINKWARDEN_BASE_URL` and `LINKWARDEN_TOKEN`. For other options, convert flag names to uppercase and replace hyphens with underscores:

```bash
export LINKWARDEN_BASE_URL=https://your-linkwarden-instance.com
export LINKWARDEN_TOKEN=your-api-token-here
export TOOLSETS=search,collection,link
export READ_ONLY=true
export CONFIG=/path/to/config.yaml
export LOG_FILE=/var/log/linkwarden-mcp-server.log

./linkwarden-mcp-server
```


## Configuration Options

### Required Options

| Option | Environment Variable | Description | Example |
|--------|---------------------|-------------|---------|
| `--base-url` | `LINKWARDEN_BASE_URL` | Your Linkwarden instance URL | `https://linkwarden.example.com` |
| `--token` | `LINKWARDEN_TOKEN` | Your Linkwarden API token | `your-api-token-here` |

### Optional Options

| Option | Environment Variable | Description | Default | Example |
|--------|---------------------|-------------|---------|---------|
| `--toolsets` | `TOOLSETS` | Comma-separated list of toolsets to enable | `all` | `search,collection,link` |
| `--read-only` | `READ_ONLY` | Enable read-only mode (disables write operations) | `false` | `true` |
| `--log-file` | `LOG_FILE` | Path to log file | - | `/var/log/linkwarden-mcp-server.log` |

## Configuration Priority

Configuration is applied in this order (higher priority overrides lower):

1. **Default values**
2. **Environment variables**
3. **Command line flags** (highest priority)

## Toolsets Configuration

### Available Toolsets

- `search`: Link searching functionality
- `collection`: Collection management operations
- `link`: Link management operations

## Available Features

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

### Examples

#### Enable All Toolsets (Default)
```bash
./linkwarden-mcp-server \
  --base-url https://your-linkwarden-instance.com \
  --token your-api-token-here
```

#### Enable Specific Toolsets
```bash
./linkwarden-mcp-server \
  --base-url https://your-linkwarden-instance.com \
  --token your-api-token-here \
  --toolsets search,link
```

#### Read-Only Mode
```bash
./linkwarden-mcp-server \
  --base-url https://your-linkwarden-instance.com \
  --token your-api-token-here \
  --read-only
```

## Authentication

### API Token Setup

1. Log in to your Linkwarden instance
2. Navigate to Settings → API Tokens
3. Create a new API token with appropriate permissions
4. Copy the token and use it in your configuration

### Required Permissions

The API token needs permissions based on the toolsets you enable:

#### Collection Toolset
- Read collections: `collections:read`
- Create collections: `collections:write`
- Delete collections: `collections:delete`

#### Link Toolset
- Read links: `links:read`
- Create links: `links:write`
- Delete links: `links:delete`
- Archive links: `links:write`

#### Search Toolset
- Search links: `links:read`

## Environment Variables Configuration

For easy setup, you can configure the Linkwarden connection using environment variables:

```bash
export LINKWARDEN_BASE_URL=https://your-linkwarden-instance.com
export LINKWARDEN_TOKEN=your-api-token-here

./linkwarden-mcp-server stdio
```

### Environment Variable Precedence

The environment variables `LINKWARDEN_BASE_URL` and `LINKWARDEN_TOKEN` take precedence over the generic `BASE_URL` and `TOKEN` variables when both are set.

## Environment-Specific Configurations

### Development
```bash
# Use environment variables for local development
export LINKWARDEN_BASE_URL=http://localhost:3000
export LINKWARDEN_TOKEN=dev-token-here

./linkwarden-mcp-server stdio \
  --log-file ./dev.log
```

### Production
```bash
# Use production instance with read-only mode
./linkwarden-mcp-server \
  --base-url https://linkwarden.company.com \
  --token prod-token-here \
  --read-only \
  --log-file /var/log/linkwarden-mcp-server.log
```

### Testing
```bash
# Use testing instance with specific toolsets
./linkwarden-mcp-server \
  --base-url https://test-linkwarden.example.com \
  --token test-token-here \
  --toolsets search \
  --log-file ./test.log
```

## Configuration Validation

The server validates configuration on startup and will report errors for:

- Missing required options (`--base-url`, `--token`)
- Invalid URLs
- Invalid toolset names

## Security Considerations

### Token Security
- Never commit API tokens to version control
- Use environment variables or secure secret management in production
- Rotate tokens regularly
- Use tokens with minimum required permissions

### Network Security
- Use HTTPS in production environments
- Consider firewall rules for Linkwarden instance access
- Monitor API usage and access logs

## Troubleshooting

### Common Issues

#### Invalid Configuration File
```bash
ERROR: failed to load config: invalid configuration file format
```
**Solution**: Check YAML syntax and ensure all required fields are present.

#### Invalid Toolset Name
```bash
ERROR: invalid toolset name: "invalid-toolset"
```
**Solution**: Use only valid toolset names: `search`, `collection`, `link`

#### Authentication Failed
```bash
ERROR: authentication failed: invalid token
```
**Solution**: Verify your API token is correct and has appropriate permissions.

### Debug Mode

For troubleshooting, you can enable verbose logging:

```bash
./linkwarden-mcp-server \
  --base-url https://your-linkwarden-instance.com \
  --token your-api-token-here \
  --log-file ./debug.log
```

Then check the log file for detailed error messages.