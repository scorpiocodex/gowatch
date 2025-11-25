#!/bin/bash
# Development helper script for GoWatch

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
info() {
    echo -e "${BLUE}ℹ ${NC}$1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Commands
cmd_help() {
    cat <<EOF
GoWatch Development Script

Usage: ./scripts/dev.sh [command]

Commands:
  build       Build the gowatch binary
  test        Run all tests
  test-race   Run tests with race detector
  test-cover  Run tests with coverage report
  lint        Run linters (go vet, gofmt)
  fmt         Format all Go files
  clean       Clean build artifacts
  install     Build and install to \$GOPATH/bin
  run         Build and run gowatch
  demo        Run a demo with example config
  help        Show this help message

Examples:
  ./scripts/dev.sh build
  ./scripts/dev.sh test
  ./scripts/dev.sh run
EOF
}

cmd_build() {
    info "Building gowatch..."
    go build -o bin/gowatch ./cmd/gowatch
    success "Built: bin/gowatch"
}

cmd_test() {
    info "Running tests..."
    go test -v ./...
    success "Tests passed"
}

cmd_test_race() {
    info "Running tests with race detector..."
    go test -race -v ./...
    success "Tests passed (race detector enabled)"
}

cmd_test_cover() {
    info "Running tests with coverage..."
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    success "Coverage report generated: coverage.html"
    
    # Show coverage summary
    go tool cover -func=coverage.out | grep total
}

cmd_lint() {
    info "Running linters..."
    
    info "Running go vet..."
    go vet ./...
    success "go vet passed"
    
    info "Checking formatting..."
    unformatted=$(gofmt -l .)
    if [ -n "$unformatted" ]; then
        error "The following files need formatting:"
        echo "$unformatted"
        exit 1
    fi
    success "All files are properly formatted"
}

cmd_fmt() {
    info "Formatting Go files..."
    gofmt -w .
    success "Files formatted"
}

cmd_clean() {
    info "Cleaning build artifacts..."
    rm -rf bin/
    rm -f coverage.out coverage.html
    success "Cleaned"
}

cmd_install() {
    info "Installing gowatch..."
    go install ./cmd/gowatch
    success "Installed to: $(go env GOPATH)/bin/gowatch"
}

cmd_run() {
    cmd_build
    info "Running gowatch..."
    ./bin/gowatch "$@"
}

cmd_demo() {
    info "Running demo with example config..."
    
    # Create example config if it doesn't exist
    if [ ! -f "examples/gowatch.yaml" ]; then
        error "Example config not found: examples/gowatch.yaml"
        exit 1
    fi
    
    cmd_build
    info "Starting gowatch with example config..."
    info "Press Ctrl+C to stop"
    ./bin/gowatch run --config examples/gowatch.yaml
}

# Main
main() {
    case "${1:-help}" in
        build)
            cmd_build
            ;;
        test)
            cmd_test
            ;;
        test-race)
            cmd_test_race
            ;;
        test-cover)
            cmd_test_cover
            ;;
        lint)
            cmd_lint
            ;;
        fmt)
            cmd_fmt
            ;;
        clean)
            cmd_clean
            ;;
        install)
            cmd_install
            ;;
        run)
            shift
            cmd_run "$@"
            ;;
        demo)
            cmd_demo
            ;;
        help|--help|-h)
            cmd_help
            ;;
        *)
            error "Unknown command: $1"
            echo ""
            cmd_help
            exit 1
            ;;
    esac
}

main "$@"