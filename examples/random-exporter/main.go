package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d0ugal/promexporter/app"
	"github.com/d0ugal/promexporter/config"
	"github.com/d0ugal/promexporter/logging"
	promexporter_metrics "github.com/d0ugal/promexporter/metrics"
	"github.com/d0ugal/promexporter/tracing"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
)

// RandomExporterConfig extends the base configuration
type RandomExporterConfig struct {
	config.BaseConfig
	Random RandomConfig `yaml:"random"`
}

type RandomConfig struct {
	// Collection interval for random metrics
	CollectionInterval config.Duration `yaml:"collection_interval"`
	// Number of random metrics to generate
	MetricCount int `yaml:"metric_count"`
	// Enable random errors for testing error tracing
	EnableRandomErrors bool `yaml:"enable_random_errors"`
	// Error probability (0.0 to 1.0)
	ErrorProbability float64 `yaml:"error_probability"`
}

// RandomCollector implements the Collector interface
type RandomCollector struct {
	config  *RandomExporterConfig
	metrics *RandomMetrics
	app     *app.App
}

// RandomMetrics holds all the random metrics
type RandomMetrics struct {
	// Counter metrics
	RandomCounter     *prometheus.CounterVec
	RandomCounterRate *prometheus.CounterVec

	// Gauge metrics
	RandomGauge       *prometheus.GaugeVec
	RandomTemperature *prometheus.GaugeVec
	RandomMemory      *prometheus.GaugeVec

	// Histogram metrics
	RandomLatency      *prometheus.HistogramVec
	RandomResponseTime *prometheus.HistogramVec

	// Summary metrics
	RandomDuration       *prometheus.SummaryVec
	RandomProcessingTime *prometheus.SummaryVec

	// Info metrics
	RandomInfo *prometheus.GaugeVec
}

// NewRandomMetrics creates a new metrics registry
func NewRandomMetrics(registry *promexporter_metrics.Registry) *RandomMetrics {
	return &RandomMetrics{
		// Counter metrics
		RandomCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "random_counter_total",
				Help: "A random counter that increments over time",
			},
			[]string{"service", "region"},
		),
		RandomCounterRate: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "random_counter_rate_total",
				Help: "A random counter with varying increment rates",
			},
			[]string{"service", "rate_type"},
		),

		// Gauge metrics
		RandomGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "random_gauge",
				Help: "A random gauge with fluctuating values",
			},
			[]string{"instance", "type"},
		),
		RandomTemperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "random_temperature_celsius",
				Help: "Simulated temperature readings",
			},
			[]string{"sensor", "location"},
		),
		RandomMemory: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "random_memory_usage_bytes",
				Help: "Simulated memory usage",
			},
			[]string{"process", "type"},
		),

		// Histogram metrics
		RandomLatency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "random_latency_seconds",
				Help:    "Simulated request latency",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "endpoint"},
		),
		RandomResponseTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "random_response_time_seconds",
				Help:    "Simulated response times",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"method", "status"},
		),

		// Summary metrics
		RandomDuration: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "random_duration_seconds",
				Help:       "Simulated operation duration",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"operation", "priority"},
		),
		RandomProcessingTime: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "random_processing_time_seconds",
				Help:       "Simulated processing times",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"task", "worker"},
		),

		// Info metrics
		RandomInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "random_info",
				Help: "Information about the random exporter",
			},
			[]string{"version", "build_date", "go_version"},
		),
	}
}

// RegisterMetrics registers all metrics with the registry
func (rm *RandomMetrics) RegisterMetrics(registry *promexporter_metrics.Registry) {
	promRegistry := registry.GetRegistry()
	promRegistry.MustRegister(
		rm.RandomCounter,
		rm.RandomCounterRate,
		rm.RandomGauge,
		rm.RandomTemperature,
		rm.RandomMemory,
		rm.RandomLatency,
		rm.RandomResponseTime,
		rm.RandomDuration,
		rm.RandomProcessingTime,
		rm.RandomInfo,
	)
}

// NewRandomCollector creates a new random collector
func NewRandomCollector(cfg *RandomExporterConfig, metrics *RandomMetrics, app *app.App) *RandomCollector {
	return &RandomCollector{
		config:  cfg,
		metrics: metrics,
		app:     app,
	}
}

