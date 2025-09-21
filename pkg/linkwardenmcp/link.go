package linkwardenmcp

import (
	"context"

	"github.com/irfansofyana/linkwarden-mcp-server/pkg/linkwarden"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/mcpgo"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/observability"
)

// GetAllLinks returns a tool for getting all links with filtering
func GetAllLinks(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithNumber(
			"sort",
			mcpgo.Description("A numeric value to sort the results."),
		),
		mcpgo.WithNumber(
			"cursor",
			mcpgo.Description("A numeric value for pagination."),
		),
		mcpgo.WithNumber(
			"collectionId",
			mcpgo.Description("Filter by collection ID."),
		),
		mcpgo.WithNumber(
			"tagId",
			mcpgo.Description("Filter by tag ID."),
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
		validator.ValidateAndAddOptionalInt(args, "sort")
		validator.ValidateAndAddOptionalInt(args, "cursor")
		validator.ValidateAndAddOptionalInt(args, "collectionId")
		validator.ValidateAndAddOptionalInt(args, "tagId")
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

		params := &linkwarden.GetApiV1LinksParams{}
		mappings := []ParameterMapping{
			{Key: "sort", Target: &params.Sort, Type: "int"},
			{Key: "cursor", Target: &params.Cursor, Type: "int"},
			{Key: "collectionId", Target: &params.CollectionId, Type: "int"},
			{Key: "tagId", Target: &params.TagId, Type: "int"},
			{Key: "pinnedOnly", Target: &params.PinnedOnly, Type: "bool"},
			{Key: "searchQueryString", Target: &params.SearchQueryString, Type: "string"},
			{Key: "searchByName", Target: &params.SearchByName, Type: "bool"},
			{Key: "searchByUrl", Target: &params.SearchByUrl, Type: "bool"},
			{Key: "searchByDescription", Target: &params.SearchByDescription, Type: "bool"},
			{Key: "searchByTextContent", Target: &params.SearchByTextContent, Type: "bool"},
			{Key: "searchByTags", Target: &params.SearchByTags, Type: "bool"},
		}
		SetOptionalParameters(args, mappings)

		resp, err := client.GetApiV1LinksWithResponse(ctx, params)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to get links: " + err.Error()), nil
		}

		if resp.JSON200 != nil {
			return mcpgo.NewToolResultJSON(resp.JSON200)
		}

		return mcpgo.NewToolResultError("Failed to get links: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"get_all_links",
		"Gets all links with optional filtering and pagination.",
		params,
		handler,
	)
}

// GetLinkById returns a tool for getting a link by ID
func GetLinkById(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithNumber(
			"id",
			mcpgo.Description("The ID of the link to retrieve."),
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

		resp, err := client.GetLinkWithResponse(ctx, id)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to get link: " + err.Error()), nil
		}

		if resp.JSON200 != nil {
			return mcpgo.NewToolResultJSON(resp.JSON200)
		}

		return mcpgo.NewToolResultError("Failed to get link: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"get_link_by_id",
		"Gets a link by its ID.",
		params,
		handler,
	)
}

