test:
	go test $(shell go list ./... | grep -vFe'/vendor/')

init:
	dep init

vendor:
	dep ensure

gen:
	go generate $(shell go list ./... | grep -vFe'/vendor/')

.PHONY: test init vendor gen
