package profiling

import (
	"log/slog"

	"github.com/d0ugal/promexporter/v2/config"
	"github.com/grafana/pyroscope-go"
)

// Profiler wraps the Pyroscope profiler with configuration
type Profiler struct {
	config   *config.ProfilingConfig
	profiler *pyroscope.Profiler
}

// NewProfiler creates a new profiler instance with the given configuration
func NewProfiler(cfg *config.ProfilingConfig, version, commit string) (*Profiler, error) {
	slog.Debug("NewProfiler called",
		"enabled", cfg.IsEnabled(),
		"service_name", cfg.ServiceName,
		"server_address", cfg.ServerAddress,
	)

	if !cfg.IsEnabled() {
		slog.Debug("Profiling disabled, returning empty profiler")
		return &Profiler{}, nil
	}

	serviceName := cfg.ServiceName
	if serviceName == "" {
		serviceName = "promexporter-app"
	}

	serverAddress := cfg.ServerAddress
	if serverAddress == "" {
		serverAddress = "http://10.10.10.2:4040"
	}

	slog.Info("Initializing continuous profiling",
		"service_name", serviceName,
		"server_address", serverAddress)

	// Build tags map
	tags := make(map[string]string)
	if version != "" {
		tags["version"] = version
	}

	if commit != "" {
		tags["commit"] = commit
	}

	// Initialize Pyroscope with CPU and memory profiling
	profiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: serviceName,
		ServerAddress:   serverAddress,
		Logger:          pyroscope.StandardLogger,
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
		Tags: tags,
	})
	if err != nil {
		slog.Warn("Failed to initialize continuous profiling, continuing without profiling",
			"error", err,
			"server_address", serverAddress)
		// Return empty profiler rather than failing - profiling is optional
		return &Profiler{}, nil
	}

	slog.Info("Continuous profiling initialized successfully",
		"service_name", serviceName,
		"server_address", serverAddress)

	return &Profiler{
		config:   cfg,
		profiler: profiler,
	}, nil
}

// IsEnabled returns true if profiling is enabled
func (p *Profiler) IsEnabled() bool {
	return p.config != nil && p.config.IsEnabled()
}

// Stop gracefully stops the profiler
func (p *Profiler) Stop() {
	slog.Debug("Profiler stop called", "enabled", p.IsEnabled())

	if !p.IsEnabled() || p.profiler == nil {
		slog.Debug("Profiling disabled or profiler nil, skipping stop")
		return
	}

	// Pyroscope profiler doesn't require explicit stopping - it runs until the process exits
	// The profiler will automatically stop when the process terminates
	slog.Debug("Profiler marked for stop (will stop on process exit)")
}
