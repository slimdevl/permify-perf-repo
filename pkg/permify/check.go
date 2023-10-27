package permify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CheckPermission verifies if a subject has a specific permission or role on
// an entity. It returns true if the permission is granted and false otherwise.
// If there's an issue during the check, an error is returned.
func (c *client) CheckPermission(ctx context.Context, who *Subject, what *Entity, permission string) (bool, error) {
	if err := c.validatePermissionCheckInput(who, what, permission); err != nil {
		return false, err
	}

	request := PermissionCheckRequest{
		Metadata: Metadata{
			Depth: 100, // hard limit for now
		},
		Entity:     what,
		Permission: permission,
		Subject:    who,
	}

	url := c.constructURL(PermissionCheckAPIPath)
	body, err := c.sendRequest(ctx, http.MethodPost, url, request)
	if err != nil {
		return false, err
	}

	var response PermissionCheckResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return false, fmt.Errorf("failed to parse response: %w", err)
	}

	if response.ErrorResponse != nil && response.ErrorResponse.Code != 0 {
		return false, ErrUnableToCheckRelationship
	}

	return response.IsAllowed(), nil
}

// validatePermissionCheckInput validates the inputs for the CheckPermission function.
func (c *client) validatePermissionCheckInput(subject *Subject, entity *Entity, role_or_permission string) error {
	if entity == nil || entity.Id == "" || entity.Type == "" {
		return fmt.Errorf("entity is invalid")
	}
	if subject == nil || subject.Id == "" || subject.Type == "" {
		return fmt.Errorf("subject is invalid")
	}
	if role_or_permission == "" {
		return fmt.Errorf("role_or_permission is invalid")
	}
	return nil
}
