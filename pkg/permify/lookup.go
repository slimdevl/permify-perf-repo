package permify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// LookupRelationship retrieves information about a specific relationship
// between a subject and an entity. It returns the details in the form of a
// LookupRelationshipResponse structure. If the lookup fails or the relationship
// is not found, an error is returned.
func (c *client) LookupRelationship(ctx context.Context, request *LookupRelationshipRequest) (*LookupRelationshipResponse, error) {
	// hard limit depth of relationship graph
	// should only be 3 deep at the present time
	request.Metadata.Depth = 100
	url := c.constructURL(LookupRelationshipAPIPath)

	body, err := c.sendRequest(ctx, http.MethodPost, url, request)
	if err != nil {
		return nil, fmt.Errorf("permify request failed: %w", err)
	}

	var response LookupRelationshipResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse lookup relationship response: %w", err)
	}

	if response.ErrorResponse != nil && response.ErrorResponse.Code != 0 {
		return nil, ErrUnableToLookupRelationship
	}

	return &response, nil
}
