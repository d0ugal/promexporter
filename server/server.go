package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/d0ugal/promexporter/config"
	"github.com/d0ugal/promexporter/metrics"
	"github.com/d0ugal/promexporter/version"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server handles HTTP requests and serves metrics
type Server struct {
	config  *config.BaseConfig
	metrics *metrics.Registry
	server  *http.Server
	router  *gin.Engine
	name    string
}

// New creates a new server instance
func New(cfg *config.BaseConfig, metricsRegistry *metrics.Registry, exporterName string) *Server {
	// Set Gin to release mode unless debug logging is enabled
	if cfg.Logging.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(customGinLogger(), gin.Recovery())

	server := &Server{
		config:  cfg,
		metrics: metricsRegistry,
		router:  router,
		name:    exporterName,
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

func (s *Server) setupRoutes() {
	// Root endpoint with HTML dashboard
	s.router.GET("/", s.handleRoot)

	// Metrics endpoint - use our custom registry
	s.router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(s.metrics.GetRegistry(), promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})))

	// Health endpoint
	s.router.GET("/health", s.handleHealth)
}

func (s *Server) handleRoot(c *gin.Context) {
	versionInfo := version.Get()
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
func (s *Server) getConfigData() map[string]interface{} {
	return map[string]interface{}{
		"Server Host": s.config.Server.Host,
		"Server Port": s.config.Server.Port,
		"Log Level":   s.config.Logging.Level,
		"Log Format":  s.config.Logging.Format,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

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
