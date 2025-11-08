package prometheuslabels

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is a go/analysis analyzer that forbids calls to WithLabelValues
// from Prometheus metrics and suggests using With(prometheus.Labels{...}) instead.
var Analyzer = &analysis.Analyzer{
	Name: "prometheuslabels",
	Doc:  "forbids calls to WithLabelValues from Prometheus metrics, suggests With(prometheus.Labels{...}) instead",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			// Check if this is a call to WithLabelValues
			if selExpr.Sel.Name != "WithLabelValues" {
				return true
			}

			// Report the issue
			pass.Reportf(
				callExpr.Pos(),
				"forbidden: use of WithLabelValues is not allowed; use With(prometheus.Labels{...}) instead",
			)

			return true
		})
	}

	return nil, nil
}
