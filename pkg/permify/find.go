package permify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// FindRelationships identifies all relationships associated with a given subject
// or entity. It returns a collection of relationships in the FoundRelationshipsResponse
// structure. If there's an issue during the process, an error is returned.
func (c *client) FindRelationships(ctx context.Context, request *FindRelationshipsRequest) (*FindRelationshipsResponse, error) {
	// hard limit depth of relationship graph
	// should only be 3 deep at the present time
	request.Metadata.Depth = 100

	url := c.constructURL(FindRelationshipsAPIPath)

	body, err := c.sendRequest(ctx, http.MethodPost, url, request)
	if err != nil {
		return nil, fmt.Errorf("permify request failed: %w", err)
	}

	var foundResp findRelationshipsResponse
	if err := json.Unmarshal(body, &foundResp); err != nil {
		return nil, fmt.Errorf("failed to parse lookup relationship response: %w", err)
	}

	if foundResp.ErrorResponse != nil {
		return nil, ErrUnableToFindRelationships
	}

	var response FindRelationshipsResponse
	for _, leaf := range foundResp.Tree.Children.Subjects.Subjects {
		response.EntityIDs = append(response.EntityIDs, leaf.Id)
	}

	return &response, nil
}
