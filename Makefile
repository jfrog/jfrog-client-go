# Makefile for jfrog-client-go

.PHONY: $(MAKECMDGOALS)

# Default target
help:
	@echo "Available targets:"
	@echo "  update-all           - Update all JFrog dependencies to latest versions"
	@echo "  update-build-info-go - Update build-info-go to latest main branch"
	@echo "  update-gofrog        - Update gofrog to latest main branch"
	@echo "  update-archiver      - Update archiver to latest version"
	@echo "  clean                - Clean build artifacts"
	@echo "  test                 - Run tests"
	@echo "  build                - Build the project"

# Update all JFrog dependencies
update-all: update-build-info-go update-gofrog update-archiver
	@echo "All JFrog dependencies updated successfully!"
	@GOPROXY=direct go mod tidy

# Update build-info-go to latest main branch (using direct proxy to bypass Artifactory)
update-build-info-go:
	@echo "Updating build-info-go to latest main branch..."
	@GOPROXY=direct go get github.com/jfrog/build-info-go@main
	@echo "build-info-go updated successfully!"

# Update gofrog to latest main branch
update-gofrog:
	@echo "Updating gofrog to latest main branch..."
	@GOPROXY=direct go get github.com/jfrog/gofrog@master
	@echo "gofrog updated successfully!"

# Update archiver to latest version
update-archiver:
	@echo "Updating archiver to latest version..."
	@GOPROXY=direct go get -u github.com/jfrog/archiver/v3
	@echo "archiver updated successfully!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@go clean
	@go clean -cache
	@go clean -modcache

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Build the project
build:
	@echo "Building project..."
	@go build ./...
