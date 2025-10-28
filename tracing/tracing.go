package tracing

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/d0ugal/promexporter/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Tracer wraps the OpenTelemetry tracer with additional utilities
type Tracer struct {
	tracer trace.Tracer
	config *config.TracingConfig
}

// NewTracer creates a new tracer instance with the given configuration
func NewTracer(cfg *config.TracingConfig) (*Tracer, error) {
	slog.Debug("NewTracer called",
		"enabled", cfg.IsEnabled(),
		"service_name", cfg.ServiceName,
		"endpoint", cfg.Endpoint,
		"headers", cfg.Headers,
	)

	if !cfg.IsEnabled() {
		slog.Debug("Tracing disabled, returning empty tracer")
		return &Tracer{}, nil
	}

	slog.Debug("Creating OTLP HTTP exporter", "endpoint", cfg.Endpoint)

	// Parse and validate the endpoint URL
	endpointURL, err := url.Parse(cfg.Endpoint)
	if err != nil {
		slog.Error("Invalid endpoint URL", "error", err, "endpoint", cfg.Endpoint)
		return nil, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	slog.Debug("Parsed endpoint URL",
		"scheme", endpointURL.Scheme,
		"host", endpointURL.Host,
		"path", endpointURL.Path,
	)

	// OTLP HTTP exporter expects just the host:port, not the full URL
	// We need to construct the proper endpoint format
	slog.Debug("Creating OTLP exporter with parsed components",
		"host", endpointURL.Host,
		"path", endpointURL.Path,
		"scheme", endpointURL.Scheme,
	)

	var exporter sdktrace.SpanExporter

	if endpointURL.Scheme == "https" {
		exporter, err = otlptracehttp.New(
			context.Background(),
			otlptracehttp.WithEndpoint(endpointURL.Host),
			otlptracehttp.WithURLPath(endpointURL.Path),
			otlptracehttp.WithHeaders(cfg.Headers),
		)
	} else {
		// For HTTP endpoints, use WithInsecure
		exporter, err = otlptracehttp.New(
			context.Background(),
			otlptracehttp.WithEndpoint(endpointURL.Host),
			otlptracehttp.WithURLPath(endpointURL.Path),
			otlptracehttp.WithHeaders(cfg.Headers),
			otlptracehttp.WithInsecure(),
		)
	}

	if err != nil {
		slog.Error("Failed to create OTLP exporter", "error", err, "endpoint", cfg.Endpoint)
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	slog.Debug("OTLP HTTP exporter created successfully")

	// Create resource with service information
	slog.Debug("Creating resource", "service_name", cfg.ServiceName)

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("service.version", "1.0.0"), // Could be made configurable
		),
	)
	if err != nil {
		slog.Error("Failed to create resource", "error", err)
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	slog.Debug("Resource created successfully")

	// Create trace provider
	slog.Debug("Creating trace provider")

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Could be made configurable
	)

	slog.Debug("Trace provider created successfully")

	// Set global trace provider
	slog.Debug("Setting global trace provider")
	otel.SetTracerProvider(tp)

	// Set global text map propagator
	slog.Debug("Setting global text map propagator")
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create tracer
	slog.Debug("Creating tracer", "service_name", cfg.ServiceName)
	tracer := tp.Tracer(cfg.ServiceName)

	slog.Info("Tracing initialized",
		"service_name", cfg.ServiceName,
		"endpoint", cfg.Endpoint,
	)
	slog.Debug("Tracer setup completed successfully")

	return &Tracer{
		tracer: tracer,
		config: cfg,
	}, nil
}

// IsEnabled returns true if tracing is enabled
func (t *Tracer) IsEnabled() bool {
	return t.config != nil && t.config.IsEnabled()
}

// StartSpan creates a new span with the given name and options
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if !t.IsEnabled() {
		return ctx, trace.SpanFromContext(ctx)
	}

	return t.tracer.Start(ctx, name, opts...)
}

// StartSpanWithAttributes creates a new span with attributes
func (t *Tracer) StartSpanWithAttributes(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	if !t.IsEnabled() {
		return ctx, trace.SpanFromContext(ctx)
	}

	return t.tracer.Start(ctx, name, trace.WithAttributes(attrs...))
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent(name, trace.WithAttributes(attrs...))
	}
}

// SetSpanAttributes sets attributes on the current span
func SetSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(attrs...)
	}
}

// RecordSpanError records an error on the current span
func RecordSpanError(ctx context.Context, err error, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.RecordError(err, trace.WithAttributes(attrs...))
	}
}