// Start implements the Collector interface
func (rc *RandomCollector) Start(ctx context.Context) {
	// Start the collection loop
	go func() {
		ticker := time.NewTicker(rc.config.Random.CollectionInterval.Duration)
		defer ticker.Stop()
		
		slog.Info("Starting metric collection loop", 
			"interval", rc.config.Random.CollectionInterval.Duration,
			"metric_count", rc.config.Random.MetricCount,
			"enable_errors", rc.config.Random.EnableRandomErrors,
			"error_probability", rc.config.Random.ErrorProbability,
		)
		
		for {
			select {
			case <-ctx.Done():
				slog.Info("Stopping metric collection loop")
				return
			case <-ticker.C:
				if err := rc.Collect(ctx); err != nil {
					slog.Error("Failed to collect metrics", "error", err)
				}
			}
		}
	}()
}

// Stop implements the Collector interface
func (rc *RandomCollector) Stop() {
	slog.Info("Stopping random collector")
}

// Collect implements the Collector interface
func (rc *RandomCollector) Collect(ctx context.Context) error {
	// Get tracer for this collection cycle
	tracer := rc.app.GetTracer()
	var collectorSpan *tracing.CollectorSpan
	var spanCtx context.Context

	if tracer != nil && tracer.IsEnabled() {
		collectorSpan = tracer.NewCollectorSpan(ctx, "random-collector", "collect-metrics")
		spanCtx = collectorSpan.Context()
		defer collectorSpan.End()
	} else {
		spanCtx = ctx
	}

	// Add event for collection start
	if collectorSpan != nil {
		collectorSpan.AddEvent("collection_started",
			attribute.String("metric_count", fmt.Sprintf("%d", rc.config.Random.MetricCount)))
	}

	// Check for random errors
	if rc.config.Random.EnableRandomErrors && rand.Float64() < rc.config.Random.ErrorProbability {
		err := fmt.Errorf("random error occurred during collection")
		if collectorSpan != nil {
			collectorSpan.RecordError(err, attribute.String("error_type", "random"))
		}
		return err
	}

	// Generate random metrics
	if err := rc.generateCounterMetrics(spanCtx, collectorSpan); err != nil {
		return err
	}

	if err := rc.generateGaugeMetrics(spanCtx, collectorSpan); err != nil {
		return err
	}

	if err := rc.generateHistogramMetrics(spanCtx, collectorSpan); err != nil {
		return err
	}

	if err := rc.generateSummaryMetrics(spanCtx, collectorSpan); err != nil {
		return err
	}

	if err := rc.generateInfoMetrics(spanCtx, collectorSpan); err != nil {
		return err
	}

	// Add event for successful collection
	if collectorSpan != nil {
		collectorSpan.AddEvent("collection_completed",
			attribute.String("metrics_generated", fmt.Sprintf("%d", rc.config.Random.MetricCount)))
	}

	return nil
}

// generateCounterMetrics generates random counter metrics
func (rc *RandomCollector) generateCounterMetrics(ctx context.Context, span *tracing.CollectorSpan) error {
	// Create a sub-span for counter metrics
	tracer := rc.app.GetTracer()
	var counterSpan *tracing.CollectorSpan

	if tracer != nil && tracer.IsEnabled() {
		counterSpan = tracer.NewCollectorSpan(ctx, "random-collector", "generate-counters")
		defer counterSpan.End()
	}

	services := []string{"api", "database", "cache", "queue"}
	regions := []string{"us-east-1", "us-west-2", "eu-west-1", "ap-southeast-1"}
	rateTypes := []string{"low", "medium", "high", "burst"}

	for i := 0; i < rc.config.Random.MetricCount/4; i++ {
		service := services[rand.Intn(len(services))]
		region := regions[rand.Intn(len(regions))]
		rateType := rateTypes[rand.Intn(len(rateTypes))]

		// Increment counters with random values
		increment := rand.Float64() * 10
		rc.metrics.RandomCounter.With(prometheus.Labels{
			"service": service,
			"region":  region,
		}).Add(increment)

		rateIncrement := rand.Float64() * 5
		rc.metrics.RandomCounterRate.With(prometheus.Labels{
			"service":  service,
			"rate_type": rateType,
		}).Add(rateIncrement)
	}

	if counterSpan != nil {
		counterSpan.AddEvent("counters_generated",
			attribute.Int("count", rc.config.Random.MetricCount/4))
	}

	return nil
}

