#!/bin/bash

# Linkwarden SDK Generation Script
# This script generates the Go SDK from the OpenAPI specification using oapi-codegen

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
OPENAPI_SPEC="api/linkwarden.openapi.yaml"
OUTPUT_DIR="pkg/linkwarden"
PACKAGE_NAME="linkwarden"

echo -e "${YELLOW}Generating Linkwarden Go SDK...${NC}"

# Check if oapi-codegen is installed
if ! command -v oapi-codegen &> /dev/null; then
    echo -e "${RED}Error: oapi-codegen is not installed${NC}"
    echo "Install it with: go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest"
    exit 1
fi

# Check if OpenAPI spec exists
if [ ! -f "$OPENAPI_SPEC" ]; then
    echo -e "${RED}Error: OpenAPI specification not found at $OPENAPI_SPEC${NC}"
    exit 1
fi

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Clean up old generated files
echo -e "${YELLOW}Cleaning up old generated files...${NC}"
rm -f "$OUTPUT_DIR"/*.go

# Generate the complete SDK using config file
echo -e "${YELLOW}Generating SDK using oapi-codegen config...${NC}"
oapi-codegen -config linkwarden.codegen.yaml "$OPENAPI_SPEC"

# Add package documentation
cat > "$OUTPUT_DIR/doc.go" << 'EOF'
// Package linkwarden provides a Go client for the Linkwarden API.
//
// This package is generated from the Linkwarden OpenAPI specification
// using oapi-codegen. It provides types and client methods for interacting
// with a Linkwarden instance.
//
// Basic usage:
//
//	client, err := linkwarden.NewClientWithResponses("https://your-linkwarden-instance.com")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Use the client to make API calls
//	resp, err := client.GetLinksWithResponse(context.Background(), &linkwarden.GetLinksParams{})
//	if err != nil {
//		log.Fatal(err)
//	}
package linkwarden
EOF

echo -e "${GREEN}✓ SDK generation completed successfully!${NC}"
echo -e "${GREEN}Generated files:${NC}"
ls -la "$OUTPUT_DIR"/*.go | while read -r line; do
    echo -e "  - ${line##*/}"
done

# Run go mod tidy to ensure dependencies are correct
echo -e "${YELLOW}Running go mod tidy...${NC}"
go mod tidy

echo -e "${GREEN}✓ All done! SDK is ready to use.${NC}"
