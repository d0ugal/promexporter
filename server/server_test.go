package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/d0ugal/promexporter/config"
	"github.com/d0ugal/promexporter/metrics"
	"github.com/d0ugal/promexporter/version"
)

// minimalConfig satisfies ConfigInterface with only the fields handleHealth
// reads — server/log/health/web-ui — keeping the test focused on the
// configured-version fallthrough.
type minimalConfig struct {
	server *config.ServerConfig
}

func (m *minimalConfig) GetDisplayConfig() map[string]interface{} {
	return map[string]interface{}{}
}
func (m *minimalConfig) GetLogging() *config.LoggingConfig {
	return &config.LoggingConfig{Level: "info", Format: "json"}
}
func (m *minimalConfig) GetServer() *config.ServerConfig    { return m.server }
func (m *minimalConfig) GetTracing() *config.TracingConfig  { return &config.TracingConfig{} }
func (m *minimalConfig) GetProfiling() *config.ProfilingConfig {
	return &config.ProfilingConfig{}
}

// TestHandleHealth_UsesConfiguredVersionInfo asserts that the /health endpoint
// reports the version supplied via WithVersionInfo rather than the build-time
// version.Get() defaults. Previously these two endpoints could disagree (the
// dashboard at / used the configured info while /health used build-time).
func TestHandleHealth_UsesConfiguredVersionInfo(t *testing.T) {
	cfg := &minimalConfig{
		server: &config.ServerConfig{Host: "127.0.0.1", Port: 0},
	}
	registry := metrics.NewRegistry("server_test_info")

	configured := &version.Info{
		Version:   "v9.9.9-test",
		Commit:    "deadbeef",
		BuildDate: "2026-05-08T00:00:00Z",
	}

	srv := New(cfg, registry, "test-exporter", configured, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if got := body["version"]; got != configured.Version {
		t.Errorf("version: want %q, got %q", configured.Version, got)
	}

	if got := body["commit"]; got != configured.Commit {
		t.Errorf("commit: want %q, got %q", configured.Commit, got)
	}

	if got := body["build_date"]; got != configured.BuildDate {
		t.Errorf("build_date: want %q, got %q", configured.BuildDate, got)
	}
}
