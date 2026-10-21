.PHONY := all build-all clean test
.DEFAULT_GOAL := all
SHELL := /usr/bin/bash

BINARY_DIR  := bin
GO_FILES    := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
MODULE_NAME := $(shell go list -m)

VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null)
LDFLAGS    := -w -s -X 'main.Version=$(VERSION)'

BUILD_TARGETS :=
CLEAN_TARGETS :=
INSTALL_TARGETS :=

all: build-all

include cmd/ipar/Makefile.am

build-all: $(BUILD_TARGETS)

clean: $(CLEAN_TARGETS)
	@rm -rf $(BINARY_DIR)
	@echo "Cleaned all artifacts!"

install: $(INSTALL_TARGETS)
	@echo "Installed all artifacts!"

test:
	@echo "Running global test suite..."
	@go test -v -race ./...

lint:
	@golangci-lint run ./...

fmt:
	@go fmt ./...
