package linkwardenmcp

import (
	"context"

	"github.com/irfansofyana/linkwarden-mcp-server/pkg/linkwarden"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/mcpgo"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/observability"
)

// SearchLinks returns tools for searching links
func SearchLinks(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithString(
			"searchQueryString",
			mcpgo.Description("A string to filter search results."),
		),
		mcpgo.WithNumber(
			"sort",
			mcpgo.Description("A numeric value to sort the search results."),
		),
		mcpgo.WithNumber(
			"cursor",
			mcpgo.Description("A numeric value for pagination."),
		),
		mcpgo.WithNumber(
			"collectionId",
			mcpgo.Description("Filter by collection ID"),
		),
		mcpgo.WithNumber(
			"tagId",
			mcpgo.Description("Filter by tag ID"),
		),
	}

	handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
		client, err := getClientFromContextOrDefault(ctx, client)
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}

		queryParams := make(map[string]interface{})

		validator := NewValidator(&req)
		validator.ValidateAndAddOptionalString(queryParams, "searchQueryString")
		validator.ValidateAndAddOptionalInt(queryParams, "sort")
		validator.ValidateAndAddOptionalInt(queryParams, "cursor")
		validator.ValidateAndAddOptionalInt(queryParams, "collectionId")
		validator.ValidateAndAddOptionalInt(queryParams, "tagId")

		if result, err := validator.HandleErrorsIfAny(); result != nil {
			return result, err
		}

		params := &linkwarden.SearchLinksParams{}
		mappings := []ParameterMapping{
			{Key: "searchQueryString", Target: &params.SearchQueryString, Type: "string"},
			{Key: "sort", Target: &params.Sort, Type: "int"},
			{Key: "cursor", Target: &params.Cursor, Type: "int"},
			{Key: "collectionId", Target: &params.CollectionId, Type: "int"},
			{Key: "tagId", Target: &params.TagId, Type: "int"},
		}
		SetOptionalParameters(queryParams, mappings)

		resp, err := client.SearchLinksWithResponse(ctx, params)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to search links: " + err.Error()), nil
		}

		if resp.JSON200 != nil {
			return mcpgo.NewToolResultJSON(resp.JSON200)
		}

		return mcpgo.NewToolResultError("Failed to search links: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"search_links",
		"Searches for links based on some query parameters.",
		params,
		handler,
	)
}
