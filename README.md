# 🕵️‍♂️ GoWatch - File Watcher & Auto-Runner

[![CI](https://github.com/scorpiocodex/gowatch/workflows/CI/badge.svg)](https://github.com/scorpiocodex/gowatch/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A production-ready, cross-platform file watcher that automatically runs configured commands when files change. Perfect for development workflows, automated testing, and continuous deployment.

## ✨ Features

- 🔍 **Smart File Watching** - Recursive directory watching with configurable ignore patterns
- ⚡ **Debounced Events** - Prevents command spam from rapid file changes
- 🎯 **Flexible Commands** - Run any command with placeholder support for paths and events
- 🔄 **Concurrent Execution** - Run commands in parallel or sequentially
- 🎨 **Beautiful Output** - Colorized, prefixed logs with timestamps
- ⚙️ **Config Files** - YAML configuration with support for multiple watch paths
- 🛡️ **Safe Execution** - Timeouts, context cancellation, and graceful shutdown
- 🧪 **Well Tested** - Comprehensive unit and integration tests
- 🌍 **Cross-Platform** - Works on Linux, macOS, and Windows

## 📦 Installation

### From Source

```bash
go install github.com/scorpiocodex/gowatch/cmd/gowatch@latest
```

### Build from Repository

**Linux/macOS:**

```bash
git clone https://github.com/scorpiocodex/gowatch.git
cd gowatch
go build -o gowatch ./cmd/gowatch
```

**Windows (PowerShell):**

```powershell
git clone https://github.com/scorpiocodex/gowatch.git
cd gowatch
.\scripts\build-windows.ps1 build
```

**Windows (Command Prompt):**

```cmd
git clone https://github.com/scorpiocodex/gowatch.git
cd gowatch
scripts\build-windows.bat build
```

### Using Development Scripts

**Linux/macOS:**

```bash
./scripts/dev.sh build
./scripts/dev.sh install
```

**Windows:**

```powershell
.\scripts\build-windows.ps1 build
.\scripts\build-windows.ps1 install
```

### 💻 Platform-Specific Notes

- **Windows**: See [WINDOWS.md](WINDOWS.md) for complete Windows guide
- **macOS**: Works natively, requires Xcode Command Line Tools
- **Linux**: Works on all major distributions

## 🚀 Quick Start

### 1. Initialize Configuration

```bash
gowatch init
```

This creates two files:

- `gowatch.yaml` - Main configuration file
- `.gowatchignore` - Patterns to ignore (like .gitignore)

### 2. Edit Configuration

Edit `gowatch.yaml` to configure your watch paths and commands:

```yaml
watch:
  - path: "./"
    recursive: true
    ignore:
      - "vendor/**"
      - ".git/**"
      - "**/*.tmp"

on_change:
  commands:
    - cmd: ["go", "test", "./..."]
      run: sequential
      timeout: "60s"

debounce: "250ms"
max_concurrency: 2
```

### 3. Run GoWatch

```bash
# Using config file
gowatch run --config gowatch.yaml

# Quick CLI usage (no config needed)
gowatch run --path . --cmd "go test ./..."

# Dry run to see what would execute
gowatch run --config gowatch.yaml --dry-run
```

### 4. See It in Action

```
╔════════════════════════════════════════════════════════════╗
║  🕵️  GoWatch - File Watcher & Auto-Runner              ║
║                            v1.0.0                          ║
╚════════════════════════════════════════════════════════════╝

── Configuration ──
15:04:05 [INFO ] Loading config from: gowatch.yaml
15:04:05 [✓ OK ] Configuration loaded successfully

── Starting Watcher ──
15:04:05 [✓ OK ] Watcher started successfully
15:04:05 [INFO ] Watching for file changes... (Press Ctrl+C to stop)
────────────────────────────────────────────────────────────

15:04:12 [WATCH] WRITE → internal/watcher/watcher.go
────────────────────────────────────────────────────────────
15:04:12 [EXEC ] File change detected
15:04:12 [INFO ] Path:  internal/watcher/watcher.go
15:04:12 [INFO ] Event: WRITE
────────────────────────────────────────────────────────────
▶ Running: go test ./...
  │ ok      github.com/scorpiocodex/gowatch/internal/config   0.123s
  │ ok      github.com/scorpiocodex/gowatch/internal/runner   0.456s
  │ ok      github.com/scorpiocodex/gowatch/internal/watcher  0.789s
✓ Completed: go test ./... (1.37s)
────────────────────────────────────────────────────────────
15:04:13 [✓ OK ] All commands completed successfully (1/1)
────────────────────────────────────────────────────────────
```

## 📖 Usage Examples

### Development Workflow

Watch Go files and run tests on changes:

```bash
gowatch run --path ./src --cmd "go test -v ./..." --debounce 300ms
```

### Build Automation

Automatically rebuild when source changes:

```yaml
watch:
  - path: "./src"
    recursive: true
    ignore:
      - "**/*_test.go"
      - "**/vendor/**"

on_change:
  commands:
    - cmd: ["go", "build", "-o", "bin/myapp", "./cmd/myapp"]
      timeout: "30s"
    - cmd: ["echo", "Build complete!"]
```

### Multi-Command Pipeline

Run linting, testing, and building in sequence:

```yaml
on_change:
  commands:
    - cmd: ["golangci-lint", "run"]
      run: sequential
      timeout: "60s"
    - cmd: ["go", "test", "./..."]
      run: sequential
      timeout: "120s"
    - cmd: ["go", "build", "./..."]
      run: sequential
      timeout: "30s"
```

### Using Placeholders

Commands support placeholders for dynamic values:

```yaml
on_change:
  commands:
    - cmd: ["echo", "File changed: {path}"]
    - cmd: ["echo", "Event type: {event}"]
```

Available placeholders:

- `{path}` - Full path of the changed file
- `{event}` - Event type (WRITE, CREATE, REMOVE, RENAME, CHMOD)

### Platform-Specific Commands

**Windows (cmd.exe):**

```yaml
commands:
  - cmd: ["cmd", "/C", "echo File changed"]
  - cmd: ["cmd", "/C", "build.bat"]
```

**Windows (PowerShell):**

```yaml
commands:
  - cmd: ["powershell", "-Command", "Write-Host 'Building...'"]
  - cmd: ["powershell", "-File", "build.ps1"]
```

**Linux/macOS (bash):**

```yaml
commands:
  - cmd: ["bash", "-c", "echo 'File changed'"]
  - cmd: ["./build.sh"]
```

## 🎛️ Configuration Reference

### Watch Paths

```yaml
watch:
  - path: "./src"          # Path to watch
    recursive: true        # Watch subdirectories
    ignore:               # Glob patterns to ignore
      - "**/*.tmp"
      - "vendor/**"
      - ".git/**"
```

### Commands

```yaml
on_change:
  commands:
    - cmd: ["go", "test"]  # Command as array (safer)
      run: sequential      # 'sequential' or 'parallel'
      timeout: "60s"       # Maximum execution time
```

### Global Settings

```yaml
debounce: "250ms"        # Wait time after last change
max_concurrency: 2       # Max parallel commands
```

## 🎨 CLI Reference

### Commands

```bash
gowatch run          # Start watching and running commands
gowatch init         # Create example configuration files
gowatch test-config  # Validate and display configuration
gowatch help         # Show help information
```

### Flags (run command)

```bash
--config, -c         Config file path (default: gowatch.yaml)
--path, -p           Path to watch
--cmd                Command to run on change
--debounce, -d       Debounce duration (default: 250ms)
--timeout            Command timeout (default: 60s)
--sequential         Run commands sequentially
--max-concurrency    Maximum concurrent commands (default: 2)
--dry-run            Show what would run without executing
--verbose, -v        Verbose logging
--no-color           Disable colored output
```

## 🎯 Example Output

```
15:04:05 [INFO] 🕵️‍♂️  GoWatch v1.0.0 - File Watcher & Auto-Runner

15:04:05 [INFO] Loading config from: gowatch.yaml
15:04:05 [WATCH] Started watching 1 path(s)
15:04:05 [SUCCESS] Watching for changes... (Press Ctrl+C to stop)

15:04:12 [WATCH] Event: WRITE internal/watcher/watcher.go
15:04:12 [RUNNER] Triggered by: internal/watcher/watcher.go (WRITE)
15:04:12 [RUNNER] Executing: go test ./...
15:04:13   │ ok      github.com/scorpiocodex/gowatch/internal/config   0.123s
15:04:13   │ ok      github.com/scorpiocodex/gowatch/internal/runner   0.456s
15:04:13   │ ok      github.com/scorpiocodex/gowatch/internal/watcher  0.789s
15:04:13 [SUCCESS] Command completed successfully (1.37s)
```

## 🔒 Security Best Practices

### Command Execution Safety

⚠️ **Important Security Considerations:**

1. **Avoid Shell Interpretation**: Use command arrays instead of shell strings:

   ```yaml
   # ✅ Safe - direct execution
   cmd: ["go", "test", "./..."]
   
   # ⚠️ Less safe - shell interpretation
   cmd: ["sh", "-c", "go test ./..."]
   ```

2. **Review Commands**: Always review commands in config files before running, especially from untrusted sources.

3. **Use Dry Run**: Test configurations with `--dry-run` first:

   ```bash
   gowatch run --config gowatch.yaml --dry-run
   ```

4. **Limit Privileges**: Don't run GoWatch with elevated privileges unless absolutely necessary.

5. **Validate Input**: Be cautious with user-provided paths and commands.

### Example Safe Configuration

```yaml
# Safe command execution
on_change:
  commands:
    # Direct execution - no shell
    - cmd: ["go", "build", "-o", "bin/app", "./cmd/app"]
    
    # Specific commands only
    - cmd: ["make", "test"]
    
    # Avoid wildcards and injections
    # - cmd: ["sh", "-c", "rm -rf *"]  # ❌ DANGEROUS
```

## 🧪 Testing

### Run Tests

```bash
# All tests
go test ./...

# With race detector
go test -race ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Using Development Script

```bash
./scripts/dev.sh test        # Run tests
./scripts/dev.sh test-race   # With race detector
./scripts/dev.sh test-cover  # With coverage report
./scripts/dev.sh lint        # Run linters
```

## 🛠️ Development

### Project Structure

```
gowatch/
├── cmd/gowatch/           # Main application entry point
│   └── main.go
├── internal/
│   ├── config/           # Configuration loading and validation
│   ├── logger/           # Structured logging
│   ├── runner/           # Command execution
│   └── watcher/          # File system watching
├── examples/             # Example configurations
├── scripts/              # Development scripts
└── .github/workflows/    # CI/CD configuration
```

### Building

```bash
# Development build
go build -o gowatch ./cmd/gowatch

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o gowatch-linux ./cmd/gowatch
GOOS=darwin GOARCH=amd64 go build -o gowatch-darwin ./cmd/gowatch
GOOS=windows GOARCH=amd64 go build -o gowatch.exe ./cmd/gowatch
```

### Development Workflow

```bash
# Install dependencies
go mod download

# Format code
go fmt ./...
gofmt -w .

# Run linters
go vet ./...

# Run tests
go test -v ./...

# Build and run
./scripts/dev.sh build
./scripts/dev.sh run
```

## 🚧 Extending GoWatch

### Potential Extensions

1. **Desktop Notifications**: Add OS-native notifications when commands complete

   ```go
   // Add to internal/notifier/notifier.go
   func SendNotification(title, message string) error {
       // Platform-specific notification
   }
   ```

2. **Webhook Support**: Send HTTP callbacks on file changes

   ```yaml
   on_change:
     webhook: "https://api.example.com/notify"
     commands: [...]
   ```

3. **Plugin System**: Load and run custom plugins

   ```yaml
   plugins:
     - path: "./plugins/custom-processor.so"
       config: {...}
   ```

4. **TUI Dashboard**: Interactive terminal UI showing:
   - Active watches
   - Running commands
   - Recent events
   - Command history

5. **Remote Watching**: Watch files over SSH/network

   ```yaml
   watch:
     - path: "ssh://server/path/to/watch"
   ```

## 📝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Contribution Guidelines

- Write tests for new features
- Follow Go best practices and idioms
- Run `go fmt` and `go vet` before committing
- Update documentation for user-facing changes
- Keep commits atomic and well-described

## 🐛 Troubleshooting

### Common Issues

**Issue**: Commands not running

- Check config with `gowatch test-config`
- Verify commands work independently
- Try `--dry-run` to see what would execute

**Issue**: Too many events

- Increase debounce duration: `--debounce 500ms`
- Add more ignore patterns
- Check for loops (command modifying watched files)

**Issue**: Permission errors

- Verify file permissions
- Check if watched paths exist
- Avoid watching system directories

**Issue**: Commands timing out

- Increase timeout: `timeout: "120s"`
- Check if commands hang
- Use `--verbose` for detailed logs

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [fsnotify](https://github.com/fsnotify/fsnotify) - Cross-platform file system notifications
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [viper](https://github.com/spf13/viper) - Configuration management
- [color](https://github.com/fatih/color) - Colorized terminal output

## 📞 Support

- 🐛 [Report Bug](https://github.com/scorpiocodex/gowatch/issues)
- 💡 [Request Feature](https://github.com/scorpiocodex/gowatch/issues)
- 📧 [Email Support](mailto:support@example.com)

---

Made with ❤️ by the GoWatch team
