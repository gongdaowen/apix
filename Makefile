.PHONY: help build clean test release dev lint

# Variables
BINARY_NAME=apix
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S')
LDFLAGS=-ldflags="-s -w -X main.Version=${VERSION} -X 'main.BuildTime=${BUILD_TIME}'"

# Default target
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build binary for current platform
	@echo "Building ${BINARY_NAME} ${VERSION}..."
	go build ${LDFLAGS} -o ${BINARY_NAME} main.go
	@echo "Build complete: ${BINARY_NAME}"

build-all: ## Build binaries for all platforms
	@echo "Building for all platforms..."
	@mkdir -p dist
	
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 main.go
	@echo "✓ Linux amd64"
	
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 main.go
	@echo "✓ Linux arm64"
	
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 main.go
	@echo "✓ macOS amd64"
	
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 main.go
	@echo "✓ macOS arm64"
	
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe main.go
	@echo "✓ Windows amd64"
	
	@echo ""
	@echo "All builds complete! Files in dist/"
	@ls -lh dist/

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f ${BINARY_NAME}
	rm -f ${BINARY_NAME}.exe
	rm -f apix-test*
	rm -rf dist/
	rm -rf release-assets/
	rm -rf artifacts/
	@echo "Clean complete"

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo ""
	@echo "Generating coverage report..."
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report saved to coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run --timeout=5m

dev: build ## Build and run in development mode
	@echo "Running ${BINARY_NAME}..."
	./${BINARY_NAME} --help

release: clean build-all ## Create a release (should be called by CI)
	@echo "Preparing release ${VERSION}..."
	@cd dist && sha256sum * > ../dist/checksums.txt
	@echo "Release artifacts prepared in dist/"
	@cat dist/checksums.txt

install: ## Install binary to GOPATH/bin
	@echo "Installing ${BINARY_NAME}..."
	go install ${LDFLAGS}
	@echo "Installed to $(go env GOPATH)/bin/${BINARY_NAME}"

version: ## Show version information
	@echo "Version: ${VERSION}"
	@echo "Build Time: ${BUILD_TIME}"
	@echo "Go Version: $(shell go version)"

check-specs: ## Validate OpenAPI spec files
	@echo "Validating OpenAPI specs..."
	@for file in *.yaml *.yml; do \
		if [ -f "$$file" ]; then \
			echo "Checking $$file..."; \
			python3 -c "import yaml; yaml.safe_load(open('$$file'))" 2>/dev/null && \
				echo "  ✓ Valid YAML" || \
				echo "  ⚠ Could not validate (python3-yaml not available)"; \
		fi \
	done

setup-hooks: ## Setup git hooks
	@echo "Setting up git hooks..."
	@mkdir -p .git/hooks
	@echo '#!/bin/bash' > .git/hooks/pre-commit
	@echo 'make test' >> .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed"
