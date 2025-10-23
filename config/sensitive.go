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