// generateGaugeMetrics generates random gauge metrics
func (rc *RandomCollector) generateGaugeMetrics(ctx context.Context, span *tracing.CollectorSpan) error {
	tracer := rc.app.GetTracer()
	var gaugeSpan *tracing.CollectorSpan

	if tracer != nil && tracer.IsEnabled() {
		gaugeSpan = tracer.NewCollectorSpan(ctx, "random-collector", "generate-gauges")
		defer gaugeSpan.End()
	}

	instances := []string{"instance-1", "instance-2", "instance-3"}
	types := []string{"cpu", "memory", "disk"}
	sensors := []string{"temp-01", "temp-02", "temp-03"}
	locations := []string{"server-room", "datacenter", "office"}
	processes := []string{"nginx", "postgres", "redis"}
	memoryTypes := []string{"heap", "stack", "cache"}

	for i := 0; i < rc.config.Random.MetricCount/4; i++ {
		instance := instances[rand.Intn(len(instances))]
		metricType := types[rand.Intn(len(types))]
		sensor := sensors[rand.Intn(len(sensors))]
		location := locations[rand.Intn(len(locations))]
		process := processes[rand.Intn(len(processes))]
		memType := memoryTypes[rand.Intn(len(memoryTypes))]

		// Set gauge values
		rc.metrics.RandomGauge.With(prometheus.Labels{
			"instance": instance,
			"type":     metricType,
		}).Set(rand.Float64() * 100)
		rc.metrics.RandomTemperature.With(prometheus.Labels{
			"sensor":   sensor,
			"location": location,
		}).Set(rand.Float64()*30 + 20) // 20-50°C
		rc.metrics.RandomMemory.With(prometheus.Labels{
			"process": process,
			"type":    memType,
		}).Set(rand.Float64() * 1024 * 1024 * 1024) // Up to 1GB
	}

	if gaugeSpan != nil {
		gaugeSpan.AddEvent("gauges_generated",
			attribute.Int("count", rc.config.Random.MetricCount/4))
	}

	return nil
}

// generateHistogramMetrics generates random histogram metrics
func (rc *RandomCollector) generateHistogramMetrics(ctx context.Context, span *tracing.CollectorSpan) error {
	tracer := rc.app.GetTracer()
	var histogramSpan *tracing.CollectorSpan
	
	if tracer != nil && tracer.IsEnabled() {
		histogramSpan = tracer.NewCollectorSpan(ctx, "random-collector", "generate-histograms")
		defer histogramSpan.End()
	}

	services := []string{"api", "database", "cache"}
	endpoints := []string{"/api/users", "/api/orders", "/api/products"}
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	statuses := []string{"200", "404", "500"}

	for i := 0; i < rc.config.Random.MetricCount/4; i++ {
		service := services[rand.Intn(len(services))]
		endpoint := endpoints[rand.Intn(len(endpoints))]
		method := methods[rand.Intn(len(methods))]
		status := statuses[rand.Intn(len(statuses))]

		// Generate random latency values
		latency := rand.Float64() * 2.0 // 0-2 seconds
		rc.metrics.RandomLatency.With(prometheus.Labels{
			"service":  service,
			"endpoint": endpoint,
		}).Observe(latency)

		responseTime := rand.Float64() * 1.0 // 0-1 seconds
		rc.metrics.RandomResponseTime.With(prometheus.Labels{
			"method": method,
			"status": status,
		}).Observe(responseTime)
	}

	if histogramSpan != nil {
		histogramSpan.AddEvent("histograms_generated",
			attribute.Int("count", rc.config.Random.MetricCount/4))
	}

	return nil
}

