package permify_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/slimdevl/repro/pkg/permify"

	"github.com/stretchr/testify/assert"
)

func TestDeleteRelationship(t *testing.T) {
	ctx := context.Background()

	// Positive Scenario
	t.Run("Delete Relationship Successfully", func(t *testing.T) {
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(`{"snap_token": "foobar"}`, http.StatusOK)
		client := permify.NewClient(config)
		filter := &permify.DeleteRelationshipRequest{
			Filter: permify.RelationshipFilter{
				Entity: permify.EntityIDSet{
					Type: "doc",
					Ids:  []string{"doc1"},
				},
				Relation: "owner",
				Subject: permify.SubjectIDSet{
					Type:     "user",
					Ids:      []string{"user1"},
					Relation: "member",
				},
			},
		}
		err := client.DeleteRelationship(ctx, filter)
		assert.NoError(t, err)
	})

	// Negative Scenarios
	t.Run("Nil Filter", func(t *testing.T) {
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(`{"code": 1}`, http.StatusBadRequest)
		client := permify.NewClient(config)
		err := client.DeleteRelationship(ctx, nil)
		assert.Error(t, err)
		assert.Equal(t, "filter is nil", err.Error())
	})

	t.Run("Invalid Entity in Filter", func(t *testing.T) {
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(`{"code": 1}`, http.StatusBadRequest)
		client := permify.NewClient(config)
		filter := &permify.DeleteRelationshipRequest{
			Filter: permify.RelationshipFilter{
				Entity:   permify.EntityIDSet{},
				Relation: "owner",
				Subject: permify.SubjectIDSet{
					Type:     "user",
					Ids:      []string{"user1"},
					Relation: "member",
				},
			},
		}
		err := client.DeleteRelationship(ctx, filter)
		assert.Error(t, err)
		assert.Equal(t, "invalid entity in filter", err.Error())
	})

	t.Run("Relation Not Specified in Filter", func(t *testing.T) {
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(`{"code": 1}`, http.StatusBadRequest)
		client := permify.NewClient(config)

		filter := &permify.DeleteRelationshipRequest{
			Filter: permify.RelationshipFilter{
				Entity: permify.EntityIDSet{
					Type: "doc",
					Ids:  []string{"doc1"},
				},
				Subject: permify.SubjectIDSet{
					Type:     "user",
					Ids:      []string{"user1"},
					Relation: "member",
				},
			},
		}
		err := client.DeleteRelationship(ctx, filter)
		assert.Error(t, err)
		assert.Equal(t, "relation is not specified in filter", err.Error())
	})
}
