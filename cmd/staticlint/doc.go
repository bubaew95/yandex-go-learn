/*
Package staticlint is a custom static analysis checker combining:

- Go standard analysis passes (e.g. `assign`, `shadow`, `printf`)
- All SA class analyzers from staticcheck.io
- At least one non-SA analyzer from staticcheck.io (e.g. ST1000)
- Two public analyzers (add their source & purpose here)
- Custom analyzer "noosexit": blocks usage of os.Exit in main.main

## Usage

Build the checker:

	go build -o staticlint ./cmd/staticlint

Run it on your code:

	./staticlint ./...
*/
package main
