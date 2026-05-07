package app

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/d0ugal/promexporter/config"
	"github.com/d0ugal/promexporter/metrics"
)

// TestRun_GracefulShutdownReturnsNil verifies that Run() returns nil on a
// graceful shutdown (which produces http.ErrServerClosed from
// http.Server.ListenAndServe). Previously this was treated as a hard error,
// causing every consumer to exit non-zero on a clean SIGTERM.
func TestRun_GracefulShutdownReturnsNil(t *testing.T) {
	cfg := &mockConfig{
		logging: &config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		server: &config.ServerConfig{
			Host: "127.0.0.1",
			Port: 0, // ask the OS for a free port
		},
	}

	metricsRegistry := metrics.NewRegistry("test_exporter_info")

	a := New("test-exporter").
		WithConfig(cfg).
		WithMetrics(metricsRegistry).
		Build()

	runErr := make(chan error, 1)
	go func() {
		runErr <- a.Run()
	}()

	// Give the server a moment to bind. With Port: 0 we don't have a way to
	// observe "ready" from outside; 200ms is comfortably enough for ListenAndServe
	// to reach the accept loop on any reasonable test host.
	time.Sleep(200 * time.Millisecond)

	// Signal the test process. Run() registers a SIGTERM handler internally
	// which intercepts this and triggers graceful shutdown.
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("FindProcess: %v", err)
	}

	if err := p.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("Signal SIGTERM: %v", err)
	}

	select {
	case err := <-runErr:
		if err != nil {
			t.Fatalf("expected nil error on graceful shutdown, got: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for Run() to return after SIGTERM")
	}
}
