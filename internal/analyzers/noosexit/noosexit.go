// Пакет noosexit реализует анализатор, запрещающий прямой вызов os.Exit
// в функции main пакета main.
package noosexit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// NewAnalyzer возвращает анализатор, который находит вызовы os.Exit
// в функции main главного пакета.
func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "noosexit",
		Doc:  "reports direct usage of os.Exit in main.main",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if pass.Pkg.Name() != "main" {
			continue
		}

		filename := pass.Fset.Position(file.Pos()).Filename
		if strings.Contains(filename, "go-build") {
			continue
		}

		if strings.HasSuffix(filename, "_test.go") {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Exit" {
				return true
			}

			ident, ok := sel.X.(*ast.Ident)
			if !ok || ident.Name != "os" {
				return true
			}

			pass.Reportf(call.Pos(), "direct call to os.Exit in main.main is not allowed")
			return true
		})
	}

	return nil, nil
}
