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
