# Project metadata
APP_NAME=sales-admin
PKG=app/services/tooling/$(APP_NAME)
BIN=bin/$(APP_NAME)

# Default target
.PHONY: all
all: build

# Build the sales-admin binary
.PHONY: build
build:
	go build -o $(BIN) $(PKG)/main.go

# Run the migration command
.PHONY: migrate
migrate:
	go run $(PKG)/main.go migrate

# Run the seed command
.PHONY: seed
seed:
	go run $(PKG)/main.go seed

# Run the tool interactively (optional default behavior)
.PHONY: run
run:
	go run $(PKG)/main.go

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/

# Tidy and vendor Go modules
.PHONY: deps
deps:
	go mod tidy
	go mod vendor

# Test everything
.PHONY: test
test:
	go test ./...

# Format all Go files
.PHONY: fmt
fmt:
	go fmt ./...

# Lint (you can add golangci-lint or similar tools here)
.PHONY: lint
lint:
	go vet ./...

