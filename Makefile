.DEFAULT_GOAL := help
.PHONY: help

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

lint:
	@yamllint . \
	&& go vet ./... \
	&& golangci-lint run

test:
	@go test ./...
