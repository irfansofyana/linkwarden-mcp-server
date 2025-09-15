package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/irfansofyana/linkwarden-mcp-server/pkg/linkwarden"
)

func main() {
	// Get configuration from environment variables
	baseURL := os.Getenv("LINKWARDEN_BASE_URL")
	if baseURL == "" {
		baseURL = "https://linkwarden.irfansp.dev" // fallback for testing
	}

	apiToken := os.Getenv("LINKWARDEN_API_TOKEN")
	if apiToken == "" {
		log.Fatal("LINKWARDEN_API_TOKEN environment variable is required")
	}

	fmt.Printf("Testing Linkwarden SDK with base URL: %s\n", baseURL)

	// Create a new Linkwarden client
	client, err := linkwarden.NewClient(baseURL)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a client with responses for easier handling
	clientWithResponses := &linkwarden.ClientWithResponses{
		ClientInterface: client,
	}

	// Create context
	ctx := context.Background()

	// Set up authentication by adding Bearer token to request headers
	authRequestEditor := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+apiToken)
		return nil
	}

	fmt.Println("\n=== Testing Get Links ===")
	testGetLinks(ctx, clientWithResponses, authRequestEditor)

	fmt.Println("\n=== Testing Get Collections ===")
	testGetCollections(ctx, clientWithResponses, authRequestEditor)

	fmt.Println("\n=== Testing Search Links ===")
	testSearchLinks(ctx, clientWithResponses, authRequestEditor)
}

func testGetLinks(ctx context.Context, client *linkwarden.ClientWithResponses, authEditor linkwarden.RequestEditorFn) {
	// Get links with optional parameters
	params := &linkwarden.GetApiV1LinksParams{
		// You can add optional parameters here:
		// Sort:         &sortValue,      // e.g., 0 for newest first
		// Cursor:       &cursorValue,    // for pagination
		// CollectionId: &collectionId,   // filter by collection
		// PinnedOnly:   &pinnedOnly,     // show only pinned links
	}

	// Make the API call to get links
	response, err := client.GetApiV1LinksWithResponse(ctx, params, authEditor)
	if err != nil {
		log.Printf("Failed to get links: %v", err)
		return
	}

	// Check response status
	if response.StatusCode() != 200 {
		log.Printf("Get links failed with status %d: %s", response.StatusCode(), string(response.Body))
		return
	}

	// Handle successful response
	if response.JSON200 != nil && response.JSON200.Response != nil {
		links := *response.JSON200.Response
		fmt.Printf("Successfully retrieved %d links:\n", len(links))

		// Print details of first few links (limit to 3 for testing)
		limit := len(links)
		if limit > 3 {
			limit = 3
		}

		for i := 0; i < limit; i++ {
			link := links[i]
			fmt.Printf("\nLink %d:\n", i+1)

			if link.Name != nil {
				fmt.Printf("  Name: %s\n", *link.Name)
			}
			if link.Url != nil {
				fmt.Printf("  URL: %s\n", *link.Url)
			}
			if link.Description != nil && *link.Description != "" {
				fmt.Printf("  Description: %s\n", *link.Description)
			}
			if link.Collection != nil && link.Collection.Name != nil {
				fmt.Printf("  Collection: %s\n", *link.Collection.Name)
			}
			if link.CreatedAt != nil {
				fmt.Printf("  Created: %s\n", link.CreatedAt.Format("2006-01-02 15:04:05"))
			}
		}
		if len(links) > 3 {
			fmt.Printf("... and %d more links\n", len(links)-3)
		}
	} else {
		fmt.Println("No links found or unexpected response format")
	}
}

func testGetCollections(ctx context.Context, client *linkwarden.ClientWithResponses, authEditor linkwarden.RequestEditorFn) {
	// Make the API call to get collections
	response, err := client.GetAllCollectionsWithResponse(ctx, authEditor)
	if err != nil {
		log.Printf("Failed to get collections: %v", err)
		return
	}

	// Check response status
	if response.StatusCode() != 200 {
		log.Printf("Get collections failed with status %d: %s", response.StatusCode(), string(response.Body))
		return
	}

	// Handle successful response
	if response.JSON200 != nil && response.JSON200.Response != nil {
		collections := *response.JSON200.Response
		fmt.Printf("Successfully retrieved %d collections:\n", len(collections))

		for i, collection := range collections {
			fmt.Printf("\nCollection %d:\n", i+1)
			if collection.Name != nil {
				fmt.Printf("  Name: %s\n", *collection.Name)
			}
			if collection.Description != nil && *collection.Description != "" {
				fmt.Printf("  Description: %s\n", *collection.Description)
			}
			if collection.Id != nil {
				fmt.Printf("  ID: %d\n", *collection.Id)
			}
			if collection.Color != nil {
				fmt.Printf("  Color: %s\n", *collection.Color)
			}
		}
	} else {
		fmt.Println("No collections found or unexpected response format")
	}
}

func testSearchLinks(ctx context.Context, client *linkwarden.ClientWithResponses, authEditor linkwarden.RequestEditorFn) {
	// Test search with a simple query
	query := "github"
	params := &linkwarden.SearchLinksParams{
		SearchQueryString: &query,
	}

	response, err := client.SearchLinksWithResponse(ctx, params, authEditor)
	if err != nil {
		log.Printf("Failed to search links: %v", err)
		return
	}

	// Check response status
	if response.StatusCode() != 200 {
		log.Printf("Search links failed with status %d: %s", response.StatusCode(), string(response.Body))
		return
	}

	// Handle successful response
	if response.JSON200 != nil && response.JSON200.Data != nil && response.JSON200.Data.Links != nil {
		searchResults := *response.JSON200.Data.Links
		fmt.Printf("Search for '%s' returned %d results:\n", query, len(searchResults))

		// Show first 2 results
		limit := len(searchResults)
		if limit > 2 {
			limit = 2
		}

		for i := 0; i < limit; i++ {
			result := searchResults[i]
			fmt.Printf("\nResult %d:\n", i+1)

			if result.Name != nil {
				fmt.Printf("  Name: %s\n", *result.Name)
			}
			if result.Url != nil {
				fmt.Printf("  URL: %s\n", *result.Url)
			}
			if result.Description != nil && *result.Description != "" {
				fmt.Printf("  Description: %s\n", *result.Description)
			}
		}
		if len(searchResults) > 2 {
			fmt.Printf("... and %d more results\n", len(searchResults)-2)
		}
	} else {
		fmt.Println("No search results found or unexpected response format")
		if response.JSON200 != nil && response.JSON200.Message != nil {
			fmt.Printf("Message: %s\n", *response.JSON200.Message)
		}
	}
}
