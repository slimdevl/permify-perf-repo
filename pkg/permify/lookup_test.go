package permify_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/slimdevl/repro/pkg/permify"
	"github.com/stretchr/testify/assert"
)

func TestLookupRelationship(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful Lookup Relationship", func(t *testing.T) {
		mockResponse := `{
			"someField": "someValue", 
			"otherField": "otherValue"
		}`

		config := permify.NewDefaultConfig()
		config.Client = newMockClient(mockResponse, http.StatusOK)
		client := permify.NewClient(config)

		req := &permify.LookupRelationshipRequest{
			Metadata:   permify.Metadata{Schema: "v1", Snap: "snap_example", Depth: 1},
			EntityType: "doc",
			Permission: "read",
			Subject:    &permify.Subject{Type: "user", Id: "user123"},
		}

		resp, err := client.LookupRelationship(ctx, req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("Server Returns Error", func(t *testing.T) {
		config := permify.NewDefaultConfig()
		config.Client = newMockClient("Internal server error", http.StatusInternalServerError)
		client := permify.NewClient(config)

		req := &permify.LookupRelationshipRequest{
			Metadata:   permify.Metadata{Schema: "v1", Snap: "snap_example", Depth: 1},
			EntityType: "doc",
			Permission: "read",
			Subject:    &permify.Subject{Type: "user", Id: "user123"},
		}

		_, err := client.LookupRelationship(ctx, req)
		assert.NotNil(t, err)
	})

	t.Run("JSON Unmarshaling Error", func(t *testing.T) {
		invalidJSON := `{"someField":`
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(invalidJSON, http.StatusOK)
		client := permify.NewClient(config)

		req := &permify.LookupRelationshipRequest{
			Metadata:   permify.Metadata{Schema: "v1", Snap: "snap_example", Depth: 1},
			EntityType: "doc",
			Permission: "read",
			Subject:    &permify.Subject{Type: "user", Id: "user123"},
		}

		_, err := client.LookupRelationship(ctx, req)
		assert.NotNil(t, err)
	})

	t.Run("Missing Fields in Server Response", func(t *testing.T) {
		missingFieldsResponse := `{}`
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(missingFieldsResponse, http.StatusOK)
		client := permify.NewClient(config)

		req := &permify.LookupRelationshipRequest{
			Metadata:   permify.Metadata{Schema: "v1", Snap: "snap_example", Depth: 1},
			EntityType: "doc",
			Permission: "read",
			Subject:    &permify.Subject{Type: "user", Id: "user123"},
		}

		resp, err := client.LookupRelationship(ctx, req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}
