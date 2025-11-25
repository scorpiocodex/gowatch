# Building GoWatch

Complete guide to building GoWatch from source.

## Prerequisites

- **Go**: Version 1.21 or later
- **Git**: For cloning the repository
- **Make**: Optional but recommended

### Checking Prerequisites

```bash
# Check Go version
go version

# Check Git
git --version

# Check Make (optional)
make --version
```

## Quick Build

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/yourname/gowatch.git
cd gowatch

# Download dependencies
go mod download

# Build
go build -o gowatch ./cmd/gowatch

# Run
./gowatch --help
```

## Using Make

The project includes a Makefile for convenient building:

```bash
# Initialize dependencies
make init

# Build binary
make build

# Build is created at: bin/gowatch
./bin/gowatch --help
```

### Makefile Targets

```bash
make build       # Build the binary
make install     # Install to $GOPATH/bin
make test        # Run tests
make test-race   # Run tests with race detector
make test-cover  # Run tests with coverage
make clean       # Clean build artifacts
make fmt         # Format code
make lint        # Run linters
make build-all   # Build for all platforms
make help        # Show help
```

## Development Build

For development with debug symbols:

```bash
go build -gcflags="all=-N -l" -o gowatch-debug ./cmd/gowatch
```

## Production Build

Optimized build for production:

```bash
go build -ldflags="-s -w" -o gowatch ./cmd/gowatch
```

Flags explained:

- `-s`: Omit symbol table
- `-w`: Omit DWARF symbol table
- Result: Smaller binary size

## Cross-Platform Builds

### Linux

```bash
# AMD64
GOOS=linux GOARCH=amd64 go build -o gowatch-linux-amd64 ./cmd/gowatch

# ARM64
GOOS=linux GOARCH=arm64 go build -o gowatch-linux-arm64 ./cmd/gowatch

# ARM (Raspberry Pi)
GOOS=linux GOARCH=arm GOARM=7 go build -o gowatch-linux-arm ./cmd/gowatch
```

### macOS

```bash
# Intel
GOOS=darwin GOARCH=amd64 go build -o gowatch-darwin-amd64 ./cmd/gowatch

# Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o gowatch-darwin-arm64 ./cmd/gowatch
```

### Windows

```bash
# AMD64
GOOS=windows GOARCH=amd64 go build -o gowatch-windows-amd64.exe ./cmd/gowatch

# 386
GOOS=windows GOARCH=386 go build -o gowatch-windows-386.exe ./cmd/gowatch
```

### Build All Platforms

```bash
make build-all
```

This creates binaries in `bin/` for:

- Linux (amd64)
- macOS (amd64, arm64)
- Windows (amd64)

## Installing

### Install to GOPATH

```bash
go install ./cmd/gowatch
```

Binary will be installed to `$GOPATH/bin/gowatch`

### Install System-Wide (Unix/Linux)

```bash
sudo cp gowatch /usr/local/bin/
```

### Install System-Wide (Windows)

Copy `gowatch.exe` to a directory in your PATH.

## Building with Custom Module Path

If you've forked the project:

```bash
# Update go.mod
sed -i 's|github.com/yourname/gowatch|github.com/YOURUSER/gowatch|g' go.mod

# Update imports in Go files
find . -name "*.go" -type f -exec sed -i 's|github.com/yourname/gowatch|github.com/YOURUSER/gowatch|g' {} +

# Download dependencies
go mod tidy

# Build
go build -o gowatch ./cmd/gowatch
```

## Development Workflow

### 1. Initial Setup

```bash
git clone https://github.com/yourname/gowatch.git
cd gowatch
make init
```

### 2. Make Changes

Edit files in your preferred editor.

### 3. Test Changes

```bash
make test
```

### 4. Build and Test

```bash
make build
./bin/gowatch run --config examples/gowatch.yaml --dry-run
```

### 5. Run Linters

```bash
make lint
```

## Troubleshooting

### Dependency Issues

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
rm go.sum
go mod download
go mod tidy
```

### Build Errors

```bash
# Verify Go version
go version  # Should be 1.21+

# Check for syntax errors
go vet ./...

# Format code
go fmt ./...
```

### Import Errors

If you see import errors, ensure module path is correct:

```bash
# Check go.mod
cat go.mod

# Verify imports
grep -r "github.com/yourname/gowatch" .
```

## Docker Build

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o gowatch ./cmd/gowatch

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gowatch .
ENTRYPOINT ["./gowatch"]
```

Build and run:

```bash
docker build -t gowatch .
docker run -v $(pwd):/workspace -w /workspace gowatch run
```

## GitHub Actions

The project includes CI/CD in `.github/workflows/ci.yml`:

- Runs tests on push/PR
- Tests on Linux, macOS, Windows
- Tests with Go 1.21 and 1.22
- Builds binaries for all platforms
- Uploads artifacts

## Verifying Build

After building, verify the binary:

```bash
# Check binary exists
ls -lh gowatch

# Check it runs
./gowatch --version

# Run help
./gowatch --help

# Test with dry-run
./gowatch run --path . --cmd "echo test" --dry-run
```

## Binary Size

Typical binary sizes:

- **Standard build**: ~15-20 MB
- **Optimized build** (`-ldflags="-s -w"`): ~10-12 MB
- **Compressed** (UPX): ~4-5 MB

To compress further (requires UPX):

```bash
upx --best gowatch
```

## Environment Variables

Useful environment variables for building:

```bash
# Enable Go modules
export GO111MODULE=on

# Set GOPATH (if needed)
export GOPATH=$HOME/go

# Add to PATH
export PATH=$PATH:$GOPATH/bin

# Set build cache
export GOCACHE=$HOME/.cache/go-build
```

## Performance Tips

### Speed Up Builds

```bash
# Use build cache
go build -o gowatch ./cmd/gowatch

# Parallel compilation (default)
go build -p 8 -o gowatch ./cmd/gowatch
```

### Reduce Binary Size

```bash
# Strip symbols and debug info
go build -ldflags="-s -w" -o gowatch ./cmd/gowatch

# Use UPX compression
upx --best --lzma gowatch
```

## Next Steps

After building:

1. Run tests: `make test`
2. Try examples: `./gowatch run --config examples/gowatch.yaml`
3. Read docs: [README.md](README.md)
4. Start developing: [CONTRIBUTING.md](CONTRIBUTING.md)

## Getting Help

If you encounter build issues:

1. Check this document
2. Search existing [issues](https://github.com/yourname/gowatch/issues)
3. Create a new issue with:
   - Go version
   - Operating system
   - Error messages
   - Steps to reproduce
