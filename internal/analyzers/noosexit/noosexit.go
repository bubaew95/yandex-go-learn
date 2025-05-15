package noosexit

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "noosexit",
		Doc:  "reports any usage of os.Exit anywhere in the codebase",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
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
			if !ok {
				return true
			}

			obj := pass.TypesInfo.Uses[ident]
			if pkg, ok := obj.(*types.PkgName); ok && pkg.Imported().Path() == "os" {
				pass.Reportf(call.Pos(), "usage of os.Exit is not allowed")
			}

			return true
		})
	}
	return nil, nil
}