// generateSummaryMetrics generates random summary metrics
func (rc *RandomCollector) generateSummaryMetrics(ctx context.Context, span *tracing.CollectorSpan) error {
	tracer := rc.app.GetTracer()
	var summarySpan *tracing.CollectorSpan
	
	if tracer != nil && tracer.IsEnabled() {
		summarySpan = tracer.NewCollectorSpan(ctx, "random-collector", "generate-summaries")
		defer summarySpan.End()
	}

	operations := []string{"process", "validate", "transform", "save"}
	priorities := []string{"low", "medium", "high", "critical"}
	tasks := []string{"data-processing", "image-resize", "file-upload", "email-send"}
	workers := []string{"worker-1", "worker-2", "worker-3"}

	for i := 0; i < rc.config.Random.MetricCount/4; i++ {
		operation := operations[rand.Intn(len(operations))]
		priority := priorities[rand.Intn(len(priorities))]
		task := tasks[rand.Intn(len(tasks))]
		worker := workers[rand.Intn(len(workers))]

		// Generate random duration values
		duration := rand.Float64() * 5.0 // 0-5 seconds
		rc.metrics.RandomDuration.With(prometheus.Labels{
			"operation": operation,
			"priority":  priority,
		}).Observe(duration)

		processingTime := rand.Float64() * 3.0 // 0-3 seconds
		rc.metrics.RandomProcessingTime.With(prometheus.Labels{
			"task":   task,
			"worker": worker,
		}).Observe(processingTime)
	}

	if summarySpan != nil {
		summarySpan.AddEvent("summaries_generated",
			attribute.Int("count", rc.config.Random.MetricCount/4))
	}

	return nil
}

// generateInfoMetrics generates info metrics
func (rc *RandomCollector) generateInfoMetrics(ctx context.Context, span *tracing.CollectorSpan) error {
	tracer := rc.app.GetTracer()
	var infoSpan *tracing.CollectorSpan
	
	if tracer != nil && tracer.IsEnabled() {
		infoSpan = tracer.NewCollectorSpan(ctx, "random-collector", "generate-info")
		defer infoSpan.End()
	}

	// Set info metrics (these don't change often)
	rc.metrics.RandomInfo.With(prometheus.Labels{
		"version":    "1.0.0",
		"build_date": "2024-01-01",
		"go_version": "1.21",
	}).Set(1)

	if infoSpan != nil {
		infoSpan.AddEvent("info_metrics_set")
	}

	return nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv() (*RandomExporterConfig, error) {
	cfg := &RandomExporterConfig{}
	
	// Load base configuration from environment
	baseConfig := &cfg.BaseConfig

	// Server configuration
	if host := os.Getenv("SERVER_HOST"); host != "" {
		baseConfig.Server.Host = host
	} else {
		baseConfig.Server.Host = "0.0.0.0"
	}

	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := parseInt(portStr); err != nil {
			return nil, fmt.Errorf("invalid server port: %w", err)
		} else {
			baseConfig.Server.Port = port
		}
	} else {
		baseConfig.Server.Port = 8080
	}

	// Logging configuration
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		baseConfig.Logging.Level = level
	} else {
		baseConfig.Logging.Level = "info"
	}

	if format := os.Getenv("LOG_FORMAT"); format != "" {
		baseConfig.Logging.Format = format
	} else {
		baseConfig.Logging.Format = "json"
	}

	// Tracing configuration
	if enabledStr := os.Getenv("TRACING_ENABLED"); enabledStr != "" {
		enabled := enabledStr == "true"
		baseConfig.Tracing.Enabled = &enabled
	}

	if serviceName := os.Getenv("TRACING_SERVICE_NAME"); serviceName != "" {
		baseConfig.Tracing.ServiceName = serviceName
	}

	if endpoint := os.Getenv("TRACING_ENDPOINT"); endpoint != "" {
		baseConfig.Tracing.Endpoint = endpoint
	}

	// Random configuration
	if intervalStr := os.Getenv("RANDOM_COLLECTION_INTERVAL"); intervalStr != "" {
		if interval, err := time.ParseDuration(intervalStr); err != nil {
			return nil, fmt.Errorf("invalid collection interval: %w", err)
		} else {
			cfg.Random.CollectionInterval = config.Duration{Duration: interval}
		}
	} else {
		cfg.Random.CollectionInterval = config.Duration{Duration: time.Second * 10}
	}
	
	if countStr := os.Getenv("RANDOM_METRIC_COUNT"); countStr != "" {
		if count, err := parseInt(countStr); err != nil {
			return nil, fmt.Errorf("invalid metric count: %w", err)
		} else {
			cfg.Random.MetricCount = count
		}
	} else {
		cfg.Random.MetricCount = 20
	}
	
	if enableErrorsStr := os.Getenv("RANDOM_ENABLE_ERRORS"); enableErrorsStr != "" {
		cfg.Random.EnableRandomErrors = enableErrorsStr == "true"
	} else {
		cfg.Random.EnableRandomErrors = false
	}
	
	if probStr := os.Getenv("RANDOM_ERROR_PROBABILITY"); probStr != "" {
		if prob, err := parseFloat(probStr); err != nil {
			return nil, fmt.Errorf("invalid error probability: %w", err)
		} else {
			cfg.Random.ErrorProbability = prob
		}
	} else {
		cfg.Random.ErrorProbability = 0.1
	}

	cfg.BaseConfig = *baseConfig
	return cfg, nil
}

