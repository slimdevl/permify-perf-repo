package permify

const (
	// Content type for JSON requests
	ContentTypeHeader = "Content-Type"
	ContentTypeJSON   = "application/json"
)

const (
	// Base path for managing schema
	SchemaWriteAPIPath = "/%s/tenants/%s/schemas/write"
	SchemaReadAPIPath  = "/%s/tenants/%s/schemas/read"

	// Base path for tenant management
	TenantCreateAPIPath = "/%s/tenants/create"
	TenantDeleteAPIPath = "/%s/tenants/%s"
	TenantListAPIPath   = "/%s/tenants/list"

	// Base path for the check permissions API endpoint
	PermissionCheckAPIPath = "/%s/tenants/%s/permissions/check"
	// Base path for the relationship ADD API endpoint
	RelationshipAPIPath = "/%s/tenants/%s/relationships/write"
	// Base path for the lookup FIND relationship API endpoint
	LookupRelationshipAPIPath = "/%s/tenants/%s/permissions/lookup-entity"
	// Base path for the LIST/find relationship API endpoint
	FindRelationshipsAPIPath = "/%s/tenants/%s/permissions/expand"
	// Base path for the DELETE relationship API endpoint
	DeleteRelationshipAPIPath = "/%s/tenants/%s/relationships/delete"
)

const (
	// Value returned from the auth server when the permission is granted
	CheckResponseAllowed = "CHECK_RESULT_ALLOWED"
)
