prcheck:
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
prres:
	curl http://localhost:8080/debug/pprof/heap > profiles/result.pprof
imp:
	goimports -local "github.com/bubaew95/yandex-go-learn" -w .