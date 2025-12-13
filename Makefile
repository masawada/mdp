BINARY_NAME := mdp
GO := go
GOFLAGS :=
MAIN_PATH := ./cmd/mdp/main.go

.PHONY: all build test e2e-test lint

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

test:
	@echo "Running tests..."
	$(GO) test ./...
	@echo "Tests completed."

e2e-test: build
	@echo "Running e2e tests..."
	./e2e/run.sh
	@echo "E2E tests completed."

lint:
	@echo "Running linter..."
	golangci-lint run
	@echo "Lint completed."
