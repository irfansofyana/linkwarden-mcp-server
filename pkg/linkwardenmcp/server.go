package linkwardenmcp

import (
	"context"
	"fmt"

	"github.com/irfansofyana/linkwarden-mcp-server/pkg/contextkey"
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/linkwarden"
)

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
