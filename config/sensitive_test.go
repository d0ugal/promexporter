package config

import (
	"strings"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

type sensitiveYAMLHolder struct {
	Password SensitiveString `yaml:"password"`
}

func TestSensitiveStringUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"non-empty scalar", "password: hunter2\n", "hunter2"},
		{"quoted scalar", "password: \"p@ss: word\"\n", "p@ss: word"},
		{"empty string", "password: \"\"\n", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var h sensitiveYAMLHolder
			if err := yaml.Unmarshal([]byte(tt.in), &h); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			if h.Password.Value() != tt.want {
				t.Errorf("got %q, want %q", h.Password.Value(), tt.want)
			}
		})
	}
}

func TestSensitiveStringMarshalYAMLRedacts(t *testing.T) {
	h := sensitiveYAMLHolder{Password: NewSensitiveString("hunter2")}

	out, err := yaml.Marshal(h)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if strings.Contains(string(out), "hunter2") {
		t.Fatalf("secret leaked in YAML output: %s", out)
	}

	if !strings.Contains(string(out), "[REDACTED]") {
		t.Errorf("expected redacted marker, got: %s", out)
	}

	empty, err := yaml.Marshal(sensitiveYAMLHolder{})
	if err != nil {
		t.Fatalf("marshal empty: %v", err)
	}

	if !strings.Contains(string(empty), "[EMPTY]") {
		t.Errorf("expected [EMPTY], got: %s", empty)
	}
}