// CreateLink returns a tool for creating a new link
func CreateLink(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithString(
			"name",
			mcpgo.Description("The name of the link."),
		),
		mcpgo.WithString(
			"url",
			mcpgo.Description("The URL of the link."),
		),
		mcpgo.WithString(
			"description",
			mcpgo.Description("The description of the link."),
		),
		mcpgo.WithString(
			"type",
			mcpgo.Description("The type of the link (url, image, pdf)."),
		),
		mcpgo.WithNumber(
			"collectionId",
			mcpgo.Description("The ID of the collection to add the link to."),
		),
		mcpgo.WithString(
			"collectionName",
			mcpgo.Description("The name of the collection to add the link to."),
		),
		mcpgo.WithArray(
			"tags",
			mcpgo.Description("List of tags to add to the link. Each tag should have 'id' and 'name' fields."),
		),
	}

	handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
		client, err := getClientFromContextOrDefault(ctx, client)
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}

		args := make(map[string]interface{})

		validator := NewValidator(&req)
		validator.ValidateAndAddRequiredString(args, "name")
		validator.ValidateAndAddRequiredString(args, "url")
		validator.ValidateAndAddOptionalString(args, "description")
		validator.ValidateAndAddOptionalString(args, "type")
		validator.ValidateAndAddOptionalInt(args, "collectionId")
		validator.ValidateAndAddOptionalString(args, "collectionName")
		validator.ValidateAndAddOptionalArray(args, "tags")

		if result, err := validator.HandleErrorsIfAny(); result != nil {
			return result, err
		}

		body := &linkwarden.CreateLinkJSONRequestBody{}

		// Set required fields
		name := args["name"].(string)
		url := args["url"].(string)
		body.Name = &name
		body.Url = &url

		// Set optional fields
		if description, ok := args["description"].(string); ok {
			body.Description = &description
		} else {
			body.Description = &name
		}

		if typeVal, ok := args["type"].(string); ok {
			typeEnum := linkwarden.CreateLinkJSONBodyType(typeVal)
			body.Type = &typeEnum
		} else {
			typeEnum := linkwarden.CreateLinkJSONBodyType("url")
			body.Type = &typeEnum
		}

		// Set collection (only if provided)
		if collectionId, ok := args["collectionId"].(int64); ok {
			id := int(collectionId)
			body.Collection = &struct {
				Id   *int    `json:"id,omitempty"`
				Name *string `json:"name,omitempty"`
			}{
				Id: &id,
			}
		} else if collectionName, ok := args["collectionName"].(string); ok {
			body.Collection = &struct {
				Id   *int    `json:"id,omitempty"`
				Name *string `json:"name,omitempty"`
			}{
				Name: &collectionName,
			}
		}

		// Set tags (only if provided)
		if tags, ok := args["tags"].([]interface{}); ok {
			tagStructs := make([]struct {
				Id   *int    `json:"id,omitempty"`
				Name *string `json:"name,omitempty"`
			}, len(tags))

			for i, tag := range tags {
				if tagMap, ok := tag.(map[string]interface{}); ok {
					if idVal, ok := tagMap["id"].(int64); ok {
						id := int(idVal)
						tagStructs[i].Id = &id
					}
					if nameVal, ok := tagMap["name"].(string); ok {
						tagStructs[i].Name = &nameVal
					}
				}
			}
			body.Tags = &tagStructs
		}

		resp, err := client.CreateLinkWithResponse(ctx, *body)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to create link: " + err.Error()), nil
		}

		if resp.JSON200 != nil {
			return mcpgo.NewToolResultJSON(resp.JSON200)
		}

		return mcpgo.NewToolResultError("Failed to create link: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"create_link",
		"Creates a new link.",
		params,
		handler,
	)
}

// DeleteLinkById returns a tool for deleting a link by ID
func DeleteLinkById(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithNumber(
			"id",
			mcpgo.Description("The ID of the link to delete."),
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

		resp, err := client.DeleteLinkWithResponse(ctx, id)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to delete link: " + err.Error()), nil
		}

		if resp.StatusCode() == 200 {
			return mcpgo.NewToolResultText("Link deleted successfully"), nil
		}

		return mcpgo.NewToolResultError("Failed to delete link: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"delete_link_by_id",
		"Deletes a link by its ID.",
		params,
		handler,
	)
}

// DeleteLinks returns a tool for deleting multiple links
func DeleteLinks(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithArray(
			"linkIds",
			mcpgo.Description("List of link IDs to delete."),
		),
	}

	handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
		client, err := getClientFromContextOrDefault(ctx, client)
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}

		args := make(map[string]interface{})

		validator := NewValidator(&req)
		validator.ValidateAndAddRequiredArray(args, "linkIds")

		if result, err := validator.HandleErrorsIfAny(); result != nil {
			return result, err
		}

		linkIdsInterface := args["linkIds"].([]interface{})
		linkIds := make([]int, len(linkIdsInterface))

		for i, id := range linkIdsInterface {
			if idVal, ok := id.(int64); ok {
				linkIds[i] = int(idVal)
			}
		}

		body := linkwarden.DeleteLinksJSONRequestBody{
			LinkIds: &linkIds,
		}

		resp, err := client.DeleteLinksWithResponse(ctx, body)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to delete links: " + err.Error()), nil
		}

		if resp.StatusCode() == 200 {
			return mcpgo.NewToolResultText("Links deleted successfully"), nil
		}

		return mcpgo.NewToolResultError("Failed to delete links: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"delete_links",
		"Deletes multiple links by their IDs.",
		params,
		handler,
	)
}

// ArchiveLink returns a tool for archiving a link
func ArchiveLink(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithNumber(
			"id",
			mcpgo.Description("The ID of the link to archive."),
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

		resp, err := client.ArchiveLinkWithResponse(ctx, id)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to archive link: " + err.Error()), nil
		}

		if resp.StatusCode() == 200 {
			return mcpgo.NewToolResultText("Link archived successfully"), nil
		}

		return mcpgo.NewToolResultError("Failed to archive link: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"archive_link",
		"Archives a link by its ID.",
		params,
		handler,
	)
}
