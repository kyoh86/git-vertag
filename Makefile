.PHONY: gen lint test install man

VERSION := `git vertag get`
COMMIT  := `git rev-parse HEAD`

.SUFFIXES: .y .go
.y.go:
	goyacc -o $@ $<
	gofmt -w $@
	rm y.output

gen: ./internal/semver/parse.go
	go generate ./...

lint: gen
	golangci-lint run

test: lint
	go test -v --race ./...

install: test
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT)" ./...

man:
	go run main.go --help-man > git-vertag.1
