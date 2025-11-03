package config

import (
	"fmt"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

// BaseConfig provides common configuration for all exporters
type BaseConfig struct {
	Server    ServerConfig    `yaml:"server"`
	Logging   LoggingConfig   `yaml:"logging"`
	Metrics   MetricsConfig   `yaml:"metrics"`
	Tracing   TracingConfig   `yaml:"tracing"`
	Profiling ProfilingConfig `yaml:"profiling"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	EnableWebUI  *bool  `yaml:"enable_web_ui,omitempty"` // Enable web UI (default: true)
	EnableHealth *bool  `yaml:"enable_health,omitempty"` // Enable health endpoint (default: true)
}

// IsWebUIEnabled returns true if web UI is enabled (defaults to true)
func (s *ServerConfig) IsWebUIEnabled() bool {
	if s.EnableWebUI == nil {
		return true // default to enabled
	}

	return *s.EnableWebUI
}

// IsHealthEnabled returns true if health endpoint is enabled (defaults to true)
func (s *ServerConfig) IsHealthEnabled() bool {
	if s.EnableHealth == nil {
		return true // default to enabled
	}

	return *s.EnableHealth
}

// IsEnabled returns true if tracing is enabled (defaults to false)
func (t *TracingConfig) IsEnabled() bool {
	if t.Enabled == nil {
		return false // default to disabled
	}

	return *t.Enabled
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"` // "json" or "text"
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Collection CollectionConfig `yaml:"collection"`
}

// CollectionConfig holds collection configuration
type CollectionConfig struct {
	DefaultInterval Duration `yaml:"default_interval"`
	// Track if the value was explicitly set
	DefaultIntervalSet bool `yaml:"-"`
}

// TracingConfig holds tracing configuration
type TracingConfig struct {
	Enabled     *bool             `yaml:"enabled,omitempty"` // Enable tracing (default: false)
	ServiceName string            `yaml:"service_name"`      // Service name for traces
	Endpoint    string            `yaml:"endpoint"`          // OTLP endpoint (default: "http://localhost:4318/v1/traces")
	Headers     map[string]string `yaml:"headers"`           // Additional headers for OTLP
}

// ProfilingConfig holds profiling configuration
type ProfilingConfig struct {
	Enabled       *bool  `yaml:"enabled,omitempty"` // Enable profiling (default: false)
	ServiceName   string `yaml:"service_name"`      // Service name for profiling
	ServerAddress string `yaml:"server_address"`    // Pyroscope server address (default: "http://10.10.10.2:4040")
}

// IsEnabled returns true if profiling is enabled (defaults to false)
func (p *ProfilingConfig) IsEnabled() bool {
	if p.Enabled == nil {
		return false // default to disabled
	}

	return *p.Enabled
}

// UnmarshalYAML implements custom unmarshaling to track if the value was set
func (c *CollectionConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Create a temporary struct to unmarshal into
	type tempCollectionConfig struct {
		DefaultInterval Duration `yaml:"default_interval"`
	}

	var temp tempCollectionConfig
	if err := unmarshal(&temp); err != nil {
		return err
	}

	c.DefaultInterval = temp.DefaultInterval
	c.DefaultIntervalSet = true

	return nil
}

// LoadConfig loads configuration from either a YAML file or environment variables
// If configFromEnv is true, it will load from environment variables only
func LoadConfig(path string, configFromEnv bool) (*BaseConfig, error) {
	if configFromEnv {
		return loadFromEnv()
	}

	return Load(path)
}

// Load loads configuration from a YAML file
func Load(path string) (*BaseConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config BaseConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	setDefaults(&config)

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv() (*BaseConfig, error) {
	config := &BaseConfig{}

	// Server configuration
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	} else {
		config.Server.Host = "0.0.0.0"
	}

	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := parseInt(portStr); err != nil {
			return nil, fmt.Errorf("invalid server port: %w", err)
		} else {
			config.Server.Port = port
		}
	} else {
		config.Server.Port = 8080
	}

	// Logging configuration
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Logging.Level = level
	} else {
		config.Logging.Level = "info"
	}

	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.Logging.Format = format
	} else {
		config.Logging.Format = "json"
	}

	// Metrics configuration
	if intervalStr := os.Getenv("METRICS_DEFAULT_INTERVAL"); intervalStr != "" {
		if interval, err := time.ParseDuration(intervalStr); err != nil {
			return nil, fmt.Errorf("invalid metrics default interval: %w", err)
		} else {
			config.Metrics.Collection.DefaultInterval = Duration{interval}
			config.Metrics.Collection.DefaultIntervalSet = true
		}
	} else {
		config.Metrics.Collection.DefaultInterval = Duration{time.Second * 30}
	}

	// Tracing configuration
	if enabledStr := os.Getenv("TRACING_ENABLED"); enabledStr != "" {
		if enabled, err := parseBool(enabledStr); err != nil {
			return nil, fmt.Errorf("invalid tracing enabled value: %w", err)
		} else {
			config.Tracing.Enabled = &enabled
		}
	}

	if serviceName := os.Getenv("TRACING_SERVICE_NAME"); serviceName != "" {
		config.Tracing.ServiceName = serviceName
	}

	if endpoint := os.Getenv("TRACING_ENDPOINT"); endpoint != "" {
		config.Tracing.Endpoint = endpoint
	}

	// Profiling configuration
	if enabledStr := os.Getenv("PROFILING_ENABLED"); enabledStr != "" {
		if enabled, err := parseBool(enabledStr); err != nil {
			return nil, fmt.Errorf("invalid profiling enabled value: %w", err)
		} else {
			config.Profiling.Enabled = &enabled
		}
	}

	if serviceName := os.Getenv("PROFILING_SERVICE_NAME"); serviceName != "" {
		config.Profiling.ServiceName = serviceName
	}

	if serverAddress := os.Getenv("PROFILING_SERVER_ADDRESS"); serverAddress != "" {
		config.Profiling.ServerAddress = serverAddress
	}

	// Set defaults for any missing values
	setDefaults(config)

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// setDefaults sets default values for configuration
func setDefaults(config *BaseConfig) {
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}

	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}

	// Set default values for new options (only if not explicitly set in YAML)
	// Note: bool fields default to false, so we need to check if they were explicitly set
	// For now, we'll assume they default to true unless explicitly set to false

	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}

	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}

	if !config.Metrics.Collection.DefaultIntervalSet {
		config.Metrics.Collection.DefaultInterval = Duration{time.Second * 30}
	}

	// Tracing defaults
	if config.Tracing.ServiceName == "" {
		config.Tracing.ServiceName = "promexporter"
	}

	if config.Tracing.Endpoint == "" {
		config.Tracing.Endpoint = ""
	}

	if config.Tracing.Headers == nil {
		config.Tracing.Headers = make(map[string]string)
	}
}

