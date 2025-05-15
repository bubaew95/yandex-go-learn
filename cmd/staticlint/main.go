package main

import (
	"github.com/bubaew95/yandex-go-learn/internal/analyzers/noosexit"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "ST1000") {
			mychecks = append(mychecks, v.Analyzer)
			break
		}
	}

	// 4. Добавить два сторонних публичных анализатора (по выбору)
	mychecks = append(mychecks,
		inspect.Analyzer,
		ctrlflow.Analyzer,
	)

	// 5. Добавить свой собственный
	mychecks = append(mychecks, noosexit.NewAnalyzer())

	// Запуск
	multichecker.Main(mychecks...)
}