// HTTPMiddleware creates a Gin middleware for HTTP request tracing
func (t *Tracer) HTTPMiddleware() func(c *gin.Context) {
	slog.Debug("HTTPMiddleware called", "enabled", t.IsEnabled())

	if !t.IsEnabled() {
		slog.Debug("Tracing disabled, returning no-op middleware")

		return func(c *gin.Context) {
			c.Next()
		}
	}

	slog.Debug("Creating HTTP tracing middleware")

	return func(c *gin.Context) {
		slog.Debug("Processing HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"url", c.Request.URL.String(),
		)

		// Extract trace context from headers
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Start span for the HTTP request
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		ctx, span := t.StartSpanWithAttributes(ctx, spanName,
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.user_agent", c.Request.UserAgent()),
		)

		slog.Debug("HTTP span created",
			"span_name", spanName,
			"span_id", span.SpanContext().SpanID().String(),
			"trace_id", span.SpanContext().TraceID().String(),
		)

		// Store context in Gin context
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Record response information
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)

		slog.Debug("HTTP request completed",
			"status_code", c.Writer.Status(),
			"span_id", span.SpanContext().SpanID().String(),
			"trace_id", span.SpanContext().TraceID().String(),
		)

		// Record error if status indicates error
		if c.Writer.Status() >= 400 {
			slog.Debug("Recording HTTP error", "status_code", c.Writer.Status())
			span.RecordError(fmt.Errorf("HTTP %d", c.Writer.Status()))
		}

		span.End()
		slog.Debug("HTTP span ended")
	}
}

// Shutdown gracefully shuts down the tracer
func (t *Tracer) Shutdown(ctx context.Context) error {
	slog.Debug("Tracer shutdown called", "enabled", t.IsEnabled())

	if !t.IsEnabled() {
		slog.Debug("Tracing disabled, skipping shutdown")
		return nil
	}

	tp := otel.GetTracerProvider()
	if sdkTp, ok := tp.(*sdktrace.TracerProvider); ok {
		slog.Debug("Shutting down SDK trace provider")

		err := sdkTp.Shutdown(ctx)
		if err != nil {
			slog.Error("Error during tracer shutdown", "error", err)
		} else {
			slog.Debug("Tracer shutdown completed successfully")
		}

		return err
	}

	slog.Debug("No SDK trace provider found, skipping shutdown")

	return nil
}

// CollectorSpan wraps common collector operations with tracing
type CollectorSpan struct {
	span trace.Span
	ctx  context.Context
}

// NewCollectorSpan creates a new collector span
func (t *Tracer) NewCollectorSpan(ctx context.Context, collectorName, operation string) *CollectorSpan {
	slog.Debug("NewCollectorSpan called",
		"enabled", t.IsEnabled(),
		"collector_name", collectorName,
		"operation", operation,
	)

	if !t.IsEnabled() {
		slog.Debug("Tracing disabled, returning empty collector span")
		return &CollectorSpan{ctx: ctx}
	}

	slog.Debug("Creating collector span",
		"collector_name", collectorName,
		"operation", operation,
	)

	ctx, span := t.StartSpanWithAttributes(ctx, operation,
		attribute.String("service.name", t.config.ServiceName),
		attribute.String("collector.name", collectorName),
		attribute.String("collector.operation", operation),
	)

	slog.Debug("Collector span created successfully",
		"span_id", span.SpanContext().SpanID().String(),
		"trace_id", span.SpanContext().TraceID().String(),
	)

	return &CollectorSpan{
		span: span,
		ctx:  ctx,
	}
}

// Context returns the context with the span
func (cs *CollectorSpan) Context() context.Context {
	return cs.ctx
}

// End ends the span
func (cs *CollectorSpan) End() {
	if cs.span != nil && cs.span.IsRecording() {
		slog.Debug("Ending collector span",
			"span_id", cs.span.SpanContext().SpanID().String(),
			"trace_id", cs.span.SpanContext().TraceID().String(),
		)
		cs.span.End()
		slog.Debug("Collector span ended successfully")
	} else {
		slog.Debug("Skipping span end - span is nil or not recording")
	}
}

// RecordError records an error on the span
func (cs *CollectorSpan) RecordError(err error, attrs ...attribute.KeyValue) {
	if cs.span != nil && cs.span.IsRecording() {
		cs.span.RecordError(err, trace.WithAttributes(attrs...))
	}
}

// SetAttributes sets attributes on the span
func (cs *CollectorSpan) SetAttributes(attrs ...attribute.KeyValue) {
	if cs.span != nil && cs.span.IsRecording() {
		cs.span.SetAttributes(attrs...)
	}
}

// AddEvent adds an event to the span
func (cs *CollectorSpan) AddEvent(name string, attrs ...attribute.KeyValue) {
	if cs.span != nil && cs.span.IsRecording() {
		slog.Debug("Adding event to collector span",
			"event_name", name,
			"span_id", cs.span.SpanContext().SpanID().String(),
			"trace_id", cs.span.SpanContext().TraceID().String(),
		)
		cs.span.AddEvent(name, trace.WithAttributes(attrs...))
	} else {
		slog.Debug("Skipping event add - span is nil or not recording", "event_name", name)
	}
}

// RecordDuration records the duration of an operation
func (cs *CollectorSpan) RecordDuration(startTime time.Time, operation string) {
	if cs.span != nil && cs.span.IsRecording() {
		duration := time.Since(startTime)
		cs.span.SetAttributes(
			attribute.String("operation", operation),
			attribute.Int64("duration_ms", duration.Milliseconds()),
		)
	}
}
