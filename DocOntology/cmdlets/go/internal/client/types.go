package client

import "time"

// Status represents the schema lifecycle status
type Status string

const (
	StatusDraft      Status = "draft"
	StatusActive     Status = "active"
	StatusDeprecated Status = "deprecated"
)

// Visibility represents the schema visibility level
type Visibility string

const (
	VisibilityPublic    Visibility = "public"
	VisibilityCommunity Visibility = "community"
	VisibilityPrivate   Visibility = "private"
)

// SchemaType represents the type of JSON schema
type SchemaType string

const (
	SchemaTypeStandard SchemaType = "standard"
	SchemaTypeRegex    SchemaType = "regex"
)

// DocType represents a document type
type DocType struct {
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Country represents a country (ISO 3166-1 alpha-2)
type Country struct {
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Schema represents a JSON schema
type Schema struct {
	PublicID        string                 `json:"publicId"`
	PublicVersionID string                 `json:"publicVersionId"`
	Name            string                 `json:"name"`
	Description     *string                `json:"description,omitempty"`
	Version         int                    `json:"version"`
	Content         map[string]interface{} `json:"content"`
	SchemaType      SchemaType             `json:"schemaType"`
	Status          Status                 `json:"status"`
	CustomerID      *string                `json:"customerId,omitempty"`
	Visibility      Visibility             `json:"visibility"`
	DocTypeCode     string                 `json:"docTypeCode"`
	CountryCode     *string                `json:"countryCode,omitempty"`
	ValidFrom       *time.Time             `json:"validFrom,omitempty"`
	ValidTo         *time.Time             `json:"validTo,omitempty"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
}

// SchemaWithRelations extends Schema with related entities
type SchemaWithRelations struct {
	Schema
	DocType DocType  `json:"docType"`
	Country *Country `json:"country,omitempty"`
}

// Pagination holds pagination information
type Pagination struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"hasMore"`
}

// PaginatedSchemaList represents a paginated list of schemas
type PaginatedSchemaList struct {
	Data       []SchemaWithRelations `json:"data"`
	Pagination Pagination            `json:"pagination"`
}

// CreateDocTypeRequest is the request body for creating a doc type
type CreateDocTypeRequest struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// UpdateDocTypeRequest is the request body for updating a doc type
type UpdateDocTypeRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"isActive,omitempty"`
}

// CreateCountryRequest is the request body for creating a country
type CreateCountryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// UpdateCountryRequest is the request body for updating a country
type UpdateCountryRequest struct {
	Name     *string `json:"name,omitempty"`
	IsActive *bool   `json:"isActive,omitempty"`
}

// CreateSchemaRequest is the request body for creating a schema
type CreateSchemaRequest struct {
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	Content     map[string]interface{} `json:"content"`
	DocTypeCode string                 `json:"docTypeCode"`
	CountryCode *string                `json:"countryCode,omitempty"`
	Visibility  *Visibility            `json:"visibility,omitempty"`
	SchemaType  *SchemaType            `json:"schemaType,omitempty"`
	CustomerID  *string                `json:"customerId,omitempty"`
}

// UpdateSchemaRequest is the request body for updating a schema
type UpdateSchemaRequest struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Content     map[string]interface{} `json:"content,omitempty"`
	DocTypeCode *string                `json:"docTypeCode,omitempty"`
	CountryCode *string                `json:"countryCode,omitempty"`
	Visibility  *Visibility            `json:"visibility,omitempty"`
	SchemaType  *SchemaType            `json:"schemaType,omitempty"`
}

// FindBestRequest is the request body for finding the best schema
type FindBestRequest struct {
	DocTypeCode string  `json:"docTypeCode"`
	CountryCode *string `json:"countryCode,omitempty"`
	CustomerID  *string `json:"customerId,omitempty"`
}

// FindBestResponse is the response from find-best endpoint
type FindBestResponse struct {
	Schema    *SchemaWithRelations `json:"schema"`
	MatchType *string              `json:"matchType"`
}

// ReferenceDataResponse contains all reference data
type ReferenceDataResponse struct {
	DocTypes  []DocType  `json:"docTypes"`
	Countries []Country  `json:"countries"`
}

// HealthResponse is the response from the health endpoint
type HealthResponse struct {
	Status    string    `json:"status"`
	Database  string    `json:"database"`
	Timestamp time.Time `json:"timestamp"`
}

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Classification represents document classification result
type Classification struct {
	DocType string `json:"doctype"`
	Country string `json:"country"`
	Pages   []int  `json:"pages"`
}

// MatchedSchema represents a schema returned from match endpoint
type MatchedSchema struct {
	PublicID        string                 `json:"publicId"`
	PublicVersionID string                 `json:"publicVersionId"`
	Name            string                 `json:"name"`
	SchemaType      SchemaType             `json:"schemaType"`
	Content         map[string]interface{} `json:"content"`
}

// ExtractResponse is the response from the match endpoint
type ExtractResponse struct {
	Classification Classification `json:"classification"`
	Schema         *MatchedSchema `json:"schema"`
}

// ListSchemasOptions contains options for listing schemas
type ListSchemasOptions struct {
	Status     *Status
	DocType    *string
	Country    *string
	Visibility *Visibility
	CustomerID *string
	Limit      int
	Offset     int
}

// GenerateSchemaRequest contains options for generating a schema
type GenerateSchemaRequest struct {
	FilePath    string  // Path to PDF or JPEG file (optional if Text is provided)
	Text        string  // Raw text content (optional if FilePath is provided)
	DocTypeCode string  // Required: document type code
	CountryCode string  // Required: country code (2 chars)
	UseOCR      bool    // Whether to use OCR on file (default: true)
}

// GeneratedSchema represents the schema object in generate response
type GeneratedSchema struct {
	Content   map[string]interface{} `json:"content"`
	Generated bool                   `json:"generated"`
}

// GenerateResponse is the response from the generate endpoint
type GenerateResponse struct {
	DocType string          `json:"docType"`
	Country string          `json:"country"`
	Schema  GeneratedSchema `json:"schema"`
}
