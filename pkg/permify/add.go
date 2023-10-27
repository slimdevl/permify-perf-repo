package permify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// AddRelationship establishes a relationship between a subject and an entity.
// On success, it returns a snapshot of the relationship graph, represented by
// the RelationshipSnap structure. If the addition fails, an error is returned.
func (c *client) AddRelationship(ctx context.Context, request *AddRelationshipRequest) (*RelationshipSnap, error) {
	if err := c.validateRelationshipRequest(request); err != nil {
		return nil, err
	}

	url := c.constructURL(RelationshipAPIPath)
	body, err := c.sendRequest(ctx, http.MethodPost, url, request)
	if err != nil {
		return nil, fmt.Errorf("permify request failed: %w", err)
	}

	var response RelationshipSnap
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, ErrBodyDecodeFailure
	}

	if response.ErrorResponse != nil || response.SnapToken == "" {
		return nil, ErrUnableToCreateRelationship
	}

	return &response, nil
}

// validateRelationshipRequest checks the validity of the RelationshipRequest
func (c *client) validateRelationshipRequest(request *AddRelationshipRequest) error {
	if request == nil {
		return fmt.Errorf("request is nil")
	}
	if len(request.Relationships) == 0 {
		return fmt.Errorf("request contains no relationships")
	}
	for i, r := range request.Relationships {
		if err := validate(r); err != nil {
			return fmt.Errorf("relationship %d validation failed: %w", i, err)
		}
	}
	return nil
}

// validate checks if all required fields of a relationship are set correctly.
func validate(r *Relationship) error {
	// Check the Entity fields
	if r.Entity == nil || r.Entity.Type == "" {
		return errors.New("relationship entity type is missing")
	}
	if r.Entity.Id == "" {
		return errors.New("relationship entity ID is missing")
	}

	// Check the Relation field
	if r.Relation == "" {
		return errors.New("relationship relation is missing")
	}

	// Check the Subject fields
	if r.Subject == nil || r.Subject.Type == "" {
		return errors.New("relationship subject type is missing")
	}
	if r.Subject.Id == "" {
		return errors.New("relationship subject ID is missing")
	}

	return nil
}
