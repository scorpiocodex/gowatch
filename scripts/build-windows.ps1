# PowerShell Build Script for GoWatch
# Usage: .\scripts\build-windows.ps1 [command]

param(
    [string]$Command = "help"
)

# Colors
function Write-Info($message) {
    Write-Host "[INFO] $message" -ForegroundColor Cyan
}

function Write-Success($message) {
    Write-Host "[SUCCESS] $message" -ForegroundColor Green
}

function Write-Error-Custom($message) {
    Write-Host "[ERROR] $message" -ForegroundColor Red
}

function Write-Header($message) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Yellow
    Write-Host "  $message" -ForegroundColor Yellow
    Write-Host "========================================" -ForegroundColor Yellow
    Write-Host ""
}

# Check Go installation
function Check-Go {
    if (-not (Get-Command "go" -ErrorAction SilentlyContinue)) {
        Write-Error-Custom "Go is not installed or not in PATH"
        Write-Host "Please install Go from https://golang.org/dl/"
        exit 1
    }
    Write-Info "Go version:"
    go version
    Write-Host ""
}

# Build command
function Build {
    Write-Header "Building GoWatch"
    Check-Go
    
    if (-not (Test-Path "bin")) {
        New-Item -ItemType Directory -Path "bin" | Out-Null
    }
    
    Write-Info "Building gowatch.exe..."
    go build -o bin\gowatch.exe .\cmd\gowatch
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Built: bin\gowatch.exe"
        $size = (Get-Item "bin\gowatch.exe").Length / 1MB
        Write-Info "Binary size: $([math]::Round($size, 2)) MB"
    } else {
        Write-Error-Custom "Build failed"
        exit 1
    }
}

# Test command
function Test-Project {
    Write-Header "Running Tests"
    Check-Go
    
    Write-Info "Running all tests..."
    go test -v ./...
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "All tests passed"
    } else {
        Write-Error-Custom "Tests failed"
        exit 1
    }
}

# Test with coverage
function Test-Coverage {
    Write-Header "Running Tests with Coverage"
    Check-Go
    
    Write-Info "Running tests with coverage..."
    go test -coverprofile=coverage.out ./...
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Tests completed"
        Write-Info "Generating HTML coverage report..."
        go tool cover -html=coverage.out -o coverage.html
        Write-Success "Coverage report: coverage.html"
        
        # Show coverage summary
        Write-Info "Coverage summary:"
        go tool cover -func=coverage.out | Select-String "total"
    } else {
        Write-Error-Custom "Tests failed"
        exit 1
    }
}

# Clean command
function Clean {
    Write-Header "Cleaning Build Artifacts"
    
    if (Test-Path "bin") {
        Remove-Item -Recurse -Force "bin"
        Write-Success "Removed bin directory"
    }
    
    if (Test-Path "coverage.out") {
        Remove-Item "coverage.out"
        Write-Success "Removed coverage.out"
    }
    
    if (Test-Path "coverage.html") {
        Remove-Item "coverage.html"
        Write-Success "Removed coverage.html"
    }
    
    Write-Success "Clean complete"
}

# Install command
function Install {
    Write-Header "Installing GoWatch"
    Check-Go
    
    Write-Info "Installing to GOPATH\bin..."
    go install .\cmd\gowatch
    
    if ($LASTEXITCODE -eq 0) {
        $gopath = go env GOPATH
        Write-Success "Installed to: $gopath\bin\gowatch.exe"
        Write-Info "Make sure $gopath\bin is in your PATH"
    } else {
        Write-Error-Custom "Installation failed"
        exit 1
    }
}

# Run command
function Run {
    Build
    if ($LASTEXITCODE -eq 0) {
        Write-Header "Running GoWatch"
        $remainingArgs = $args[1..($args.Length - 1)]
        & .\bin\gowatch.exe $remainingArgs
    }
}

# Format code
function Format {
    Write-Header "Formatting Code"
    Check-Go
    
    Write-Info "Running go fmt..."
    go fmt ./...
    Write-Success "Code formatted"
}

# Lint code
function Lint {
    Write-Header "Linting Code"
    Check-Go
    
    Write-Info "Running go vet..."
    go vet ./...
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Lint checks passed"
    } else {
        Write-Error-Custom "Lint checks failed"
        exit 1
    }
}

# Init project
function Init {
    Write-Header "Initializing Project"
    Check-Go
    
    Write-Info "Downloading dependencies..."
    go mod download
    
    Write-Info "Tidying go.mod..."
    go mod tidy
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Project initialized"
    } else {
        Write-Error-Custom "Initialization failed"
        exit 1
    }
}

# Help command
function Show-Help {
    Write-Header "GoWatch Build Script"
    
    Write-Host "Usage: .\scripts\build-windows.ps1 [command]"
    Write-Host ""
    Write-Host "Commands:" -ForegroundColor Yellow
    Write-Host "  build      - Build the gowatch.exe binary"
    Write-Host "  test       - Run all tests"
    Write-Host "  coverage   - Run tests with coverage report"
    Write-Host "  clean      - Clean build artifacts"
    Write-Host "  install    - Install to GOPATH\bin"
    Write-Host "  run        - Build and run gowatch"
    Write-Host "  fmt        - Format code"
    Write-Host "  lint       - Run linters"
    Write-Host "  init       - Initialize project dependencies"
    Write-Host "  help       - Show this help message"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  .\scripts\build-windows.ps1 build"
    Write-Host "  .\scripts\build-windows.ps1 test"
    Write-Host "  .\scripts\build-windows.ps1 run --help"
    Write-Host ""
}

# Main script logic
Write-Header "GoWatch - Windows Build Script"

switch ($Command.ToLower()) {
    "build"    { Build }
    "test"     { Test-Project }
    "coverage" { Test-Coverage }
    "clean"    { Clean }
    "install"  { Install }
    "run"      { Run $args }
    "fmt"      { Format }
    "lint"     { Lint }
    "init"     { Init }
    "help"     { Show-Help }
    default    { Show-Help }
}

Write-Host ""