// parseInt parses a string to int
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// parseFloat parses a string to float64
func parseFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	return result, err
}

func main() {
	// Parse command line flags
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information")

	var (
		configPath    string
		configFromEnv bool
		showConfig    bool
	)

	flag.StringVar(&configPath, "config", "config.yaml", "Path to configuration file")
	flag.BoolVar(&configFromEnv, "config-from-env", false, "Load configuration from environment variables only")
	flag.BoolVar(&showConfig, "show-config", false, "Show loaded configuration and exit")
	flag.Parse()

	// Show version if requested
	if showVersion {
		fmt.Printf("random-exporter v1.0.0\n")
		fmt.Printf("Build Date: %s\n", time.Now().Format("2006-01-02"))
		fmt.Printf("Go Version: %s\n", "1.21")
		os.Exit(0)
	}

	// Load configuration
	var cfg *RandomExporterConfig
	var err error

	if configFromEnv || os.Getenv("CONFIG_FROM_ENV") == "true" {
		cfg, err = loadFromEnv()
	} else {
		// For this example, we'll use environment-based config
		cfg, err = loadFromEnv()
	}

	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Show configuration if requested
	if showConfig {
		displayConfig := cfg.BaseConfig.GetDisplayConfig()
		fmt.Printf("Configuration:\n")
		for key, value := range displayConfig {
			fmt.Printf("  %s: %v\n", key, value)
		}
		os.Exit(0)
	}

	// Configure logging using promexporter
	logging.Configure(&logging.Config{
		Level:  cfg.BaseConfig.Logging.Level,
		Format: cfg.BaseConfig.Logging.Format,
	})

	slog.Info("Starting Random Exporter",
		"version", "1.0.0",
		"config_from_env", configFromEnv || os.Getenv("CONFIG_FROM_ENV") == "true",
		"tracing_enabled", cfg.BaseConfig.Tracing.IsEnabled(),
	)

	// Initialize metrics registry using promexporter
	metricsRegistry := promexporter_metrics.NewRegistry("random_exporter_info")

	// Create random metrics
	randomMetrics := NewRandomMetrics(metricsRegistry)
	randomMetrics.RegisterMetrics(metricsRegistry)

	// Create and run application using promexporter
	application := app.New("Random Exporter").
		WithConfig(&cfg.BaseConfig).
		WithMetrics(metricsRegistry).
		WithVersionInfo("1.0.0", "example-commit", time.Now().Format("2006-01-02"))

	// Create collector with app reference for tracing
	randomCollector := NewRandomCollector(cfg, randomMetrics, application)
	application.WithCollector(randomCollector)

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		slog.Info("Received shutdown signal", "signal", sig.String())
		cancel()
	}()

	// Start the application in a goroutine
	appErr := make(chan error, 1)
	go func() {
		appErr <- application.Build().Run()
	}()

	// Wait for application to finish or context cancellation
	select {
	case err := <-appErr:
		if err != nil {
			slog.Error("Application error", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		slog.Info("Shutting down gracefully")
	}
}
