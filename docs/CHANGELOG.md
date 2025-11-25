# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-XX

### Added

- Initial release of GoWatch
- File system watching with fsnotify
- Recursive directory watching with configurable ignore patterns
- Debounced event handling to prevent command spam
- Command execution with timeout support
- Concurrent and sequential command execution modes
- Beautiful terminal UI with colorized output
- YAML configuration file support
- CLI interface with cobra
- Comprehensive test suite
- Cross-platform support (Linux, macOS, Windows)
- Graceful shutdown handling
- Placeholder support in commands ({path}, {event})
- Example configurations and ignore files
- CI/CD pipeline with GitHub Actions
- Development scripts and Makefile

### Features

- `gowatch run` - Start watching files and running commands
- `gowatch init` - Create example configuration files
- `gowatch test-config` - Validate and display configuration
- Support for glob patterns in ignore rules
- Real-time command output streaming
- Exit code reporting
- Duration tracking for command execution
- Dry-run mode for testing configurations
- Verbose logging mode

### Security

- Safe command execution with context cancellation
- Command timeout enforcement
- Input validation and sanitization
- Clear security warnings in documentation

## [Unreleased]

### Planned Features

- Desktop notifications for command completion
- Webhook support for remote notifications
- TUI dashboard with interactive controls
- Plugin system for extensibility
- Remote file watching over SSH
- Performance monitoring and statistics
- Configuration hot-reloading
- Multiple configuration file support
- Command history and replay
