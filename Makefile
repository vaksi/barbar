VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOFILES := $(wildcard *.go)
STIME := $(shell date +%s)

.PHONY: all coverage clean

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go-get

proto-barbar:
	@echo "  >  Start Generate Proto..."
	protoc -I=proto/users --go_out=plugins=grpc:. proto/users/*.proto
	protoc -I=proto/auth --go_out=plugins=grpc:. proto/auth/*.proto
	@echo "  >  Done Generate Proto..."
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"