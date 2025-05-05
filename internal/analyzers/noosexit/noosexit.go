// Пакет noosexit реализует анализатор, запрещающий прямой вызов os.Exit
// в функции main пакета main.
package noosexit

import (
	"go/ast"
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
		ast.Inspect(file, func(n ast.Node) bool {
			// Ищем вызов os.Exit в функции main
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "Exit" {
						if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "os" {
							pass.Reportf(call.Pos(), "direct call to os.Exit in main.main is not allowed")
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
