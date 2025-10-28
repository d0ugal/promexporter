package tracing

import (
	"context"
	"fmt"
	"log/slog"
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
	if !cfg.IsEnabled() {
		return &Tracer{}, nil
	}

	// Create OTLP HTTP exporter
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(cfg.Endpoint),
		otlptracehttp.WithHeaders(cfg.Headers),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("service.version", "1.0.0"), // Could be made configurable
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Could be made configurable
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	// Set global text map propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create tracer
	tracer := tp.Tracer(cfg.ServiceName)

	slog.Info("Tracing initialized", 
		"service_name", cfg.ServiceName,
		"endpoint", cfg.Endpoint,
	)

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
	if !t.IsEnabled() {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		// Extract trace context from headers
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Start span for the HTTP request
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		ctx, span := t.StartSpanWithAttributes(ctx, spanName,
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.user_agent", c.Request.UserAgent()),
		)

		// Store context in Gin context
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Record response information
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)

		// Record error if status indicates error
		if c.Writer.Status() >= 400 {
			span.RecordError(fmt.Errorf("HTTP %d", c.Writer.Status()))
		}

		span.End()
	}
}

// Shutdown gracefully shuts down the tracer
func (t *Tracer) Shutdown(ctx context.Context) error {
	if !t.IsEnabled() {
		return nil
	}

	tp := otel.GetTracerProvider()
	if sdkTp, ok := tp.(*sdktrace.TracerProvider); ok {
		return sdkTp.Shutdown(ctx)
	}

	return nil
}

// CollectorSpan wraps common collector operations with tracing
type CollectorSpan struct {
	span trace.Span
	ctx  context.Context
}

// NewCollectorSpan creates a new collector span
func (t *Tracer) NewCollectorSpan(ctx context.Context, collectorName, operation string) *CollectorSpan {
	if !t.IsEnabled() {
		return &CollectorSpan{ctx: ctx}
	}

	ctx, span := t.StartSpanWithAttributes(ctx, operation,
		attribute.String("service.name", t.config.ServiceName),
		attribute.String("collector.name", collectorName),
		attribute.String("collector.operation", operation),
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
		cs.span.End()
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
		cs.span.AddEvent(name, trace.WithAttributes(attrs...))
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
