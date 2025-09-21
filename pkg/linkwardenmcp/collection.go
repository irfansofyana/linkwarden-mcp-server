package linkwardenmcp

import (
	"context"

	"github.com/irfansofyana/linkwarden-mcp-server/pkg/linkwarden"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/mcpgo"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/observability"
)

// GetAllCollections returns a tool for getting all collections
func GetAllCollections(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{}

	handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
		client, err := getClientFromContextOrDefault(ctx, client)
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}

		resp, err := client.GetAllCollectionsWithResponse(ctx)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to get all collections: " + err.Error()), nil
		}

		if resp.JSON200 != nil {
			return mcpgo.NewToolResultJSON(resp.JSON200)
		}

		return mcpgo.NewToolResultError("Failed to get all collections: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"get_all_collections",
		"Gets all collections.",
		params,
		handler,
	)
}

// GetCollectionById returns a tool for getting a collection by ID
func GetCollectionById(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithNumber(
			"id",
			mcpgo.Description("The ID of the collection to retrieve."),
		),
	}

	handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
		client, err := getClientFromContextOrDefault(ctx, client)
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}

		args := make(map[string]interface{})

		validator := NewValidator(&req)
		validator.ValidateAndAddRequiredInt(args, "id")

		if result, err := validator.HandleErrorsIfAny(); result != nil {
			return result, err
		}

		id := int(args["id"].(int64))

		resp, err := client.GetCollectionByIdWithResponse(ctx, id)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to get collection: " + err.Error()), nil
		}

		if resp.JSON200 != nil {
			return mcpgo.NewToolResultJSON(resp.JSON200)
		}

		return mcpgo.NewToolResultError("Failed to get collection: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"get_collection_by_id",
		"Gets a collection by its ID.",
		params,
		handler,
	)
}

// CreateCollection returns tools for creating collections
func CreateCollection(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithString(
			"name",
			mcpgo.Description("The name of the collection."),
		),
		mcpgo.WithString(
			"description",
			mcpgo.Description("The description of the collection."),
		),
		mcpgo.WithString(
			"color",
			mcpgo.Description("The color of the collection."),
		),
		mcpgo.WithString(
			"icon",
			mcpgo.Description("The icon of the collection."),
		),
		mcpgo.WithString(
			"iconWeight",
			mcpgo.Description("The weight of the collection's icon."),
		),
		mcpgo.WithNumber(
			"parentId",
			mcpgo.Description("The ID of the parent collection, if applicable."),
		),
	}

	handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
		client, err := getClientFromContextOrDefault(ctx, client)
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}

		args := make(map[string]interface{})

		validator := NewValidator(&req)
		validator.ValidateAndAddOptionalString(args, "name")
		validator.ValidateAndAddOptionalString(args, "description")
		validator.ValidateAndAddOptionalString(args, "color")
		validator.ValidateAndAddOptionalString(args, "icon")
		validator.ValidateAndAddOptionalString(args, "iconWeight")
		validator.ValidateAndAddOptionalInt(args, "parentId")

		if result, err := validator.HandleErrorsIfAny(); result != nil {
			return result, err
		}

		body := &linkwarden.CreateCollectionJSONRequestBody{}
		mappings := []ParameterMapping{
			{Key: "name", Target: &body.Name, Type: "string"},
			{Key: "description", Target: &body.Description, Type: "string"},
			{Key: "color", Target: &body.Color, Type: "string"},
			{Key: "icon", Target: &body.Icon, Type: "string"},
			{Key: "iconWeight", Target: &body.IconWeight, Type: "string"},
			{Key: "parentId", Target: &body.ParentId, Type: "int"},
		}
		SetOptionalParameters(args, mappings)

		resp, err := client.CreateCollectionWithResponse(ctx, *body)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to create collection: " + err.Error()), nil
		}

		if resp.JSON200 != nil {
			return mcpgo.NewToolResultJSON(resp.JSON200)
		}

		return mcpgo.NewToolResultError("Failed to create collection: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"create_collection",
		"Creates a new collection.",
		params,
		handler,
	)
}


// DeleteCollectionById returns a tool for deleting a collection by ID
func DeleteCollectionById(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithNumber(
			"id",
			mcpgo.Description("The ID of the collection to delete."),
		),
	}

	handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
		client, err := getClientFromContextOrDefault(ctx, client)
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}

		args := make(map[string]interface{})

		validator := NewValidator(&req)
		validator.ValidateAndAddRequiredInt(args, "id")

		if result, err := validator.HandleErrorsIfAny(); result != nil {
			return result, err
		}

		id := int(args["id"].(int64))

		resp, err := client.DeleteCollectionByIdWithResponse(ctx, id)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to delete collection: " + err.Error()), nil
		}

		if resp.StatusCode() == 200 {
			return mcpgo.NewToolResultText("Collection deleted successfully"), nil
		}

		return mcpgo.NewToolResultError("Failed to delete collection: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"delete_collection_by_id",
		"Deletes a collection by its ID.",
		params,
		handler,
	)
}

