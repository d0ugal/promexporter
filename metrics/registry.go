package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Registry provides a wrapper around Prometheus registry with metric info tracking
type Registry struct {
	// The underlying Prometheus registry
	registry *prometheus.Registry

	// Version info metric (standard across all exporters)
	VersionInfo *prometheus.GaugeVec

	// Metric information for UI
	metricInfo []MetricInfo
}

// NewRegistry creates a new metrics registry
func NewRegistry() *Registry {
	registry := prometheus.NewRegistry()
	factory := promauto.With(registry)

	r := &Registry{
		registry: registry,
		VersionInfo: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "exporter_info",
				Help: "Information about the exporter",
			},
			[]string{"version", "commit", "build_date"},
		),
	}

	// Add version info metric to the UI info
	r.addMetricInfo("exporter_info", "Information about the exporter", []string{"version", "commit", "build_date"})

	return r
}

// NewRegistryWithCustom creates a new metrics registry with a custom registry
func NewRegistryWithCustom(registry *prometheus.Registry) *Registry {
	factory := promauto.With(registry)

	r := &Registry{
		registry: registry,
		VersionInfo: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "exporter_info",
				Help: "Information about the exporter",
			},
			[]string{"version", "commit", "build_date"},
		),
	}

	// Add version info metric to the UI info
	r.addMetricInfo("exporter_info", "Information about the exporter", []string{"version", "commit", "build_date"})

	return r
}

// addMetricInfo adds metric information to the registry
func (r *Registry) addMetricInfo(name, help string, labels []string) {
	r.metricInfo = append(r.metricInfo, MetricInfo{
		Name:         name,
		Help:         help,
		Labels:       labels,
		ExampleValue: "",
	})
}

// AddMetricInfo allows external packages to add metric information
func (r *Registry) AddMetricInfo(name, help string, labels []string) {
	r.addMetricInfo(name, help, labels)
}

// GetRegistry returns the underlying Prometheus registry
func (r *Registry) GetRegistry() *prometheus.Registry {
	return r.registry
}

// GetMetricsInfo returns information about all metrics for the UI
func (r *Registry) GetMetricsInfo() []MetricInfo {
	return r.metricInfo
}
