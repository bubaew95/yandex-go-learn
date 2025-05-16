prcheck:
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
prres:
	curl http://localhost:8080/debug/pprof/heap > profiles/result.pprof
imp:
	goimports -local "github.com/bubaew95/yandex-go-learn" -w .

GIT_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
VERSION    := $(shell git describe --tags --always)

LDFLAGS := -X 'main.buildVersion=$(VERSION)' -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit=$(GIT_COMMIT)'

build:
	go build -ldflags "$(LDFLAGS)" -o bin/app ./cmd/shortener/main.go