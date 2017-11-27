.PHONY: default gen test vendor install

default:
	echo use gen, test, vendor or install

gen:
	go generate ./...

test:
	go test ./...

vendor:
	dep ensure

install:
	go install ./...
