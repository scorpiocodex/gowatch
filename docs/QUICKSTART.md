# GoWatch Quick Start Guide

Get up and running with GoWatch in 5 minutes!

## Installation

### Option 1: Install from Source

```bash
go install github.com/scorpiocodex/gowatch/cmd/gowatch@latest
```

### Option 2: Build from Repository

```bash
git clone https://github.com/scorpiocodex/gowatch.git
cd gowatch
make build
sudo cp bin/gowatch /usr/local/bin/
```

### Option 3: Download Binary

Download the latest release from [GitHub Releases](https://github.com/scorpiocodex/gowatch/releases)

## First Run - Quick Mode

Watch the current directory and run a command:

```bash
gowatch run --path . --cmd "echo File changed!"
```

Now modify any file in the current directory and watch the output!

## Setting Up a Project

### Step 1: Initialize Configuration

```bash
cd your-project
gowatch init
```

This creates:

- `gowatch.yaml` - Main configuration
- `.gowatchignore` - Files to ignore

### Step 2: Edit Configuration

Edit `gowatch.yaml`:

```yaml
watch:
  - path: "./"
    recursive: true
    ignore:
      - "vendor/**"
      - ".git/**"

on_change:
  commands:
    - cmd: ["go", "test", "./..."]
      timeout: "60s"
```

### Step 3: Start Watching

```bash
gowatch run
```

## Common Use Cases

### 1. Go Development - Auto Test

```yaml
watch:
  - path: "./"
    recursive: true
    ignore:
      - "vendor/**"
      - "**/testdata/**"

on_change:
  commands:
    - cmd: ["go", "test", "-v", "./..."]
      timeout: "60s"
```

```bash
gowatch run
```

### 2. Web Development - Auto Reload

```yaml
watch:
  - path: "./src"
    recursive: true

on_change:
  commands:
    - cmd: ["npm", "run", "build"]
      timeout: "30s"
    - cmd: ["echo", "Build complete!"]
```

### 3. Documentation - Auto Preview

```yaml
watch:
  - path: "./docs"
    recursive: true
    ignore:
      - "**/*.html"

on_change:
  commands:
    - cmd: ["mkdocs", "build"]
      timeout: "10s"
```

### 4. Continuous Build

```bash
gowatch run --path ./src --cmd "make build" --debounce 500ms
```

## Command Line Usage

### Basic Commands

```bash
# Start watching with config file
gowatch run

# Watch specific path
gowatch run --path ./src

# Run specific command
gowatch run --path . --cmd "make test"

# Test configuration
gowatch test-config

# Initialize project
gowatch init
```

### Useful Flags

```bash
# Dry run (see what would execute)
gowatch run --dry-run

# Verbose output
gowatch run --verbose

# Sequential execution
gowatch run --sequential

# Custom debounce
gowatch run --debounce 500ms

# Disable colors
gowatch run --no-color
```

## Configuration Tips

### Ignore Patterns

Use glob patterns to ignore files:

```yaml
ignore:
  - "vendor/**"          # All files in vendor
  - ".git/**"            # All git files
  - "**/*.tmp"           # All .tmp files
  - "**/node_modules/**" # Node modules anywhere
  - "build/"             # Build directory
```

### Placeholders

Use placeholders in commands:

```yaml
commands:
  - cmd: ["echo", "Changed: {path}"]
  - cmd: ["process", "{path}", "--event={event}"]
```

Available placeholders:

- `{path}` - Full path of changed file
- `{event}` - Event type (WRITE, CREATE, etc.)

### Timeouts

Set timeouts for long-running commands:

```yaml
commands:
  - cmd: ["go", "test", "./..."]
    timeout: "120s"  # 2 minutes
```

## Troubleshooting

### Commands Not Running

```bash
# Test your configuration
gowatch test-config

# Try dry-run mode
gowatch run --dry-run

# Enable verbose logging
gowatch run --verbose
```

### Too Many Events

```bash
# Increase debounce time
gowatch run --debounce 1s

# Add more ignore patterns
```

Edit `gowatch.yaml`:

```yaml
ignore:
  - "**/*.log"
  - "**/*.tmp"
  - "**/cache/**"
```

### Command Timeouts

Increase timeout in config:

```yaml
commands:
  - cmd: ["your-command"]
    timeout: "300s"  # 5 minutes
```

## Examples

### Example 1: Go TDD Workflow

```yaml
watch:
  - path: "./"
    recursive: true
    ignore:
      - "vendor/**"
      - "*.exe"

on_change:
  commands:
    - cmd: ["go", "test", "-v", "-race", "./..."]
      timeout: "60s"
    - cmd: ["go", "vet", "./..."]
      timeout: "30s"

debounce: "300ms"
```

### Example 2: Multi-Stage Build

```yaml
watch:
  - path: "./src"
    recursive: true

on_change:
  commands:
    - cmd: ["echo", "Starting build pipeline..."]
    - cmd: ["go", "fmt", "./..."]
      timeout: "10s"
    - cmd: ["go", "build", "-o", "bin/app", "./cmd/app"]
      timeout: "60s"
    - cmd: ["echo", "Build complete!"]

debounce: "500ms"
max_concurrency: 1
```

### Example 3: Parallel Testing

```yaml
watch:
  - path: "./internal"
    recursive: true

on_change:
  commands:
    - cmd: ["go", "test", "./internal/watcher"]
      timeout: "30s"
    - cmd: ["go", "test", "./internal/runner"]
      timeout: "30s"
    - cmd: ["go", "test", "./internal/config"]
      timeout: "30s"

max_concurrency: 3  # Run all tests in parallel
```

## Next Steps

- Read the [full documentation](README.md)
- Check out [examples/](examples/)
- Join our community
- Report issues or request features

## Getting Help

```bash
# Show help
gowatch --help

# Show command help
gowatch run --help

# Show version
gowatch --version
```

Happy watching! 🕵️‍♂️