// parseInt parses a string to int
func parseInt(s string) (int, error) {
	var i int

	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return 0, err
	}
	// Check if there are any remaining characters (like decimal points)
	if len(fmt.Sprintf("%d", i)) != len(s) {
		return 0, fmt.Errorf("invalid integer format: %s", s)
	}

	return i, nil
}

// parseBool parses a string to bool
func parseBool(s string) (bool, error) {
	switch s {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}

// Validate performs comprehensive validation of the configuration
func (c *BaseConfig) Validate() error {
	// Validate server configuration
	if err := c.validateServerConfig(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	// Validate logging configuration
	if err := c.validateLoggingConfig(); err != nil {
		return fmt.Errorf("logging config: %w", err)
	}

	// Validate metrics configuration
	if err := c.validateMetricsConfig(); err != nil {
		return fmt.Errorf("metrics config: %w", err)
	}

	// Validate tracing configuration
	if err := c.validateTracingConfig(); err != nil {
		return fmt.Errorf("tracing config: %w", err)
	}

	return nil
}

// GetDefaultInterval returns the default collection interval
func (c *BaseConfig) GetDefaultInterval() int {
	return c.Metrics.Collection.DefaultInterval.Seconds()
}

// GetDisplayConfig returns configuration data safe for display
// This method can be overridden by exporters to include their own configuration
func (c *BaseConfig) GetDisplayConfig() map[string]interface{} {
	config := map[string]interface{}{
		"Server Host":    c.Server.Host,
		"Server Port":    c.Server.Port,
		"Web UI Enabled": c.Server.IsWebUIEnabled(),
		"Health Enabled": c.Server.IsHealthEnabled(),
		"Log Level":      c.Logging.Level,
		"Log Format":     c.Logging.Format,
	}

	// Add tracing info if enabled
	if c.Tracing.IsEnabled() {
		config["Tracing Enabled"] = true
		config["Tracing Service Name"] = c.Tracing.ServiceName
		config["Tracing Endpoint"] = c.Tracing.Endpoint
	} else {
		config["Tracing Enabled"] = false
	}

	return config
}

// GetLogging returns the logging configuration
func (c *BaseConfig) GetLogging() *LoggingConfig {
	return &c.Logging
}

// GetServer returns the server configuration
func (c *BaseConfig) GetServer() *ServerConfig {
	return &c.Server
}

// GetProfiling returns the profiling configuration
func (c *BaseConfig) GetProfiling() *ProfilingConfig {
	return &c.Profiling
}

// GetTracing returns the tracing configuration
func (c *BaseConfig) GetTracing() *TracingConfig {
	return &c.Tracing
}

func (c *BaseConfig) validateServerConfig() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Server.Port)
	}

	return nil
}

func (c *BaseConfig) validateLoggingConfig() error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("invalid logging level: %s", c.Logging.Level)
	}

	validFormats := map[string]bool{
		"json": true,
		"text": true,
	}
	if !validFormats[c.Logging.Format] {
		return fmt.Errorf("invalid logging format: %s", c.Logging.Format)
	}

	return nil
}

func (c *BaseConfig) validateMetricsConfig() error {
	if c.Metrics.Collection.DefaultInterval.Seconds() < 1 {
		return fmt.Errorf("default interval must be at least 1 second, got %d", c.Metrics.Collection.DefaultInterval.Seconds())
	}

	if c.Metrics.Collection.DefaultInterval.Seconds() > 86400 {
		return fmt.Errorf("default interval must be at most 86400 seconds (24 hours), got %d", c.Metrics.Collection.DefaultInterval.Seconds())
	}

	return nil
}

func (c *BaseConfig) validateTracingConfig() error {
	// Only validate if tracing is enabled
	if !c.Tracing.IsEnabled() {
		return nil
	}

	if c.Tracing.ServiceName == "" {
		return fmt.Errorf("service name is required when tracing is enabled")
	}

	if c.Tracing.Endpoint == "" {
		return fmt.Errorf("tracing endpoint must be configured when tracing is enabled")
	}

	return nil
}
