package support
package support

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config manages application configuration.
type Config struct {
	values map[string]interface{}
}

// NewConfig creates a new configuration instance.
func NewConfig() *Config {
	return &Config{
		values: make(map[string]interface{}),
	}
}

// LoadEnv loads environment variables from .env file.
func (c *Config) LoadEnv(path string) error {
	return godotenv.Load(path)
}

// Get retrieves a configuration value with a default fallback.
func (c *Config) Get(key string, defaultValue interface{}) interface{} {
	// First check internal values
	if val, ok := c.values[key]; ok {
		return val
	}

	// Then check environment
	parts := strings.Split(key, ".")
	if len(parts) > 0 {
		// Try ENV_VAR_NAME format
		envKey := strings.ToUpper(strings.Join(parts, "_"))
		if val, ok := os.LookupEnv(envKey); ok {
			return val
		}

		// Try dot-separated format
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
	}

	return defaultValue
}

// Set sets a configuration value.
func (c *Config) Set(key string, value interface{}) {
	c.values[key] = value
}

// GetString retrieves a string configuration value.
func (c *Config) GetString(key string, defaultValue string) string {
	val := c.Get(key, defaultValue)
	switch v := val.(type) {
	case string:
		return v
	case nil:
		return defaultValue
	default:
		return defaultValue
	}
}

// GetInt retrieves an integer configuration value.
func (c *Config) GetInt(key string, defaultValue int) int {
	val := c.Get(key, nil)
	switch v := val.(type) {
	case int:
		return v
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

// GetBool retrieves a boolean configuration value.
func (c *Config) GetBool(key string, defaultValue bool) bool {
	val := c.Get(key, nil)
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "true" || v == "1"
	}
	return defaultValue
}

// GetFloat retrieves a float configuration value.
func (c *Config) GetFloat(key string, defaultValue float64) float64 {
	val := c.Get(key, nil)
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	case int:
		return float64(v)
	}
	return defaultValue
}

// All returns all configuration values.
func (c *Config) All() map[string]interface{} {
	return c.values
}
