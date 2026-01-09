package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	APIBaseURL string
	APIKey     string
	Timeout    int
}

// DefaultTimeout is the default request timeout in seconds
const DefaultTimeout = 30

// DefaultAPIURL is the default Schema Registry API URL
const DefaultAPIURL = "https://api.docdigitizer.com/registry"

// Load loads configuration from environment variables
func Load() (*Config, error) {
	apiURL := os.Getenv("SCHEMACTL_API_URL")
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}

	cfg := &Config{
		APIBaseURL: apiURL,
		APIKey:     os.Getenv("SCHEMACTL_API_KEY"),
		Timeout:    DefaultTimeout,
	}

	if timeoutStr := os.Getenv("SCHEMACTL_TIMEOUT"); timeoutStr != "" {
		timeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SCHEMACTL_TIMEOUT: %w", err)
		}
		cfg.Timeout = timeout
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.APIBaseURL == "" {
		return fmt.Errorf("API URL is required. Set SCHEMACTL_API_URL or use --api-url flag")
	}
	return nil
}

// GetDefaultAPIURL returns the default API URL
func GetDefaultAPIURL() string {
	return DefaultAPIURL
}

// WithAPIURL returns a copy of the config with the API URL set
func (c *Config) WithAPIURL(url string) *Config {
	if url != "" {
		c.APIBaseURL = url
	}
	return c
}

// WithAPIKey returns a copy of the config with the API key set
func (c *Config) WithAPIKey(key string) *Config {
	if key != "" {
		c.APIKey = key
	}
	return c
}
