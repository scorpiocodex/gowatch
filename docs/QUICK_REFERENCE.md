# 🚀 GoWatch Quick Reference Card

One-page reference for common use cases.

## 📥 Installation

```bash
# Build from source
go build -o gowatch.exe .\cmd\gowatch
```

## 🎬 Quick Start

```bash
# 1. Navigate to your project
cd my-project

# 2. Initialize (auto-detects project type!)
gowatch init

# 3. Start watching
gowatch run
```

## 🗂️ Project Types

| Project | Auto-Detected By | Commands Run |
|---------|------------------|--------------|
| **Python** | `pyproject.toml`, `setup.py`, `requirements.txt` | `pytest` |
| **Rust** | `Cargo.toml` | `cargo check`, `cargo test`, `cargo build` |
| **Node.js** | `package.json` | `npm run lint`, `npm test`, `npm run build` |
| **Go** | `go.mod` | `go fmt`, `go vet`, `go test`, `go build` |

## ⚙️ Common Commands

```bash
# Basic usage
gowatch run                          # Run with gowatch.yaml
gowatch run --config custom.yaml    # Use custom config
gowatch run --dry-run                # See what would run (no execution)

# Testing
gowatch test-config                  # Validate configuration
gowatch run --verbose                # Debug mode

# Custom options
gowatch run --debounce 1s           # Wait 1 second after changes
gowatch run --sequential            # Run commands one-by-one
gowatch run --no-color              # Disable colors
```

## 📝 Config Template

```yaml
watch:
  - path: "./"                # What to watch
    recursive: true           # Watch subdirectories
    ignore:                   # What to skip
      - "vendor/**"
      - ".git/**"

on_change:
  commands:
    - cmd: ["go", "test", "./..."]  # Command as array
      timeout: "60s"                # Max run time

debounce: "500ms"            # Wait after last change
max_concurrency: 2           # Max parallel commands
```

## 🎯 Common Patterns

### Python Django
```yaml
watch:
  - path: "./myapp"
on_change:
  commands:
    - cmd: ["python", "manage.py", "test"]
```

### React
```yaml
watch:
  - path: "./src"
ignore:
  - "**/node_modules/**"
on_change:
  commands:
    - cmd: ["npm", "run", "lint"]
    - cmd: ["npm", "test", "--", "--watchAll=false"]
```

### Rust
```yaml
watch:
  - path: "./src"
on_change:
  commands:
    - cmd: ["cargo", "check"]     # Fast
    - cmd: ["cargo", "test"]      # Full
```

### Go
```yaml
watch:
  - path: "./"
ignore:
  - "**/vendor/**"
on_change:
  commands:
    - cmd: ["go", "test", "./..."]
    - cmd: ["go", "build", "./..."]
```

## 🚫 Common Ignore Patterns

```yaml
ignore:
  # Dependencies
  - "**/node_modules/**"
  - "**/vendor/**"
  - "**/target/**"
  - "**/__pycache__/**"
  
  # Build outputs
  - "**/dist/**"
  - "**/build/**"
  - "**/bin/**"
  - "**/*.exe"
  
  # VCS
  - ".git/**"
  
  # Lock files
  - "**/package-lock.json"
  - "**/Cargo.lock"
  - "**/poetry.lock"
```

## 🐛 Quick Troubleshooting

| Problem | Solution |
|---------|----------|
| **Wrong commands running** | `rm gowatch.yaml && gowatch init` |
| **Too many events** | Increase `debounce: "1s"` |
| **Commands timing out** | Increase `timeout: "300s"` |
| **Wrong directory watched** | Use `path: "./"` not absolute paths |
| **Can't stop with Ctrl+C** | Press Ctrl+C twice |

## 💡 Pro Tips

1. **Always run from project root:**
   ```bash
   cd my-project
   gowatch run
   ```

2. **Test config first:**
   ```bash
   gowatch test-config
   gowatch run --dry-run
   ```

3. **Use project-specific config:**
   ```bash
   # Python project
   cd python-project && gowatch init
   
   # Rust project  
   cd rust-project && gowatch init
   ```

4. **Increase debounce for large projects:**
   ```yaml
   debounce: "1s"  # or "2s" for very large projects
   ```

5. **Watch only what you need:**
   ```yaml
   watch:
     - path: "./src"      # ✅ Specific
   # NOT
     - path: "./"         # ❌ Too broad
   ```

## 📂 File Locations

```
my-project/
├── gowatch.yaml         ← Your config (create with gowatch init)
├── .gowatchignore      ← Ignore patterns (optional)
└── src/                ← Your code
```

## 🎨 Example Workflows

### TDD Workflow (Go)
```yaml
watch:
  - path: "./"
    ignore: ["**/*_test.go"]  # Don't retrigger on test changes
on_change:
  commands:
    - cmd: ["go", "test", "-v", "./..."]
```

### Build on Save (Rust)
```yaml
on_change:
  commands:
    - cmd: ["cargo", "build", "--release"]
      timeout: "300s"
```

### Continuous Testing (Python)
```yaml
on_change:
  commands:
    - cmd: ["pytest", "-x"]  # Stop on first failure
```

## 🔗 Quick Links

- Full docs: `README.md`
- Windows guide: `WINDOWS.md`
- Multi-language: `MULTI_LANGUAGE_GUIDE.md`
- Examples: `examples/` directory

## ❓ Getting Help

```bash
gowatch --help              # General help
gowatch run --help          # Command-specific help
gowatch --version           # Version info
```

---

**Remember:** 
- Use `gowatch init` for auto-detection
- Test with `--dry-run` first
- Run from project root
- Use relative paths

**Happy coding!** 🕵️‍♂️