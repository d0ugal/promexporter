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

func (m *mockConfig) GetProfiling() *config.ProfilingConfig {
	return &config.ProfilingConfig{}
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

// TestGetTracer_NonNilWhenDisabled locks in the contract that GetTracer()
// always returns a usable Tracer after Build(), even when tracing is disabled.
// This lets consumers call tracer.NewCollectorSpan(...) directly without a
// nil-check, simplifying every collector across the cohort.
func TestGetTracer_NonNilWhenDisabled(t *testing.T) {
	registry := prometheus.NewRegistry()

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

	metricsRegistry := metrics.NewRegistry("tracer_test_info")
	registry.MustRegister(metricsRegistry.VersionInfo)

	app := New("test-exporter").
		WithConfig(mockCfg).
		WithMetrics(metricsRegistry).
		Build()

	tracer := app.GetTracer()
	if tracer == nil {
		t.Fatal("GetTracer() returned nil; expected a usable no-op tracer")
	}

	if tracer.IsEnabled() {
		t.Error("expected IsEnabled() to be false when tracing config is disabled")
	}

	// Calling NewCollectorSpan on the no-op tracer must not panic and must
	// return a usable span whose End() / SetAttributes() / etc. are safe.
	ctx, span := func() (any, *struct{ ok bool }) {
		// Use deferred recovery so a panic doesn't blow up the whole test.
		var (
			c any
			s *struct{ ok bool }
		)

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("NewCollectorSpan panicked on disabled tracer: %v", r)
			}
		}()

		cs := tracer.NewCollectorSpan(nil, "test-collector", "test-op")
		if cs == nil {
			t.Fatal("NewCollectorSpan returned nil on disabled tracer")
		}

		cs.End()
		c = cs.Context()
		s = &struct{ ok bool }{ok: true}

		return c, s
	}()
	_ = ctx
	_ = span
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
