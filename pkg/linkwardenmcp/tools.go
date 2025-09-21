package linkwardenmcp

import (
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/linkwarden"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/observability"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/toolsets"
)

func NewToolSets(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
	enabledToolsets []string,
	readonly bool,
) (*toolsets.ToolsetGroup, error) {
	toolsetGroup := toolsets.NewToolsetGroup(readonly)

	search := toolsets.NewToolset("search", "Linkwarden search related tools").
		AddReadTools(SearchLinks(obs, client))

	toolsetGroup.AddToolset(search)

	if err := toolsetGroup.EnableToolsets(enabledToolsets); err != nil {
		return nil, err
	}

	return toolsetGroup, nil
}
