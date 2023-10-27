package permify_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/slimdevl/repro/pkg/permify"
	"github.com/stretchr/testify/assert"
)

func TestFindRelationships(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful Find Relationships", func(t *testing.T) {
		mockResponse := `
		{
			"tree": {
				"entity": {
					"type": "doc",
					"id": "doc123"
				},
				"permission": "read",
				"arguments": [],
				"leaf": {
					"subjects": {
						"subjects": [
							{
								"type": "entity",
								"id": "entity1",
								"relation": ""
							},
							{
								"type": "entity",
								"id": "entity2",
								"relation": ""
							}
						]
					}
				}
			}
		}`
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(mockResponse, http.StatusOK)
		client := permify.NewClient(config)

		req := &permify.FindRelationshipsRequest{
			Metadata: permify.Metadata{
				Schema: "v1",
				Snap:   "snap_example",
				Depth:  1,
			},
			Entity: &permify.Entity{
				Type: "doc",
				Id:   "doc123",
			},
			Permission: "read",
		}

		resp, err := client.FindRelationships(ctx, req)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, []string{"entity1", "entity2"}, resp.EntityIDs)
	})

	t.Run("Server Returns Error", func(t *testing.T) {
		config := permify.NewDefaultConfig()
		config.Client = newMockClient("Internal server error", http.StatusInternalServerError)
		client := permify.NewClient(config)

		req := &permify.FindRelationshipsRequest{
			Metadata: permify.Metadata{
				Schema: "v1",
				Snap:   "snap_example",
				Depth:  1,
			},
			Entity: &permify.Entity{
				Type: "doc",
				Id:   "doc123",
			},
			Permission: "read",
		}

		_, err := client.FindRelationships(ctx, req)
		assert.NotNil(t, err)
	})

	t.Run("JSON Unmarshaling Error", func(t *testing.T) {
		invalidJSON := `{"tree":`
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(invalidJSON, http.StatusOK)
		client := permify.NewClient(config)

		req := &permify.FindRelationshipsRequest{
			Metadata: permify.Metadata{
				Schema: "v1",
				Snap:   "snap_example",
				Depth:  1,
			},
			Entity: &permify.Entity{
				Type: "doc",
				Id:   "doc123",
			},
			Permission: "read",
		}

		_, err := client.FindRelationships(ctx, req)
		assert.NotNil(t, err)
	})

	t.Run("Missing Fields in Server Response", func(t *testing.T) {
		missingFieldsResponse := `{"tree": {}}`
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(missingFieldsResponse, http.StatusOK)
		client := permify.NewClient(config)

		req := &permify.FindRelationshipsRequest{
			Metadata: permify.Metadata{
				Schema: "v1",
				Snap:   "snap_example",
				Depth:  1,
			},
			Entity: &permify.Entity{
				Type: "doc",
				Id:   "doc123",
			},
			Permission: "read",
		}

		resp, err := client.FindRelationships(ctx, req)
		assert.Nil(t, err)
		assert.Empty(t, resp.EntityIDs) // Since there are no subjects, this should be empty
	})
}
