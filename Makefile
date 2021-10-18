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
	protoc --go_out=proto/users/ --go_opt=paths=source_relative --go-grpc_out=proto/users/ --go-grpc_opt=paths=source_relative --proto_path=proto/users/ proto/users/users.general.proto
	protoc --go_out=proto/users/ --go_opt=paths=source_relative --go-grpc_out=proto/users/ --go-grpc_opt=paths=source_relative --proto_path=proto/users/ proto/users/users.service.proto
	protoc --go_out=proto/auth/ --go_opt=paths=source_relative --go-grpc_out=proto/auth/ --go-grpc_opt=paths=source_relative --proto_path=proto/auth/ proto/auth/auth.general.proto
	protoc --go_out=proto/auth/ --go_opt=paths=source_relative --go-grpc_out=proto/auth/ --go-grpc_opt=paths=source_relative --proto_path=proto/auth/ proto/auth/auth.service.proto
	@echo "  >  Done Generate Proto..."
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"