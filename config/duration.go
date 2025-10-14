package config

import (
	"fmt"
	"time"
)

// Duration represents a time duration that can be parsed from strings
type Duration struct {
	time.Duration
}

// UnmarshalYAML implements custom unmarshaling for duration strings
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}

	switch v := value.(type) {
	case string:
		duration, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("invalid duration format '%s': %w", v, err)
		}

		d.Duration = duration
	case int:
		// Backward compatibility: treat as seconds
		d.Duration = time.Duration(v) * time.Second
	case int64:
		// Backward compatibility: treat as seconds
		d.Duration = time.Duration(v) * time.Second
	default:
		return fmt.Errorf("duration must be a string (e.g., '60s', '1h') or integer (seconds)")
	}

	return nil
}

// Seconds returns the duration in seconds
func (d *Duration) Seconds() int {
	return int(d.Duration.Seconds())
}
