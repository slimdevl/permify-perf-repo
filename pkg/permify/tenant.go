package permify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (c *client) CreateTenant(ctx context.Context, tenant *CreateTenantRequest) (*CreateTenantResponse, error) {
	url := fmt.Sprintf("%s://%s"+TenantCreateAPIPath,
		c.config.Protocol, c.config.Host, c.config.APIVersion)

	if tenant.Tenant == "" {
		tenant.Tenant = c.config.Tenant
	}
	if tenant.ID == "" {
		tenant.ID = uuid.New().String()
	}
	if tenant.CreatedAt == "" {
		tenant.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	body, err := c.sendRequest(ctx, http.MethodPost, url, tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to save model schema: %w", err)
	}

	var response CreateTenantResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, ErrBodyDecodeFailure
	}

	if response.ErrorResponse != nil && response.ErrorResponse.Code != 0 {
		return nil, ErrUnableToCreateTenant
	}

	return &response, nil
}

func (c *client) DeleteTenant(ctx context.Context, tenantID string) (*DeleteTenantResponse, error) {
	url := fmt.Sprintf("%s://%s"+TenantDeleteAPIPath,
		c.config.Protocol, c.config.Host, c.config.APIVersion, tenantID)

	body, err := c.sendRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to save model schema: %w", err)
	}

	var response DeleteTenantResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, ErrBodyDecodeFailure
	}

	if response.ErrorResponse != nil && response.ErrorResponse.Code != 0 {
		return nil, ErrUnableToDeleteTenant
	}

	return &response, nil
}

func (c *client) ListTenants(ctx context.Context, request *ListTenantsRequest) (*ListTenantsResponse, error) {
	url := fmt.Sprintf("%s://%s"+TenantListAPIPath,
		c.config.Protocol,
		c.config.Host,
		c.config.APIVersion)
	body, err := c.sendRequest(ctx, http.MethodPost, url, request)
	if err != nil {
		return nil, fmt.Errorf("failed to save model schema: %w", err)
	}

	var response ListTenantsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, ErrBodyDecodeFailure
	}

	if response.ErrorResponse != nil && response.ErrorResponse.Code != 0 {
		return nil, ErrUnableToListTenant
	}

	return &response, nil
}
