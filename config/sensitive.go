package config

import (
	"encoding/json"
	"strings"
)

// SensitiveString represents a string value that should be treated as sensitive
// and redacted when displayed in configuration or logs
type SensitiveString struct {
	value string
}

// NewSensitiveString creates a new SensitiveString with the given value
func NewSensitiveString(value string) SensitiveString {
	return SensitiveString{value: value}
}

// String returns a redacted representation of the sensitive string
func (s SensitiveString) String() string {
	if s.value == "" {
		return "[EMPTY]"
	}
	return "[REDACTED]"
}

// Value returns the actual sensitive value (use with caution)
func (s SensitiveString) Value() string {
	return s.value
}

// IsEmpty returns true if the sensitive string is empty
func (s SensitiveString) IsEmpty() bool {
	return s.value == ""
}

// MarshalJSON implements json.Marshaler interface
func (s SensitiveString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON implements json.Unmarshaler interface
func (s *SensitiveString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	s.value = str
	return nil
}

// ConfigDisplay represents configuration data that can be safely displayed
// It automatically handles sensitive fields
type ConfigDisplay struct {
	fields map[string]interface{}
}

// NewConfigDisplay creates a new ConfigDisplay
func NewConfigDisplay() *ConfigDisplay {
	return &ConfigDisplay{
		fields: make(map[string]interface{}),
	}
}

// Add adds a regular configuration field
func (cd *ConfigDisplay) Add(key, value string) *ConfigDisplay {
	cd.fields[key] = value
	return cd
}

// AddSensitive adds a sensitive configuration field
func (cd *ConfigDisplay) AddSensitive(key string, value SensitiveString) *ConfigDisplay {
	cd.fields[key] = value
	return cd
}

// AddInt adds an integer configuration field
func (cd *ConfigDisplay) AddInt(key string, value int) *ConfigDisplay {
	cd.fields[key] = value
	return cd
}

// AddBool adds a boolean configuration field
func (cd *ConfigDisplay) AddBool(key string, value bool) *ConfigDisplay {
	cd.fields[key] = value
	return cd
}

// GetFields returns the configuration fields for display
func (cd *ConfigDisplay) GetFields() map[string]interface{} {
	return cd.fields
}

// ToDisplayMap converts the configuration to a map suitable for display
func (cd *ConfigDisplay) ToDisplayMap() map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range cd.fields {
		switch v := value.(type) {
		case SensitiveString:
			result[key] = v.String()
		default:
			result[key] = value
		}
	}
	return result
}

// IsSensitiveField checks if a field name indicates it contains sensitive data
func IsSensitiveField(fieldName string) bool {
	upperField := strings.ToUpper(fieldName)
	sensitivePatterns := []string{
		"PASSWORD", "PASS", "SECRET", "KEY", "TOKEN", "AUTH",
		"CREDENTIAL", "CRED", "PRIVATE", "SENSITIVE",
	}

	for _, pattern := range sensitivePatterns {
		if strings.Contains(upperField, pattern) {
			return true
		}
	}
	return false
}

// CreateSensitiveStringFromEnv creates a SensitiveString from an environment variable
// if the field name indicates it's sensitive
func CreateSensitiveStringFromEnv(key, value string) interface{} {
	if IsSensitiveField(key) {
		return NewSensitiveString(value)
	}
	return value
}
