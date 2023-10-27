package permify_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/slimdevl/repro/pkg/permify"
	"github.com/stretchr/testify/assert"
)

func TestSaveModelSchema(t *testing.T) {
	tests := []struct {
		name     string
		mockResp string
		expected *permify.SaveSchemaResponse
	}{
		{
			name:     "successful request",
			mockResp: `{"schema_version":"1.0"}`,
			expected: &permify.SaveSchemaResponse{
				SchemaVersion: "1.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := permify.NewDefaultConfig()
			config.Client = newMockClient(tt.mockResp, http.StatusOK)
			client := permify.NewClient(config).(permify.SchemaManagerClient)
			resp, err := client.SaveModelSchema(context.TODO(),
				&permify.SaveSchemaRequest{
					Schema: `{ "schema": "entity user{}\n"}`,
				})
			assert.NoError(t, err)
			if resp.SchemaVersion != tt.expected.SchemaVersion {
				t.Fatalf("expected schema version %s, got %s", tt.expected.SchemaVersion, resp.SchemaVersion)
			}
		})
	}
}
