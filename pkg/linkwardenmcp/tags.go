package linkwardenmcp

import (
	"context"

	"github.com/irfansofyana/linkwarden-mcp-server/pkg/linkwarden"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/mcpgo"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/observability"
)

// GetAllTags returns a tool for getting all tags
func GetAllTags(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{}

	handler := func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.ToolResult, error) {
		client, err := getClientFromContextOrDefault(ctx, client)
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}

		resp, err := client.GetTagsWithResponse(ctx)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to get all tags: " + err.Error()), nil
		}

		if resp.JSON200 != nil {
			return mcpgo.NewToolResultJSON(resp.JSON200)
		}

		return mcpgo.NewToolResultError("Failed to get all tags: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"get_all_tags",
		"Gets all tags.",
		params,
		handler,
	)
}

// DeleteTagById returns a tool for deleting a tag by ID
func DeleteTagById(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
) mcpgo.Tool {
	params := []mcpgo.ToolParameter{
		mcpgo.WithNumber(
			"id",
			mcpgo.Description("The ID of the tag to delete."),
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

		resp, err := client.DeleteTagWithResponse(ctx, id)
		if err != nil {
			return mcpgo.NewToolResultError("Failed to delete tag: " + err.Error()), nil
		}

		if resp.StatusCode() == 200 {
			return mcpgo.NewToolResultText("Tag deleted successfully"), nil
		}

		return mcpgo.NewToolResultError("Failed to delete tag: " + resp.Status()), nil
	}

	return mcpgo.NewTool(
		"delete_tag_by_id",
		"Deletes a tag by its ID.",
		params,
		handler,
	)
}