// GetPublicCollectionsLinks returns a tool for getting public collection links
func GetPublicCollectionsLinks(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithNumber(
			"collectionId",
			mcpgo.Description("The ID of the collection to retrieve links for."),
		),
		mcpgo.WithNumber(
			"sort",
			mcpgo.Description("A numeric value to sort the results."),
		),
		mcpgo.WithNumber(
			"cursor",
			mcpgo.Description("A numeric value for pagination."),
		),
		mcpgo.WithBoolean(
			"pinnedOnly",
			mcpgo.Description("Whether to return only pinned links."),
		),
		mcpgo.WithString(
			"searchQueryString",
			mcpgo.Description("A string to filter search results."),
		),
		mcpgo.WithBoolean(
			"searchByName",
			mcpgo.Description("Whether to search by name."),
		),
		mcpgo.WithBoolean(
			"searchByUrl",
			mcpgo.Description("Whether to search by URL."),
		),
		mcpgo.WithBoolean(
			"searchByDescription",
			mcpgo.Description("Whether to search by description."),
		),
		mcpgo.WithBoolean(
			"searchByTextContent",
			mcpgo.Description("Whether to search by text content."),
		),
		mcpgo.WithBoolean(
			"searchByTags",
			mcpgo.Description("Whether to search by tags."),
		),
	}

	handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
		client, err := getClientFromContextOrDefault(ctx, client)
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}

		args := make(map[string]interface{})

		validator := NewValidator(&req)
		validator.ValidateAndAddRequiredInt(args, "collectionId")
		validator.ValidateAndAddOptionalInt(args, "sort")
		validator.ValidateAndAddOptionalInt(args, "cursor")
		validator.ValidateAndAddOptionalBool(args, "pinnedOnly")
		validator.ValidateAndAddOptionalString(args, "searchQueryString")
		validator.ValidateAndAddOptionalBool(args, "searchByName")
		validator.ValidateAndAddOptionalBool(args, "searchByUrl")
		validator.ValidateAndAddOptionalBool(args, "searchByDescription")
		validator.ValidateAndAddOptionalBool(args, "searchByTextContent")
		validator.ValidateAndAddOptionalBool(args, "searchByTags")

		if result, err := validator.HandleErrorsIfAny(); result != nil {
			return result, err
		}

		params := &linkwarden.GetApiV1PublicCollectionsLinksParams{}
		mappings := []ParameterMapping{
			{Key: "collectionId", Target: &params.CollectionId, Type: "int"},
			{Key: "sort", Target: &params.Sort, Type: "int"},
			{Key: "cursor", Target: &params.Cursor, Type: "int"},
			{Key: "pinnedOnly", Target: &params.PinnedOnly, Type: "bool"},
			{Key: "searchQueryString", Target: &params.SearchQueryString, Type: "string"},
			{Key: "searchByName", Target: &params.SearchByName, Type: "bool"},
			{Key: "searchByUrl", Target: &params.SearchByUrl, Type: "bool"},
			{Key: "searchByDescription", Target: &params.SearchByDescription, Type: "bool"},
			{Key: "searchByTextContent", Target: &params.SearchByTextContent, Type: "bool"},
			{Key: "searchByTags", Target: &params.SearchByTags, Type: "bool"},
		}
		SetOptionalParameters(args, mappings)

		resp, err := client.GetApiV1PublicCollectionsLinksWithResponse(ctx, params)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to get public collection links: " + err.Error()), nil
		}

		if resp.JSON200 != nil {
			return mcpgo.NewToolResultJSON(resp.JSON200)
		}

		return mcpgo.NewToolResultError("Failed to get public collection links: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"get_public_collections_links",
		"Gets links from a public collection.",
		params,
		handler,
	)
}

// GetPublicCollectionsTags returns a tool for getting public collection tags
func GetPublicCollectionsTags(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithNumber(
			"collectionId",
			mcpgo.Description("The ID of the collection to retrieve tags for."),
		),
	}

	handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
		client, err := getClientFromContextOrDefault(ctx, client)
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}

		args := make(map[string]interface{})

		validator := NewValidator(&req)
		validator.ValidateAndAddRequiredInt(args, "collectionId")

		if result, err := validator.HandleErrorsIfAny(); result != nil {
			return result, err
		}

		params := &linkwarden.GetApiV1PublicCollectionsTagsParams{}
		mappings := []ParameterMapping{
			{Key: "collectionId", Target: &params.CollectionId, Type: "int"},
		}
		SetOptionalParameters(args, mappings)

		resp, err := client.GetApiV1PublicCollectionsTagsWithResponse(ctx, params)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to get public collection tags: " + err.Error()), nil
		}

		if resp.JSON200 != nil {
			return mcpgo.NewToolResultJSON(resp.JSON200)
		}

		return mcpgo.NewToolResultError("Failed to get public collection tags: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"get_public_collections_tags",
		"Gets tags from a public collection.",
		params,
		handler,
	)
}

// GetPublicCollectionById returns a tool for getting a public collection by ID
func GetPublicCollectionById(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithNumber(
			"id",
			mcpgo.Description("The ID of the public collection to retrieve."),
		),
	}

	handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
		client, err := getClientFromContextOrDefault(ctx, client)
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}

		args := make(map[string]interface{})

		validator := NewValidator(&req)
		validator.ValidateAndAddRequiredInt(args, "id")

		if result, err := validator.HandleErrorsIfAny(); result != nil {
			return result, err
		}

		id := int(args["id"].(int64))

		resp, err := client.GetApiV1PublicCollectionsIdWithResponse(ctx, id)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to get public collection: " + err.Error()), nil
		}

		if resp.JSON200 != nil {
			return mcpgo.NewToolResultJSON(resp.JSON200)
		}

		return mcpgo.NewToolResultError("Failed to get public collection: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"get_public_collection_by_id",
		"Gets a public collection by its ID.",
		params,
		handler,
	)
}
