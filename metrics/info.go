package metrics

// MetricInfo contains information about a metric for the UI
type MetricInfo struct {
	Name         string
	Help         string
	Labels       []string
	ExampleValue string
}
