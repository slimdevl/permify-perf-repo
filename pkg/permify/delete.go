package permify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *client) DeleteRelationship(ctx context.Context, filter *DeleteRelationshipRequest) error {
	if err := c.validateDeleteFilter(filter); err != nil {
		return err
	}

	url := c.constructURL(DeleteRelationshipAPIPath)
	body, err := c.sendRequest(ctx, http.MethodPost, url, filter)
	if err != nil {
		return fmt.Errorf("permify request failed: %w", err)
	}

	var response RelationshipSnap
	if err := json.Unmarshal(body, &response); err != nil {
		return ErrBodyDecodeFailure
	}

	if response.ErrorResponse != nil || response.SnapToken == "" {
		return ErrUnableToDeleteRelationship
	}

	return nil
}

func (c *client) validateDeleteFilter(request *DeleteRelationshipRequest) error {
	if request == nil {
		return fmt.Errorf("filter is nil")
	}
	if request.Filter.Entity.Type == "" || len(request.Filter.Entity.Ids) == 0 {
		return fmt.Errorf("invalid entity in filter")
	}
	if request.Filter.Relation == "" {
		return fmt.Errorf("relation is not specified in filter")
	}
	return nil
}
