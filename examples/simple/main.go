package main

import (
	"context"
	"log"
	"time"

	"github.com/d0ugal/promexporter/app"
	"github.com/d0ugal/promexporter/config"
	"github.com/d0ugal/promexporter/metrics"
)

// SimpleCollector is a basic collector that increments a counter
type SimpleCollector struct {
	registry *metrics.Registry
}

// NewSimpleCollector creates a new simple collector
func NewSimpleCollector(registry *metrics.Registry) *SimpleCollector {
	return &SimpleCollector{
		registry: registry,
	}
}

// Start implements the Collector interface
func (c *SimpleCollector) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Println("Collecting metrics...")
			// In a real collector, you would update metrics here
		}
	}
}

// Stop implements the Collector interface
func (c *SimpleCollector) Stop() {
	log.Println("Stopping collector...")
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml", false)
	if err != nil {
		log.Fatal(err)
	}

	// Create metrics registry
	metricsRegistry := metrics.NewRegistry()

	// Create collector
	collector := NewSimpleCollector(metricsRegistry)

	// Create and run application
	application := app.New("simple-exporter").
		WithConfig(cfg).
		WithMetrics(metricsRegistry).
		WithCollector(collector).
		Build()

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
