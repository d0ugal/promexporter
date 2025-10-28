package app

import (
	"strings"
	"testing"

	"github.com/d0ugal/promexporter/config"
	"github.com/d0ugal/promexporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// mockConfig implements ConfigInterface for testing
type mockConfig struct {
	logging *config.LoggingConfig
	server  *config.ServerConfig
	tracing *config.TracingConfig
}

func (m *mockConfig) GetDisplayConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func (m *mockConfig) GetLogging() *config.LoggingConfig {
	return m.logging
}

func (m *mockConfig) GetServer() *config.ServerConfig {
	return m.server
}

func (m *mockConfig) GetTracing() *config.TracingConfig {
	if m.tracing == nil {
		return &config.TracingConfig{}
	}

	return m.tracing
}

func TestWithVersionInfo(t *testing.T) {
	// Create a custom registry for this test
	registry := prometheus.NewRegistry()

	// Create a mock config
	mockCfg := &mockConfig{
		logging: &config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		server: &config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
	}

	// Create metrics registry
	metricsRegistry := metrics.NewRegistry("test_exporter_info")

	// Register the metrics with our custom registry
	registry.MustRegister(metricsRegistry.VersionInfo)

	// Test with custom version info
	app := New("test-exporter").
		WithConfig(mockCfg).
		WithMetrics(metricsRegistry).
		WithVersionInfo("v1.2.3", "abc123", "2023-01-01T00:00:00Z").
		Build()

	if app == nil {
		t.Fatal("Expected app to be created, got nil")
	}

	// Check that the version info metric was set correctly
	expectedMetric := `# HELP test_exporter_info Information about the exporter
# TYPE test_exporter_info gauge
test_exporter_info{build_date="2023-01-01T00:00:00Z",commit="abc123",version="v1.2.3"} 1
`

	if err := testutil.GatherAndCompare(registry, strings.NewReader(expectedMetric), "test_exporter_info"); err != nil {
		t.Errorf("Version info metric mismatch:\n%s", err)
	}
}

func TestWithoutVersionInfo(t *testing.T) {
	// Create a custom registry for this test
	registry := prometheus.NewRegistry()

	// Create a mock config
	mockCfg := &mockConfig{
		logging: &config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		server: &config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
	}

	// Create metrics registry
	metricsRegistry := metrics.NewRegistry("test_exporter_info")

	// Register the metrics with our custom registry
	registry.MustRegister(metricsRegistry.VersionInfo)

	// Test without custom version info (should use default)
	app := New("test-exporter").
		WithConfig(mockCfg).
		WithMetrics(metricsRegistry).
		Build()

	if app == nil {
		t.Fatal("Expected app to be created, got nil")
	}

	// Check that the version info metric was set with default values
	expectedMetric := `# HELP test_exporter_info Information about the exporter
# TYPE test_exporter_info gauge
test_exporter_info{build_date="unknown",commit="unknown",version="dev"} 1
`

	if err := testutil.GatherAndCompare(registry, strings.NewReader(expectedMetric), "test_exporter_info"); err != nil {
		t.Errorf("Version info metric mismatch:\n%s", err)
	}
}
