package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/d0ugal/promexporter/config"
	"github.com/d0ugal/promexporter/metrics"
	"github.com/d0ugal/promexporter/tracing"
	"github.com/d0ugal/promexporter/version"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ConfigInterface defines the interface that configuration types must implement
type ConfigInterface interface {
	GetDisplayConfig() map[string]interface{}
	GetLogging() *config.LoggingConfig
	GetServer() *config.ServerConfig
	GetTracing() *config.TracingConfig
}

// CustomConfigRenderer allows exporters to provide custom HTML fragments for specific config keys
type CustomConfigRenderer interface {
	RenderConfigHTML(key string, value interface{}) (string, bool)
}

// Server handles HTTP requests and serves metrics
type Server struct {
	config      ConfigInterface
	metrics     *metrics.Registry
	server      *http.Server
	router      *gin.Engine
	name        string
	versionInfo *version.Info
	tracer      *tracing.Tracer
}

// New creates a new server instance
func New(cfg ConfigInterface, metricsRegistry *metrics.Registry, exporterName string, customVersionInfo *version.Info, tracer *tracing.Tracer) *Server {
	// Set Gin to release mode unless debug logging is enabled
	loggingConfig := cfg.GetLogging()
	if loggingConfig.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	
	// Add tracing middleware if tracer is available
	if tracer != nil && tracer.IsEnabled() {
		router.Use(tracer.HTTPMiddleware())
	}
	
	router.Use(customGinLogger(), gin.Recovery())

	server := &Server{
		config:      cfg,
		metrics:     metricsRegistry,
		router:      router,
		name:        exporterName,
		versionInfo: customVersionInfo,
		tracer:      tracer,
	}

	server.setupRoutes()

	return server
}

// customGinLogger creates a custom Gin logger that uses slog
func customGinLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Use slog to log the request
		slog.Info("HTTP request",
			"method", param.Method,
			"path", param.Path,
			"status", param.StatusCode,
			"latency", param.Latency,
			"client_ip", param.ClientIP,
			"user_agent", param.Request.UserAgent(),
		)

		return "" // Return empty string since slog handles the output
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	serverConfig := s.config.GetServer()
	addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)

	s.server = &http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadHeaderTimeout: 30 * time.Second,
	}

	slog.Info("Starting exporter server",
		"name", s.name,
		"address", addr,
	)

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			slog.Error("Server shutdown error", "error", err)
			return err
		} else {
			slog.Info("Server shutdown gracefully")
		}
	}

	return nil
}

func (s *Server) setupRoutes() {
	// Root endpoint with HTML dashboard (optional)
	if s.config.GetServer().IsWebUIEnabled() {
		s.router.GET("/", s.handleRoot)
	}

	// Metrics endpoint - use our custom registry
	s.router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(s.metrics.GetRegistry(), promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})))

	// Health endpoint (optional)
	if s.config.GetServer().IsHealthEnabled() {
		s.router.GET("/health", s.handleHealth)
	}
}

func (s *Server) handleRoot(c *gin.Context) {
	var versionInfo *version.Info
	if s.versionInfo != nil {
		versionInfo = s.versionInfo
	} else {
		defaultVersion := version.Get()
		versionInfo = &defaultVersion
	}

	metricsInfo := s.metrics.GetMetricsInfo()

	// Convert metrics to template data
	metrics := make([]MetricData, 0, len(metricsInfo))
	for _, metric := range metricsInfo {
		metrics = append(metrics, MetricData{
			Name:         metric.Name,
			Help:         metric.Help,
			Labels:       metric.Labels,
			ExampleValue: metric.ExampleValue,
		})
	}

	data := TemplateData{
		ExporterName: s.name,
		Version:      versionInfo.Version,
		Commit:       versionInfo.Commit,
		BuildDate:    versionInfo.BuildDate,
		Status:       "ready",
		Config:       s.getConfigData(),
		Metrics:      metrics,
	}

	c.Header("Content-Type", "text/html")

	if err := mainTemplate.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func (s *Server) handleHealth(c *gin.Context) {
	versionInfo := version.Get()
	c.JSON(http.StatusOK, gin.H{
		"status":     "healthy",
		"timestamp":  time.Now().Unix(),
		"service":    s.name,
		"version":    versionInfo.Version,
		"commit":     versionInfo.Commit,
		"build_date": versionInfo.BuildDate,
	})
}

// getConfigData returns configuration data for the template
// Uses the BaseConfig's GetDisplayConfig method and adds sensitivity information
func (s *Server) getConfigData() map[string]interface{} {
	config := s.config.GetDisplayConfig()

	// Add sensitivity information and custom HTML to each config value
	for key, value := range config {
		// Check if the value implements SensitiveValue interface
		if sensitiveValue, ok := value.(interface{ IsSensitive() bool }); ok && sensitiveValue.IsSensitive() {
			// Wrap sensitive values with metadata
			config[key] = map[string]interface{}{
				"value":     value,
				"sensitive": true,
			}
		} else {
			// For non-sensitive values, preserve the original data structure
			// but add sensitivity metadata
			config[key] = map[string]interface{}{
				"value":     value,
				"sensitive": false,
			}
		}

		// Check if the config implements CustomConfigRenderer
		if renderer, ok := s.config.(CustomConfigRenderer); ok {
			if customHTML, hasCustom := renderer.RenderConfigHTML(key, value); hasCustom {
				// Add custom HTML fragment to the config value as template.HTML to prevent escaping
				// SECURITY NOTE: This bypasses HTML escaping. The RenderConfigHTML method should only
				// return trusted HTML content from the application's own configuration, not user input.
				if configMap, ok := config[key].(map[string]interface{}); ok {
					configMap["custom_html"] = template.HTML(customHTML) //nolint:gosec // Trusted HTML from config renderer
				}
			}
		}
	}

	return config
}
