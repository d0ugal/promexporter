package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d0ugal/promexporter/config"
	"github.com/d0ugal/promexporter/logging"
	"github.com/d0ugal/promexporter/metrics"
	"github.com/d0ugal/promexporter/server"
	"github.com/d0ugal/promexporter/tracing"
	"github.com/d0ugal/promexporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

// ConfigInterface defines the interface that configuration types must implement
type ConfigInterface interface {
	GetDisplayConfig() map[string]interface{}
	GetLogging() *config.LoggingConfig
	GetServer() *config.ServerConfig
	GetTracing() *config.TracingConfig
}

// App represents the main application
type App struct {
	name        string
	config      ConfigInterface
	metrics     *metrics.Registry
	server      *server.Server
	collectors  []Collector
	versionInfo *VersionInfo
	tracer      *tracing.Tracer
}

// VersionInfo holds version information for the application
type VersionInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

// Collector interface for data collection
type Collector interface {
	Start(ctx context.Context)
	Stop()
}

// New creates a new application instance
func New(name string) *App {
	return &App{
		name: name,
	}
}

// WithConfig sets the configuration
func (a *App) WithConfig(cfg ConfigInterface) *App {
	a.config = cfg
	return a
}

// WithMetrics sets the metrics registry
func (a *App) WithMetrics(registry *metrics.Registry) *App {
	a.metrics = registry
	return a
}

// WithCollector adds a collector to the application
func (a *App) WithCollector(collector Collector) *App {
	a.collectors = append(a.collectors, collector)
	return a
}

// WithVersionInfo sets custom version information for the application
func (a *App) WithVersionInfo(version, commit, buildDate string) *App {
	a.versionInfo = &VersionInfo{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}

	return a
}

// GetTracer returns the tracer instance (may be nil if tracing is disabled)
func (a *App) GetTracer() *tracing.Tracer {
	return a.tracer
}

// Build finalizes the application setup
func (a *App) Build() *App {
	// Configure logging
	loggingConfig := a.config.GetLogging()
	logging.Configure(&logging.Config{
		Level:  loggingConfig.Level,
		Format: loggingConfig.Format,
	})

	// Initialize tracing
	tracingConfig := a.config.GetTracing()
	if tracingConfig.IsEnabled() {
		tracer, err := tracing.NewTracer(tracingConfig)
		if err != nil {
			slog.Error("Failed to initialize tracing", "error", err)
			// Continue without tracing rather than failing
		} else {
			a.tracer = tracer

			slog.Info("Tracing enabled", "service_name", tracingConfig.ServiceName)

			// Start runtime metrics collection (requires tracing to be enabled)
			// This automatically collects Go runtime metrics (GC, memory, goroutines, etc.)
			if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
				slog.Warn("Failed to start runtime metrics collection", "error", err)
				// Continue without runtime metrics rather than failing
			} else {
				slog.Info("Runtime metrics collection enabled")
			}
		}
	}

	// Set version info metric
	if a.versionInfo != nil {
		// Use custom version info if provided
		a.metrics.VersionInfo.With(prometheus.Labels{
			"version":    a.versionInfo.Version,
			"commit":     a.versionInfo.Commit,
			"build_date": a.versionInfo.BuildDate,
		}).Set(1)
	} else {
		// Fall back to default version info
		slog.Warn("No custom version info provided, falling back to build defaults")

		versionInfo := version.Get()
		a.metrics.VersionInfo.With(prometheus.Labels{
			"version":    versionInfo.Version,
			"commit":     versionInfo.Commit,
			"build_date": versionInfo.BuildDate,
		}).Set(1)
	}

	// Create server
	var serverVersionInfo *version.Info
	if a.versionInfo != nil {
		serverVersionInfo = &version.Info{
			Version:   a.versionInfo.Version,
			Commit:    a.versionInfo.Commit,
			BuildDate: a.versionInfo.BuildDate,
		}
	}

	a.server = server.New(a.config, a.metrics, a.name, serverVersionInfo, a.tracer)

	return a
}

// Run starts the application
func (a *App) Run() error {
	// Start collectors
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, collector := range a.collectors {
		collector.Start(ctx)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		slog.Info("Shutting down gracefully...")
		cancel()

		// Stop collectors
		for _, collector := range a.collectors {
			collector.Stop()
		}

		// Shutdown tracing
		if a.tracer != nil {
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()

			if err := a.tracer.Shutdown(shutdownCtx); err != nil {
				slog.Error("Failed to shutdown tracing gracefully", "error", err)
			}
		}

		// Shutdown server
		if err := a.server.Shutdown(); err != nil {
			slog.Error("Failed to shutdown server gracefully", "error", err)
		}
	}()

	// Start server
	if err := a.server.Start(); err != nil {
		slog.Error("Server failed", "error", err)
		return err
	}

	return nil
}
