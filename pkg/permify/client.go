package permify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/time/rate"
)

var (
	ErrUnableToCreateRelationship = errors.New("failed to create relationship")
	ErrUnableToDeleteRelationship = errors.New("failed to delete relationship")
	ErrUnableToFindRelationships  = errors.New("failed to find relationships")
	ErrUnableToLookupRelationship = errors.New("failed to lookup relationship")
	ErrUnableToCheckRelationship  = errors.New("failed to check relationship")
	ErrUnableToWriteSchema        = errors.New("failed to update model")
	ErrUnableToCreateTenant       = errors.New("failed to create tenant")
	ErrUnableToDeleteTenant       = errors.New("failed to delete tenant")
	ErrUnableToListTenant         = errors.New("failed to list tenants")
	ErrBodyDecodeFailure          = errors.New("failed to decode response body")
	ErrRateLimitExceeded          = errors.New("rate limit exceeded")
)

const (
	// match default rate limit of the service
	// and default Postgres max connections
	// (pushing to the limit :-D )
	DefaultRateLimit = 100
)

// RelationshipClient represents the behavior of a client managing relationships.
type RelationshipClient interface {
	// AddRelationship establishes a relationship between a subject and an entity.
	// On success, it returns a snapshot of the relationship graph, represented by
	// the RelationshipSnap structure. If the addition fails, an error is returned.
	AddRelationship(ctx context.Context, request *AddRelationshipRequest) (*RelationshipSnap, error)

	// LookupRelationship retrieves information about a specific relationship
	// between a subject and an entity. It returns the details in the form of a
	// LookupRelationshipResponse structure. If the lookup fails or the relationship
	// is not found, an error is returned.
	LookupRelationship(ctx context.Context, request *LookupRelationshipRequest) (*LookupRelationshipResponse, error)

	// FindRelationships identifies all relationships associated with a given subject
	// or entity. It returns a collection of relationships in the FoundRelationshipsResponse
	// structure. If there's an issue during the process, an error is returned.
	FindRelationships(ctx context.Context, request *FindRelationshipsRequest) (*FindRelationshipsResponse, error)

	// DeleteRelationship removes a specific relationship based on the given criteria.
	// If successful, it returns nil. If the deletion fails or the relationship
	// isn't found, an error is returned.
	DeleteRelationship(ctx context.Context, filter *DeleteRelationshipRequest) error

	// CheckPermission verifies if a subject has a specific permission or role on
	// an entity. It returns true if the permission is granted and false otherwise.
	// If there's an issue during the check, an error is returned.
	CheckPermission(ctx context.Context, subject *Subject, entity *Entity, roleOrPermission string) (bool, error)
}

// SchemaManagerClient represents the behavior of a client managing schemas.
// This will typically be used by our deployment application to create and
// manage out and authorization schema model.
type SchemaManagerClient interface {
	// Creates a new Tenant in the Permify Authz service.
	CreateTenant(ctx context.Context, tenant *CreateTenantRequest) (*CreateTenantResponse, error)

	// Deletes an existing Tenant at the Permify Authz service.
	DeleteTenant(ctx context.Context, tenantID string) (*DeleteTenantResponse, error)

	// Lists active Tenants in the Permify Authz service.
	ListTenants(ctx context.Context, request *ListTenantsRequest) (*ListTenantsResponse, error)

	// SaveModelSchema saves the provided schema to the permify authz service.
	SaveModelSchema(ctx context.Context, schema *SaveSchemaRequest) (*SaveSchemaResponse, error)
}

var _ RelationshipClient = (*client)(nil)
var _ SchemaManagerClient = (*client)(nil)

type client struct {
	config  *Config      // Client configuration
	client  *http.Client // HTTP client for making requests
	limiter *rate.Limiter
}

// Config defines the configuration parameters for the client.
type Config struct {
	Protocol   string       // Communication protocol: "http" or "https"
	APIVersion string       // API version like "v1", "v2" etc.
	Host       string       // Host address
	APIKey     string       // Authentication or API key
	Tenant     string       // Tenant identifier
	Client     *http.Client // useful for mocking
	RateLimit  int
}

// NewDefaultConfig returns a default configuration for the client.
func NewDefaultConfig() *Config {
	return &Config{
		Host:       "localhost:3476",
		Tenant:     "t1",
		Protocol:   "http",
		APIVersion: "v1",
		Client:     http.DefaultClient,
		RateLimit:  DefaultRateLimit,
	}
}

// NewClient initializes and returns
// a new RelationshipClient with a built in rate limiter
func NewClient(config *Config) RelationshipClient {
	if config.Host == "" {
		panic("must provide host")
	}
	if config.Tenant == "" {
		config.Tenant = "t1"
	}
	if config.Protocol == "" {
		config.Protocol = "http"
	}
	if config.APIVersion == "" {
		config.APIVersion = "v1"
	}
	if config.Client == nil {
		config.Client = http.DefaultClient
	}
	if config.RateLimit <= 0 {
		config.RateLimit = DefaultRateLimit
	}

	return &client{
		config:  config,
		client:  config.Client,
		limiter: rate.NewLimiter(rate.Limit(config.RateLimit), 1),
	}
}

// constructURL constructs the API endpoint URL based on the client's config and the provided path format.
func (c *client) constructURL(pathFormat string) string {
	return fmt.Sprintf("%s://%s"+pathFormat, c.config.Protocol, c.config.Host, c.config.APIVersion, c.config.Tenant)
}

// sendRequest sends a JSON request and returns the response body.
// this is the central point where all APIs make their requests.
// and we ratelimit
func (c *client) sendRequest(ctx context.Context, method, url string, payload interface{}) ([]byte, error) {

	if err := c.limiter.Wait(ctx); err != nil {
		return nil, ErrRateLimitExceeded
	}

	// good to proceed
	var requestBody []byte
	var err error

	// Only marshal and set the request body if the payload is not nil
	if payload != nil {
		requestBody, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Add(ContentTypeHeader, ContentTypeJSON)

	// TODO: Set authentication header

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, nil
}
