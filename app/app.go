package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/d0ugal/promexporter/config"
	"github.com/d0ugal/promexporter/logging"
	"github.com/d0ugal/promexporter/metrics"
	"github.com/d0ugal/promexporter/server"
	"github.com/d0ugal/promexporter/version"
)

// App represents the main application
type App struct {
	name     string
	config   *config.BaseConfig
	metrics  *metrics.Registry
	server   *server.Server
	collectors []Collector
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
func (a *App) WithConfig(cfg *config.BaseConfig) *App {
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

// Build finalizes the application setup
func (a *App) Build() *App {
	// Configure logging
	logging.Configure(&logging.Config{
		Level:  a.config.Logging.Level,
		Format: a.config.Logging.Format,
	})

	// Set version info metric
	versionInfo := version.Get()
	a.metrics.VersionInfo.WithLabelValues(versionInfo.Version, versionInfo.Commit, versionInfo.BuildDate).Set(1)

	// Create server
	a.server = server.New(a.config, a.metrics, a.name)

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
