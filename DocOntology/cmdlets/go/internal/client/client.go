package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/miguel-bandeira-infosistema/schemactl/internal/config"
)

// Client is the HTTP client for the Schema Registry API
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// New creates a new API client
func New(cfg *config.Config) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(cfg.APIBaseURL, "/"),
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// doRequest performs an HTTP request and decodes the response
func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
		}
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Health checks the API health
func (c *Client) Health() (*HealthResponse, error) {
	var result HealthResponse
	if err := c.doRequest(http.MethodGet, "/health", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetReferenceData gets all reference data
func (c *Client) GetReferenceData() (*ReferenceDataResponse, error) {
	var result ReferenceDataResponse
	if err := c.doRequest(http.MethodGet, "/reference-data", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListDocTypes lists all doc types (admin endpoint includes inactive)
func (c *Client) ListDocTypes(includeInactive bool) ([]DocType, error) {
	path := "/doc-types"
	if includeInactive {
		path = "/admin/doc-types"
	}
	var result []DocType
	if err := c.doRequest(http.MethodGet, path, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetDocType gets a doc type by code
func (c *Client) GetDocType(code string) (*DocType, error) {
	var result DocType
	if err := c.doRequest(http.MethodGet, "/admin/doc-types/"+url.PathEscape(code), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateDocType creates a new doc type
func (c *Client) CreateDocType(req *CreateDocTypeRequest) (*DocType, error) {
	var result DocType
	if err := c.doRequest(http.MethodPost, "/admin/doc-types", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateDocType updates a doc type
func (c *Client) UpdateDocType(code string, req *UpdateDocTypeRequest) (*DocType, error) {
	var result DocType
	if err := c.doRequest(http.MethodPatch, "/admin/doc-types/"+url.PathEscape(code), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteDocType soft deletes a doc type
func (c *Client) DeleteDocType(code string) error {
	return c.doRequest(http.MethodDelete, "/admin/doc-types/"+url.PathEscape(code), nil, nil)
}

// ListCountries lists all countries (admin endpoint includes inactive)
func (c *Client) ListCountries(includeInactive bool) ([]Country, error) {
	path := "/countries"
	if includeInactive {
		path = "/admin/countries"
	}
	var result []Country
	if err := c.doRequest(http.MethodGet, path, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetCountry gets a country by code
func (c *Client) GetCountry(code string) (*Country, error) {
	var result Country
	if err := c.doRequest(http.MethodGet, "/admin/countries/"+url.PathEscape(code), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateCountry creates a new country
func (c *Client) CreateCountry(req *CreateCountryRequest) (*Country, error) {
	var result Country
	if err := c.doRequest(http.MethodPost, "/admin/countries", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateCountry updates a country
func (c *Client) UpdateCountry(code string, req *UpdateCountryRequest) (*Country, error) {
	var result Country
	if err := c.doRequest(http.MethodPatch, "/admin/countries/"+url.PathEscape(code), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteCountry soft deletes a country
func (c *Client) DeleteCountry(code string) error {
	return c.doRequest(http.MethodDelete, "/admin/countries/"+url.PathEscape(code), nil, nil)
}

// ListSchemas lists schemas with optional filtering
func (c *Client) ListSchemas(opts *ListSchemasOptions) (*PaginatedSchemaList, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Status != nil {
			params.Set("status", string(*opts.Status))
		}
		if opts.DocType != nil {
			params.Set("docTypeCode", *opts.DocType)
		}
		if opts.Country != nil {
			params.Set("countryCode", *opts.Country)
		}
		if opts.Visibility != nil {
			params.Set("visibility", string(*opts.Visibility))
		}
		if opts.CustomerID != nil {
			params.Set("customerId", *opts.CustomerID)
		}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Set("offset", strconv.Itoa(opts.Offset))
		}
	}

	path := "/admin/schemas"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result PaginatedSchemaList
	if err := c.doRequest(http.MethodGet, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSchema gets a schema by ID (publicId or publicVersionId)
func (c *Client) GetSchema(id string) (*SchemaWithRelations, error) {
	var result SchemaWithRelations
	if err := c.doRequest(http.MethodGet, "/admin/schemas/"+url.PathEscape(id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSchemaVersion gets a specific schema version by version ID
func (c *Client) GetSchemaVersion(versionID string) (*SchemaWithRelations, error) {
	var result SchemaWithRelations
	if err := c.doRequest(http.MethodGet, "/admin/schemas/versions/"+url.PathEscape(versionID), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSchemaVersions gets all versions of a schema
func (c *Client) GetSchemaVersions(id string) ([]SchemaWithRelations, error) {
	var result []SchemaWithRelations
	if err := c.doRequest(http.MethodGet, "/admin/schemas/"+url.PathEscape(id)+"/versions", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateSchema creates a new schema
func (c *Client) CreateSchema(req *CreateSchemaRequest) (*SchemaWithRelations, error) {
	var result SchemaWithRelations
	if err := c.doRequest(http.MethodPost, "/admin/schemas", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateSchema updates a schema
func (c *Client) UpdateSchema(id string, req *UpdateSchemaRequest) (*SchemaWithRelations, error) {
	var result SchemaWithRelations
	if err := c.doRequest(http.MethodPatch, "/admin/schemas/"+url.PathEscape(id), req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ActivateSchema activates a draft schema
func (c *Client) ActivateSchema(id string) (*SchemaWithRelations, error) {
	var result SchemaWithRelations
	if err := c.doRequest(http.MethodPost, "/admin/schemas/"+url.PathEscape(id)+"/activate", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeprecateSchema deprecates an active schema
func (c *Client) DeprecateSchema(id string) (*SchemaWithRelations, error) {
	var result SchemaWithRelations
	if err := c.doRequest(http.MethodPost, "/admin/schemas/"+url.PathEscape(id)+"/deprecate", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteSchema deletes a draft schema
func (c *Client) DeleteSchema(id string) error {
	return c.doRequest(http.MethodDelete, "/admin/schemas/"+url.PathEscape(id), nil, nil)
}

// FindBestSchema finds the best matching schema
func (c *Client) FindBestSchema(req *FindBestRequest) (*FindBestResponse, error) {
	var result FindBestResponse
	if err := c.doRequest(http.MethodPost, "/schemas/find-best", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// MatchSchema uploads a file and finds a matching schema
func (c *Client) MatchSchema(filePath string, customerID *string) (*ExtractResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/schemas/extract", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	if customerID != nil {
		req.Header.Set("X-Customer-Id", *customerID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
		}
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var result ExtractResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GenerateSchema generates a schema from a document using LLM
func (c *Client) GenerateSchema(req *GenerateSchemaRequest) (*GenerateResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add required fields
	if err := writer.WriteField("docTypeCode", req.DocTypeCode); err != nil {
		return nil, fmt.Errorf("failed to write docTypeCode: %w", err)
	}
	if err := writer.WriteField("countryCode", req.CountryCode); err != nil {
		return nil, fmt.Errorf("failed to write countryCode: %w", err)
	}

	// Add file or text
	if req.FilePath != "" {
		file, err := os.Open(req.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(req.FilePath))
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}

		if _, err := io.Copy(part, file); err != nil {
			return nil, fmt.Errorf("failed to copy file: %w", err)
		}

		// Add useOCR field
		useOCR := "true"
		if !req.UseOCR {
			useOCR = "false"
		}
		if err := writer.WriteField("useOCR", useOCR); err != nil {
			return nil, fmt.Errorf("failed to write useOCR: %w", err)
		}
	} else if req.Text != "" {
		if err := writer.WriteField("text", req.Text); err != nil {
			return nil, fmt.Errorf("failed to write text: %w", err)
		}
	} else {
		return nil, fmt.Errorf("either FilePath or Text must be provided")
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.baseURL+"/schemas/generate", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
		}
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var result GenerateResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
