package permify_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/slimdevl/repro/pkg/permify"

	"github.com/stretchr/testify/assert"
)

type MockRoundTripper struct {
	response *http.Response
	err      error
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func newMockClient(mockResponse string, statusCode int) *http.Client {
	return &http.Client{
		Transport: &MockRoundTripper{
			response: &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			},
		},
	}
}

func TestAddRelationship(t *testing.T) {
	ctx := context.Background()
	config := permify.NewDefaultConfig()
	config.Client = newMockClient(`{"snap_token":"crackle_pop"}`, http.StatusOK)
	client := permify.NewClient(config)

	// Positive Scenario
	t.Run("Add Relationship Successfully", func(t *testing.T) {

		relationship := &permify.AddRelationshipRequest{
			Metadata: permify.Metadata{},
			Relationships: []*permify.Relationship{
				{
					Entity:   &permify.Entity{Type: "doc", Id: "doc1"},
					Relation: "owner",
					Subject:  &permify.Subject{Type: "user", Id: "user1"},
				},
			},
		}
		snap, err := client.AddRelationship(ctx, relationship)
		assert.Nil(t, err)
		assert.NotNil(t, snap)
		assert.Equal(t, "crackle_pop", snap.SnapToken)
	})

	// // Negative Scenarios

	t.Run("Nil Request", func(t *testing.T) {
		client := permify.NewClient(config)
		_, err := client.AddRelationship(ctx, nil)
		assert.Error(t, err)
		assert.Equal(t, "request is nil", err.Error())
	})

	t.Run("Empty Relationships", func(t *testing.T) {
		client := permify.NewClient(config)
		_, err := client.AddRelationship(ctx, &permify.AddRelationshipRequest{})
		assert.Error(t, err)
		assert.Equal(t, "request contains no relationships", err.Error())
	})

	t.Run("Missing Entity Fields", func(t *testing.T) {
		relationship := &permify.AddRelationshipRequest{
			Metadata: permify.Metadata{},
			Relationships: []*permify.Relationship{
				{
					Entity:   &permify.Entity{},
					Relation: "owner",
					Subject:  &permify.Subject{Type: "user", Id: "user1"},
				},
			},
		}
		client := permify.NewClient(config)
		_, err := client.AddRelationship(ctx, relationship)
		assert.Error(t, err)
	})

	t.Run("Missing Relation Field", func(t *testing.T) {
		relationship := &permify.AddRelationshipRequest{
			Metadata: permify.Metadata{},
			Relationships: []*permify.Relationship{
				{
					Entity:  &permify.Entity{Type: "doc", Id: "doc1"},
					Subject: &permify.Subject{Type: "user", Id: "user1"},
				},
			},
		}
		client := permify.NewClient(config)
		_, err := client.AddRelationship(ctx, relationship)
		assert.Error(t, err)
	})

	t.Run("Missing Subject Fields", func(t *testing.T) {
		relationship := &permify.AddRelationshipRequest{
			Metadata: permify.Metadata{},
			Relationships: []*permify.Relationship{
				{
					Entity:   &permify.Entity{Type: "doc", Id: "doc1"},
					Relation: "owner",
				},
			},
		}
		client := permify.NewClient(config)
		_, err := client.AddRelationship(ctx, relationship)
		assert.Error(t, err)
	})
}
