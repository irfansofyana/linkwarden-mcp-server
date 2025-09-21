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

	collection := toolsets.NewToolset("collection", "Linkwarden collection related tools").
		AddReadTools(
			GetAllCollections(obs, client),
			GetCollectionById(obs, client),
			GetPublicCollectionsLinks(obs, client),
			GetPublicCollectionsTags(obs, client),
			GetPublicCollectionById(obs, client),
		).
		AddWriteTools(
			CreateCollection(obs, client),
			DeleteCollectionById(obs, client),
		)

	link := toolsets.NewToolset("link", "Linkwarden link related tools").
		AddReadTools(
			GetAllLinks(obs, client),
			GetLinkById(obs, client),
		).
		AddWriteTools(
			CreateLink(obs, client),
			DeleteLinkById(obs, client),
			DeleteLinks(obs, client),
			ArchiveLink(obs, client),
		)

	toolsetGroup.AddToolset(search)
	toolsetGroup.AddToolset(collection)
	toolsetGroup.AddToolset(link)

	if err := toolsetGroup.EnableToolsets(enabledToolsets); err != nil {
		return nil, err
	}

	return toolsetGroup, nil
}
