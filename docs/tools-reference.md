# Tools Reference

This document provides detailed information about all available tools in the linkwarden-mcp-server, including parameters, usage examples, and return values.

## Collection Toolset

### Read Operations

#### get_all_collections

Retrieves all collections from your Linkwarden instance.

**Parameters:** None

**Returns:**
```json
{
  "collections": [
    {
      "id": 1,
      "name": "My Bookmarks",
      "description": "Personal bookmark collection",
      "color": "#3B82F6",
      "icon": "bookmark",
      "iconWeight": "regular",
      "parentId": null,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**Example Usage:**
```json
{
  "name": "get_all_collections",
  "arguments": {}
}
```

#### get_collection_by_id

Retrieves a specific collection by its ID.

**Parameters:**
- `id` (required, number): The ID of the collection to retrieve

**Returns:**
```json
{
  "id": 1,
  "name": "My Bookmarks",
  "description": "Personal bookmark collection",
  "color": "#3B82F6",
  "icon": "bookmark",
  "iconWeight": "regular",
  "parentId": null,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Example Usage:**
```json
{
  "name": "get_collection_by_id",
  "arguments": {
    "id": 1
  }
}
```

#### get_public_collections_links

Retrieves links from a public collection with advanced filtering options.

**Parameters:**
- `collectionId` (required, number): The ID of the collection to retrieve links for
- `sort` (optional, number): A numeric value to sort the results
- `cursor` (optional, number): A numeric value for pagination
- `pinnedOnly` (optional, boolean): Whether to return only pinned links
- `searchQueryString` (optional, string): A string to filter search results
- `searchByName` (optional, boolean): Whether to search by name
- `searchByUrl` (optional, boolean): Whether to search by URL
- `searchByDescription` (optional, boolean): Whether to search by description
- `searchByTextContent` (optional, boolean): Whether to search by text content
- `searchByTags` (optional, boolean): Whether to search by tags

**Returns:**
```json
{
  "links": [
    {
      "id": 1,
      "name": "Example Website",
      "url": "https://example.com",
      "description": "An example website",
      "tags": ["example", "web"],
      "pinned": true,
      "collectionId": 1,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "cursor": 0,
  "hasMore": false
}
```

**Example Usage:**
```json
{
  "name": "get_public_collections_links",
  "arguments": {
    "collectionId": 1,
    "searchQueryString": "tutorial",
    "searchByName": true,
    "searchByUrl": false,
    "sort": 0
  }
}
```

#### get_public_collections_tags

Retrieves tags from a public collection.

**Parameters:**
- `collectionId` (required, number): The ID of the collection to retrieve tags for

**Returns:**
```json
{
  "tags": [
    {
      "id": 1,
      "name": "tutorial",
      "color": "#10B981",
      "count": 5
    }
  ]
}
```

**Example Usage:**
```json
{
  "name": "get_public_collections_tags",
  "arguments": {
    "collectionId": 1
  }
}
```

#### get_public_collection_by_id

Retrieves a public collection by its ID.

**Parameters:**
- `id` (required, number): The ID of the public collection to retrieve

**Returns:**
```json
{
  "id": 1,
  "name": "Public Collection",
  "description": "A public collection of resources",
  "color": "#8B5CF6",
  "icon": "globe",
  "iconWeight": "regular",
  "parentId": null,
  "isPublic": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Example Usage:**
```json
{
  "name": "get_public_collection_by_id",
  "arguments": {
    "id": 1
  }
}
```

### Write Operations

#### create_collection

Creates a new collection in your Linkwarden instance.

**Parameters:**
- `name` (optional, string): The name of the collection
- `description` (optional, string): The description of the collection
- `color` (optional, string): The color of the collection (hex code)
- `icon` (optional, string): The icon of the collection
- `iconWeight` (optional, string): The weight of the collection's icon
- `parentId` (optional, number): The ID of the parent collection, if applicable

**Returns:**
```json
{
  "id": 2,
  "name": "New Collection",
  "description": "A newly created collection",
  "color": "#EF4444",
  "icon": "folder",
  "iconWeight": "regular",
  "parentId": null,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**Example Usage:**
```json
{
  "name": "create_collection",
  "arguments": {
    "name": "Learning Resources",
    "description": "Collection for learning materials",
    "color": "#10B981",
    "icon": "graduation-cap"
  }
}
```

#### delete_collection_by_id

Deletes a collection by its ID.

**Parameters:**
- `id` (required, number): The ID of the collection to delete

**Returns:**
```json
{
  "message": "Collection deleted successfully"
}
```

**Example Usage:**
```json
{
  "name": "delete_collection_by_id",
  "arguments": {
    "id": 2
  }
}
```

## Search Toolset

#### search_links

Searches for links based on query parameters with advanced filtering options.

**Parameters:**
- `searchQueryString` (optional, string): A string to filter search results
- `sort` (optional, number): A numeric value to sort the search results
- `cursor` (optional, number): A numeric value for pagination
- `collectionId` (optional, number): Filter by collection ID
- `tagId` (optional, number): Filter by tag ID

**Returns:**
```json
{
  "links": [
    {
      "id": 1,
      "name": "Search Result",
      "url": "https://example.com/search-result",
      "description": "A link matching the search criteria",
      "tags": ["search", "example"],
      "pinned": false,
      "collectionId": 1,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "cursor": 0,
  "hasMore": false
}
```

**Example Usage:**
```json
{
  "name": "search_links",
  "arguments": {
    "searchQueryString": "golang tutorial",
    "collectionId": 1,
    "sort": 0
  }
}
```

## Common Response Patterns

### Success Responses

All tools return either:
- JSON data structure specific to the tool
- Success message for write operations

### Error Responses

```json
{
  "content": [
    {
      "type": "text",
      "text": "Validation errors:\n- missing required parameter: id"
    }
  ],
  "isError": true
}
```

Common error types:
- **Validation errors**: Missing or invalid parameters
- **Authentication errors**: Invalid API token
- **Network errors**: Connection issues with Linkwarden instance
- **API errors**: Linkwarden API errors (e.g., collection not found)

## Tool Usage Patterns

### 1. Collection Management Workflow

```json
// List all collections
{
  "name": "get_all_collections",
  "arguments": {}
}

// Create a new collection
{
  "name": "create_collection",
  "arguments": {
    "name": "Project Resources",
    "description": "Resources for current project"
  }
}

// Search links in a specific collection
{
  "name": "search_links",
  "arguments": {
    "collectionId": 3,
    "searchQueryString": "documentation"
  }
}
```

### 2. Public Collection Exploration

```json
// Get public collection details
{
  "name": "get_public_collection_by_id",
  "arguments": {
    "id": 1
  }
}

// Get links from public collection with filtering
{
  "name": "get_public_collections_links",
  "arguments": {
    "collectionId": 1,
    "searchQueryString": "tutorial",
    "pinnedOnly": true
  }
}

// Get tags from public collection
{
  "name": "get_public_collections_tags",
  "arguments": {
    "collectionId": 1
  }
}
```

### 3. Advanced Search Workflow

```json
// Search with multiple filters
{
  "name": "search_links",
  "arguments": {
    "searchQueryString": "golang",
    "collectionId": 1,
    "tagId": 5,
    "sort": 0
  }
}
```

## Best Practices

### 1. Error Handling
- Always check for error responses
- Handle validation errors gracefully
- Implement retry logic for network issues

### 2. Performance
- Use pagination for large result sets
- Cache frequently accessed collections
- Use specific search filters to reduce result size

### 3. Security
- Use read-only mode when only reading data
- Validate user input before passing to tools
- Handle sensitive data appropriately

### 4. Tool Selection
- Enable only necessary toolsets
- Use read-only mode in production
- Monitor API usage and costs

## Integration Examples

### Claude Desktop Integration

```json
{
  "mcpServers": {
    "linkwarden": {
      "command": "/path/to/linkwarden-mcp-server",
      "args": [
        "--base-url", "https://your-linkwarden-instance.com",
        "--token", "your-api-token-here",
        "--toolsets", "search,collection"
      ]
    }
  }
}
```

### Custom MCP Client

```python
import asyncio
from mcp import ClientSession, StdioServerParameters

async def main():
    server_params = StdioServerParameters(
        command="/path/to/linkwarden-mcp-server",
        args=[
            "--base-url", "https://your-linkwarden-instance.com",
            "--token", "your-api-token-here"
        ]
    )

    async with ClientSession(server_params) as session:
        # List all collections
        result = await session.call_tool("get_all_collections", {})
        print(result)

        # Search for links
        result = await session.call_tool("search_links", {
            "searchQueryString": "tutorial"
        })
        print(result)

if __name__ == "__main__":
    asyncio.run(main())
```