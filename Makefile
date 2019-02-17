.PHONY: default gen test install

default:
	echo use gen, test, vendor or install

gen:
	go generate ./...

test:
	go test ./...

install:
	go install ./...
