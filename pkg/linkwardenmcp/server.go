package linkwardenmcp

import (
	"context"
	"fmt"

	"github.com/irfansofyana/linkwarden-mcp-server/pkg/contextkey"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/linkwarden"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/mcpgo"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/observability"
)

func NewLinkwardenMcpServer(
	obs *observability.Observability,
	client *linkwarden.ClientWithResponses,
	enabledToolsets []string,
	readOnly bool,
	mcpOpts ...mcpgo.ServerOption,
) (mcpgo.Server, error) {
	if obs == nil {
		return nil, fmt.Errorf("observability is required")
	}

	if client == nil {
		return nil, fmt.Errorf("linkwarden client is required")
	}

	defaultOpts := []mcpgo.ServerOption{
		mcpgo.WithLogging(),
		mcpgo.WithResourceCapabilities(true, true),
		mcpgo.WithToolCapabilities(true),
		mcpgo.WithHooks(mcpgo.SetupHooks(obs))}

	// Merge with user-provided options
	mcpOpts = append(defaultOpts, mcpOpts...)

	server := mcpgo.NewMcpServer("linkwarden-mcp", "0.0.1", mcpOpts...)

	toolsets, err := NewToolSets(obs, client, enabledToolsets, readOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to create toolsets: %w", err)
	}
	toolsets.RegisterTools(server)

	return server, nil
}

// getClientFromContextOrDefault returns either the provided default
// client or gets one from context.
func getClientFromContextOrDefault(
	ctx context.Context,
	defaultClient *linkwarden.ClientWithResponses,
) (*linkwarden.ClientWithResponses, error) {
	if defaultClient != nil {
		return defaultClient, nil
	}

	clientInterface := contextkey.ClientFromContext(ctx)
	if clientInterface == nil {
		return nil, fmt.Errorf("no client found in context")
	}

	client, ok := clientInterface.(*linkwarden.ClientWithResponses)
	if !ok {
		return nil, fmt.Errorf("invalid client type in context")
	}

	return client, nil
}
