package server

import (
	"embed"
	"html/template"
)

//go:embed templates/*.html
var templateFS embed.FS

// TemplateData holds the data passed to the HTML template
type TemplateData struct {
	ExporterName string
	Version      string
	Commit       string
	BuildDate    string
	Status       string
	Config       map[string]interface{}
	Metrics      []MetricData
}

// MetricData represents a metric for template rendering
type MetricData struct {
	Name         string
	Help         string
	Labels       []string
	ExampleValue string
}

var mainTemplate = template.Must(template.New("index.html").Funcs(template.FuncMap{
	"safeHTML": func(s string) template.HTML {
		// SECURITY NOTE: This bypasses HTML escaping. Only use with trusted content.
		return template.HTML(s) //nolint:gosec // Template helper for trusted HTML
	},
}).ParseFS(templateFS, "templates/index.html"))
