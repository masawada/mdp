
BINARY_NAME := mdp
GO := go
GOFLAGS :=
MAIN_PATH := ./cmd/mdp/main.go

.PHONY: all build 

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

test:
	@echo "Running tests..."
	$(GO) test ./...
	@echo "Tests completed."
