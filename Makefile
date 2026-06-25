.PHONY: install build test

VERSION ?= v0.1.0-beta.1
LDFLAGS := -X github.com/noodl-labs/ConnorLLM/services/runtime/internal/cli.Version=$(VERSION)

install:
	cd services/runtime && go install -ldflags "$(LDFLAGS)" ./cmd/connor

build:
	cd services/runtime && go build -ldflags "$(LDFLAGS)" -o bin/connor ./cmd/connor

test:
	cd services/runtime && go test ./...
