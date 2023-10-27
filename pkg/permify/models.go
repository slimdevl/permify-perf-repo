package permify

// Entity represents a generic object with a type and identifier.
type Entity struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

// Subject is an alias for Entity, representing an entity that is the target or receiver of a relationship.
type Subject struct {
	Type string `json:"type" validate:"required"`
	Id   string `json:"id" validate:"required"`
}

// Relationship defines a relation between two entities.
type Relationship struct {
	Entity   *Entity  `json:"entity" validate:"required"`
	Relation string   `json:"relation" validate:"required"`
	Subject  *Subject `json:"subject" validate:"required"`
}

// RelationshipFilter contains filters for relationships.
type RelationshipFilter struct {
	Entity   EntityIDSet  `json:"entity"`
	Relation string       `json:"relation"`
	Subject  SubjectIDSet `json:"subject"`
}

// EntityIDSet represents a set of entity IDs.
type EntityIDSet struct {
	Type string   `json:"type"`
	Ids  []string `json:"ids"`
}

// SubjectIDSet represents a set of subject IDs.
type SubjectIDSet struct {
	Type     string   `json:"type"`
	Ids      []string `json:"ids"`
	Relation string   `json:"relation"`
}

// ErrorResponse represents the standard structure for error responses from the API.
type ErrorResponse struct {
	Code    int      `json:"code"`    // The error code
	Message string   `json:"message"` // A human-readable error message
	Details []string `json:"details"` // Additional details or reasons for the error
}

// RelationshipSnap contains a snapshot token, typically used for concurrency control or versioning.
type RelationshipSnap struct {
	*ErrorResponse `json:",inline"`
	SnapToken      string `json:"snap_token"`
}

// Metadata encapsulates meta-information related to a request or response.
type Metadata struct {
	Schema string `json:"schema_version,omitempty"` // The version of the schema used in the payload
	Snap   string `json:"snap_token,omitempty"`     // RE: https://github.com/Permify/permify/blob/master/docs/docs/reference/snap-tokens.md
	Depth  int    `json:"depth,omitempty"`          // Used when checking permissions to limit the depth of the graph search
}

type PermissionCheckResponseMetadata struct {
	CheckCount int `json:"check_count"`
}

type Tenant struct {
	ID        string `json:"id"`
	Tenant    string `json:"tenant"`
	CreatedAt string `json:"created_at,omitempty"`
}

//
// Internal only models
//

type findRelationshipsResponse struct {
	*ErrorResponse `json:",inline"`
	Tree           relationshipsTree `json:"tree"`
}

type relationshipsTree struct {
	Entity     Entity       `json:"entity"`
	Permission string       `json:"permission"`
	Children   discoverySet `json:"leaf"`
}

type discoverySet struct {
	Subjects subjectSet `json:"subjects"`
}

type subjectSet struct {
	Subjects []leaf `json:"subjects"`
}

type leaf struct {
	Type     string `json:"type"`
	Id       string `json:"id"`
	Relation string `json:"relation"`
}
