# Contributing to GoWatch

Thank you for your interest in contributing to GoWatch! We welcome contributions from everyone.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, but recommended)

### Setting Up Development Environment

1. Fork the repository on GitHub
2. Clone your fork:

   ```bash
   git clone https://github.com/YOUR_USERNAME/gowatch.git
   cd gowatch
   ```

3. Add upstream remote:

   ```bash
   git remote add upstream https://github.com/scorpiocodex/gowatch.git
   ```

4. Install dependencies:

   ```bash
   make init
   # or
   go mod download
   go mod tidy
   ```

5. Build the project:

   ```bash
   make build
   # or
   go build -o gowatch ./cmd/gowatch
   ```

## Development Workflow

### Making Changes

1. Create a new branch for your feature/fix:

   ```bash
   git checkout -b feature/my-new-feature
   ```

2. Make your changes following our coding standards

3. Run tests:

   ```bash
   make test
   # or
   go test ./...
   ```

4. Format your code:

   ```bash
   make fmt
   # or
   go fmt ./...
   ```

5. Run linters:

   ```bash
   make lint
   # or
   go vet ./...
   ```

### Commit Messages

We follow conventional commit format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Example:

```
feat(watcher): add support for symlink following

Add configuration option to follow symbolic links when watching
directories. This is useful for development setups using symlinks.

Closes #123
```

### Testing

- Write unit tests for new features
- Ensure all tests pass before submitting PR
- Aim for high test coverage
- Include integration tests where appropriate

```bash
# Run all tests
make test

# Run with race detector
make test-race

# Generate coverage report
make test-cover
```

### Code Style

- Follow Go best practices and idioms
- Use `gofmt` for formatting
- Keep functions small and focused
- Write clear comments for exported functions
- Use meaningful variable names

Example:

```go
// Good
func (w *Watcher) Start(ctx context.Context) (<-chan Event, error) {
    // Implementation
}

// Bad
func (w *Watcher) start(c context.Context) (<-chan Event, error) {
    // Implementation
}
```

## Pull Request Process

1. Update documentation if needed
2. Add tests for new features
3. Ensure all tests pass
4. Update CHANGELOG.md
5. Push to your fork
6. Submit a pull request to the main repository

### PR Checklist

- [ ] Tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linters pass (`make lint`)
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated
- [ ] Commits follow convention
- [ ] PR description explains the changes

## Reporting Issues

### Bug Reports

When reporting bugs, please include:

- GoWatch version
- Operating system and version
- Go version
- Steps to reproduce
- Expected behavior
- Actual behavior
- Error messages/logs
- Configuration file (if applicable)

### Feature Requests

When requesting features, please include:

- Clear description of the feature
- Use case/motivation
- Expected behavior
- Any implementation ideas

## Code Review Process

1. Maintainers will review your PR
2. Address any feedback or requested changes
3. Once approved, your PR will be merged
4. Your contribution will be included in the next release

## Community

- Be respectful and constructive
- Help others when you can
- Follow our Code of Conduct
- Ask questions if you're unsure

## Development Tips

### Debugging

Enable verbose logging:

```bash
./gowatch run --verbose
```

Use dry-run mode:

```bash
./gowatch run --dry-run
```

### Testing Local Changes

Build and test locally:

```bash
make build
./bin/gowatch run --config examples/gowatch.yaml
```

### Useful Commands

```bash
# Run development script
./scripts/dev.sh help

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o gowatch-linux ./cmd/gowatch

# Run specific test
go test -v ./internal/watcher -run TestDebouncer

# Check test coverage
go test -cover ./...
```

## Documentation

- Keep README.md up to date
- Document new features in code comments
- Update examples when adding new functionality
- Add inline comments for complex logic

## Questions?

If you have questions about contributing:

1. Check existing issues and documentation
2. Open a new issue with your question
3. Reach out to maintainers

Thank you for contributing to GoWatch! 🎉
