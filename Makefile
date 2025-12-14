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

# Run tests in Docker to avoid interference from user's local config files
# To run without Docker: go test ./...
test:
	@echo "Running tests in Docker..."
	docker run --rm -v $(CURDIR):/app -v mdp-go-cache:/go -w /app golang:1.25-alpine go test ./...
	@echo "Tests completed."

e2e-test: build
	@echo "Running e2e tests..."
	./e2e/run.sh
	@echo "E2E tests completed."

lint:
	@echo "Running linter..."
	golangci-lint run
	@echo "Lint completed."
