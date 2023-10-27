package permify

// AddRelationshipRequest represents a request payload to create or modify relationships.
type AddRelationshipRequest struct {
	Metadata      Metadata        `json:"metadata"` // Metadata related to the request
	Relationships []*Relationship `json:"tuples"`   // A list of relationships to be created or modified
}

type PermissionCheckRequest struct {
	Metadata   Metadata `json:"metadata"`
	Entity     *Entity  `json:"entity"`
	Permission string   `json:"permission"`
	Subject    *Subject `json:"subject"`
}

type PermissionCheckResponse struct {
	*ErrorResponse `json:",inline"`
	Can            string                          `json:"can"`
	Metadata       PermissionCheckResponseMetadata `json:"metadata"`
}

func (r *PermissionCheckResponse) IsAllowed() bool {
	return r.Can == CheckResponseAllowed
}

//
// Model Management
//

type SaveSchemaRequest struct {
	Schema string `json:"schema"`
}

type SaveSchemaResponse struct {
	*ErrorResponse `json:",inline"`
	SchemaVersion  string `json:"schema_version"`
}

type CreateTenantRequest = Tenant

type CreateTenantResponse struct {
	*ErrorResponse `json:",inline"`
	Tenant         *Tenant `json:"tenant"`
}

type DeleteTenantResponse struct {
	*ErrorResponse `json:",inline"`
	Tenant         *Tenant `json:"tenant"`
}

type ListTenantsRequest struct {
	PageSize        int    `json:"page_size"`
	ContinuousToken string `json:"continuous_token,omitempty"`
}

type ListTenantsResponse struct {
	*ErrorResponse  `json:",inline"`
	Tenants         []*Tenant `json:"tenants"`
	ContinuousToken string    `json:"continuous_token,omitempty"`
}

type LookupRelationshipRequest struct {
	Metadata   Metadata `json:"metadata"`
	EntityType string   `json:"entity_type"`
	Permission string   `json:"permission"`
	Subject    *Subject `json:"subject"`
}

type LookupRelationshipResponse struct {
	*ErrorResponse `json:",inline"`
	EntityIDs      []string `json:"entity_ids"`
}

type FindRelationshipsRequest struct {
	Metadata   Metadata `json:"metadata"`
	Entity     *Entity  `json:"entity"`
	Permission string   `json:"permission"`
}

type FindRelationshipsResponse = LookupRelationshipResponse

type DeleteRelationshipRequest struct {
	Filter RelationshipFilter `json:"filter"`
}
