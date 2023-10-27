package permify_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/slimdevl/repro/pkg/permify"
	"github.com/stretchr/testify/assert"
)

func TestCheckPermission(t *testing.T) {
	ctx := context.Background()

	t.Run("Permission Granted", func(t *testing.T) {
		mockResponse := `{"can": "CHECK_RESULT_ALLOWED", "metadata": {}}`
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(mockResponse, http.StatusOK)
		client := permify.NewClient(config)

		subject := &permify.Subject{Type: "user", Id: "user123"}
		entity := &permify.Entity{Type: "workspace", Id: "ws123"}
		allowed, err := client.CheckPermission(ctx, subject, entity, "view")

		assert.Nil(t, err)
		assert.True(t, allowed)
	})

	t.Run("Permission Denied", func(t *testing.T) {
		mockResponse := `{"allowed": false}`
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(mockResponse, http.StatusOK)
		client := permify.NewClient(config)

		subject := &permify.Subject{Type: "user", Id: "user123"}
		entity := &permify.Entity{Type: "workspace", Id: "ws123"}
		allowed, err := client.CheckPermission(ctx, subject, entity, "view")

		assert.Nil(t, err)
		assert.False(t, allowed)
	})

	t.Run("Server Error", func(t *testing.T) {
		config := permify.NewDefaultConfig()
		config.Client = newMockClient("Internal server error", http.StatusInternalServerError)
		client := permify.NewClient(config)

		subject := &permify.Subject{Type: "user", Id: "user123"}
		entity := &permify.Entity{Type: "workspace", Id: "ws123"}
		_, err := client.CheckPermission(ctx, subject, entity, "view")

		assert.NotNil(t, err)
	})

	t.Run("Unmarshaling Error", func(t *testing.T) {
		invalidJSON := `{"allowed":`
		config := permify.NewDefaultConfig()
		config.Client = newMockClient(invalidJSON, http.StatusOK)
		client := permify.NewClient(config)

		subject := &permify.Subject{Type: "user", Id: "user123"}
		entity := &permify.Entity{Type: "workspace", Id: "ws123"}
		_, err := client.CheckPermission(ctx, subject, entity, "view")

		assert.NotNil(t, err)
	})

	t.Run("Invalid Inputs", func(t *testing.T) {
		config := permify.NewDefaultConfig()
		config.Client = newMockClient("", http.StatusOK)
		client := permify.NewClient(config)

		invalidSubjects := []*permify.Subject{
			nil,
			{Type: "user"},
			{Id: "user123"},
		}
		invalidEntities := []*permify.Entity{
			nil,
			{Type: "workspace"},
			{Id: "ws123"},
		}

		for _, s := range invalidSubjects {
			for _, e := range invalidEntities {
				_, err := client.CheckPermission(ctx, s, e, "view")
				assert.NotNil(t, err)
			}
		}

		validSubject := &permify.Subject{Type: "user", Id: "user123"}
		validEntity := &permify.Entity{Type: "workspace", Id: "ws123"}
		_, err := client.CheckPermission(ctx, validSubject, validEntity, "")
		assert.NotNil(t, err)
	})
}
