package permify_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/slimdevl/repro/pkg/permify"
	"github.com/stretchr/testify/assert"
)

const TenantName = "t1"

func TestCreateTenant(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		status   int
		expected *permify.CreateTenantResponse
	}{
		{
			name:   "successful request",
			id:     TenantName,
			status: http.StatusOK,
			expected: &permify.CreateTenantResponse{
				Tenant: &permify.Tenant{
					ID: TenantName,
				},
			},
		},
		{
			name:   "failed request",
			id:     "fail",
			status: http.StatusInternalServerError,
			expected: &permify.CreateTenantResponse{
				Tenant: &permify.Tenant{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := permify.NewDefaultConfig()
			mockResp, err := marshalExpected(tt.expected)
			assert.NoError(t, err)
			config.Client = newMockClient(mockResp, tt.status) // Assuming newMockClient can handle various status codes
			client := permify.NewClient(config).(permify.SchemaManagerClient)
			resp, err := client.CreateTenant(context.TODO(),
				&permify.CreateTenantRequest{
					ID: tt.id,
				},
			)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Tenant)
			if resp.Tenant.ID != tt.expected.Tenant.ID {
				t.Fatalf("expected tenant id %s, got %s", tt.expected.Tenant.ID, resp.Tenant.ID)
			}

		})
	}
}

func TestDeleteTenant(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected *permify.DeleteTenantResponse
	}{
		{
			name: "successful delete",
			id:   TenantName,
			expected: &permify.DeleteTenantResponse{
				Tenant: &permify.Tenant{
					ID: TenantName,
				},
			},
		},
		// You can add more test cases as needed.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := permify.NewDefaultConfig()
			mockResp, err := marshalExpected(tt.expected)
			assert.NoError(t, err)
			config.Client = newMockClient(mockResp, http.StatusOK)
			client := permify.NewClient(config).(permify.SchemaManagerClient)
			resp, err := client.DeleteTenant(context.TODO(), tt.id)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			if resp.Tenant.ID != tt.expected.Tenant.ID {
				t.Fatalf("expected tenant id %s, got %s", tt.expected.Tenant.ID, resp.Tenant.ID)
			}
		})
	}
}

func TestListTenants(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected *permify.ListTenantsResponse
	}{
		{
			name:   "successful list",
			status: http.StatusOK,
			expected: &permify.ListTenantsResponse{
				Tenants: []*permify.Tenant{
					{ID: TenantName},
					{ID: "john"},
				},
			},
		},
		{
			name:   "empty list",
			status: http.StatusOK,
			expected: &permify.ListTenantsResponse{
				Tenants: []*permify.Tenant{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := permify.NewDefaultConfig()
			mockResp, err := marshalExpected(tt.expected)
			assert.NoError(t, err)
			config.Client = newMockClient(mockResp, tt.status)
			client := permify.NewClient(config).(permify.SchemaManagerClient)
			resp, err := client.ListTenants(context.TODO(), &permify.ListTenantsRequest{})
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, len(tt.expected.Tenants), len(resp.Tenants))
			for i, tenant := range resp.Tenants {
				assert.Equal(t, tt.expected.Tenants[i].ID, tenant.ID)
			}
		})
	}
}

func marshalExpected(v interface{}) (string, error) {
	marshaled, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal expected value: %w", err)
	}
	return string(marshaled), nil
}
