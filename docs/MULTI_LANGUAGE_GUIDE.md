# 🌍 Multi-Language Project Guide

Complete guide for using GoWatch with Python, Rust, Web (Node.js/TypeScript), and Go projects.

## 🎯 Quick Start by Project Type

### Python Projects 🐍

```bash
cd your-python-project
gowatch init  # Auto-detects Python and creates optimized config
gowatch run
```

**What it watches:**
- `.py` files in your project
- `pyproject.toml`, `setup.py`

**What it runs:**
- `pytest` for testing
- Optional: `mypy`, `pylint`, `black`

**What it ignores:**
- `__pycache__/`, `.pyc` files
- `venv/`, `env/`, `.venv/`
- `.pytest_cache/`, `.mypy_cache/`

### Rust Projects 🦀

```bash
cd your-rust-project
gowatch init  # Auto-detects Rust and creates optimized config
gowatch run
```

**What it watches:**
- `src/` directory
- `tests/` directory
- `Cargo.toml`

**What it runs:**
- `cargo check` (fast syntax check)
- `cargo test`
- `cargo build`

**What it ignores:**
- `target/` (build artifacts)
- `Cargo.lock` changes

### Web Projects (Node.js/React/Vue) 🌐

```bash
cd your-web-project
gowatch init  # Auto-detects Node.js and creates optimized config
gowatch run
```

**What it watches:**
- `src/` directory
- `public/` directory
- `package.json`

**What it runs:**
- `npm run lint`
- `npm test`
- `npm run build`

**What it ignores:**
- `node_modules/`
- `dist/`, `build/`, `.next/`
- `package-lock.json` changes

### Go Projects 🔵

```bash
cd your-go-project
gowatch init  # Auto-detects Go and creates optimized config
gowatch run
```

**What it watches:**
- All `.go` files
- `go.mod`

**What it runs:**
- `go fmt ./...`
- `go vet ./...`
- `go test ./...`
- `go build ./...`

**What it ignores:**
- `vendor/`
- `bin/`, `*.exe`
- Coverage files

## 📝 Configuration Examples

### Python Django Project

```yaml
# gowatch.yaml for Django
watch:
  - path: "./myapp"
    recursive: true
  - path: "./manage.py"

ignore:
  - "**/__pycache__/**"
  - "**/migrations/**"  # Don't retrigger on migrations
  - "**/static/**"      # Don't watch static files
  - "**/media/**"

on_change:
  commands:
    - cmd: ["python", "manage.py", "test"]
      timeout: "180s"
    
    # Check for missing migrations
    - cmd: ["python", "manage.py", "makemigrations", "--check", "--dry-run"]
      timeout: "30s"

debounce: "500ms"
```

### Python FastAPI Project

```yaml
# gowatch.yaml for FastAPI
watch:
  - path: "./app"
    recursive: true

ignore:
  - "**/__pycache__/**"
  - "**/.pytest_cache/**"

on_change:
  commands:
    - cmd: ["python", "-m", "pytest", "tests/", "-v"]
      timeout: "120s"
    
    - cmd: ["python", "-m", "mypy", "app/"]
      timeout: "60s"

debounce: "500ms"
```

### Rust CLI Project

```yaml
# gowatch.yaml for Rust CLI
watch:
  - path: "./src"
    recursive: true
  - path: "Cargo.toml"

on_change:
  commands:
    - cmd: ["cargo", "clippy", "--", "-D", "warnings"]
      timeout: "90s"
    
    - cmd: ["cargo", "test"]
      timeout: "180s"
    
    - cmd: ["cargo", "build", "--release"]
      timeout: "300s"

debounce: "1s"  # Rust can be slow
```

### React + TypeScript

```yaml
# gowatch.yaml for React + TypeScript
watch:
  - path: "./src"
    recursive: true

ignore:
  - "**/node_modules/**"
  - "**/build/**"

on_change:
  commands:
    # Type check
    - cmd: ["npm", "run", "type-check"]
      timeout: "60s"
    
    # Lint
    - cmd: ["npm", "run", "lint"]
      timeout: "60s"
    
    # Test
    - cmd: ["npm", "test", "--", "--watchAll=false", "--passWithNoTests"]
      timeout: "120s"
    
    # Build
    - cmd: ["npm", "run", "build"]
      timeout: "180s"

debounce: "500ms"
```

### Next.js Project

```yaml
# gowatch.yaml for Next.js
watch:
  - path: "./app"
    recursive: true
  - path: "./pages"
    recursive: true
  - path: "./components"
    recursive: true

ignore:
  - "**/node_modules/**"
  - "**/.next/**"

on_change:
  commands:
    - cmd: ["npm", "run", "lint"]
      timeout: "60s"
    
    - cmd: ["npm", "run", "build"]
      timeout: "240s"

debounce: "750ms"
```

### Go Microservice

```yaml
# gowatch.yaml for Go microservice
watch:
  - path: "./cmd"
    recursive: true
  - path: "./internal"
    recursive: true
  - path: "./pkg"
    recursive: true

ignore:
  - "**/vendor/**"
  - "**/*.exe"

on_change:
  commands:
    - cmd: ["go", "fmt", "./..."]
      timeout: "30s"
    
    - cmd: ["go", "vet", "./..."]
      timeout: "60s"
    
    - cmd: ["go", "test", "-race", "./..."]
      timeout: "180s"
    
    - cmd: ["go", "build", "-o", "bin/service", "./cmd/service"]
      timeout: "90s"

debounce: "500ms"
```

## 🔧 Advanced Multi-Project Setup

### Monorepo with Multiple Languages

```yaml
# gowatch.yaml for monorepo
watch:
  # Backend (Go)
  - path: "./services/api"
    recursive: true
    ignore:
      - "**/vendor/**"
  
  # Frontend (React)
  - path: "./web/src"
    recursive: true
    ignore:
      - "**/node_modules/**"
  
  # Shared libraries
  - path: "./libs"
    recursive: true

on_change:
  commands:
    # Test backend
    - cmd: ["go", "test", "./services/..."]
      timeout: "120s"
    
    # Test frontend
    - cmd: ["npm", "test", "--prefix", "./web"]
      timeout: "120s"
    
    # Build backend
    - cmd: ["go", "build", "-o", "bin/api", "./services/api"]
      timeout: "90s"
    
    # Build frontend
    - cmd: ["npm", "run", "build", "--prefix", "./web"]
      timeout: "180s"

debounce: "1s"
max_concurrency: 2  # Run backend and frontend in parallel
```

### Using Makefiles

```yaml
# gowatch.yaml using make
watch:
  - path: "./"
    recursive: true
    ignore:
      - "**/target/**"
      - "**/node_modules/**"
      - "**/vendor/**"

on_change:
  commands:
    # Use make for complex build logic
    - cmd: ["make", "test"]
      timeout: "300s"
    
    - cmd: ["make", "build"]
      timeout: "300s"

debounce: "750ms"
```

## 🎨 Project-Specific Tips

### Python Tips

**Poetry Projects:**
```yaml
commands:
  - cmd: ["poetry", "run", "pytest"]
  - cmd: ["poetry", "run", "mypy", "."]
```

**Using tox:**
```yaml
commands:
  - cmd: ["tox", "-e", "py39"]
    timeout: "300s"
```

**Django with hot reload:**
- Don't use GoWatch for dev server (Django has its own)
- Use GoWatch for tests and checks only

### Rust Tips

**Faster feedback with cargo-watch style:**
```yaml
commands:
  - cmd: ["cargo", "check"]  # Fast!
    timeout: "60s"
  # Run tests only if check passes
```

**Release builds:**
```yaml
commands:
  - cmd: ["cargo", "build", "--release"]
    timeout: "600s"  # Can be slow!
```

**With cargo-nextest (faster tests):**
```yaml
commands:
  - cmd: ["cargo", "nextest", "run"]
    timeout: "120s"
```

### Web Tips

**Webpack projects:**
```yaml
# Don't watch build output!
ignore:
  - "**/dist/**"
  - "**/*.bundle.js"
```

**Development servers:**
- Don't use GoWatch to run dev servers
- Dev servers have their own hot reload
- Use GoWatch for tests and linting

**Turbo/Rush monorepos:**
```yaml
commands:
  - cmd: ["turbo", "run", "test"]
  - cmd: ["turbo", "run", "build"]
```

### Go Tips

**Working on GoWatch itself:**
```yaml
watch:
  - path: "./internal"
    recursive: true
  - path: "./cmd"
    recursive: true

commands:
  - cmd: ["go", "test", "./internal/..."]
  - cmd: ["go", "build", "-o", "bin/gowatch.exe", "./cmd/gowatch"]
```

**Large projects - watch specific packages:**
```yaml
watch:
  - path: "./pkg/myfeature"
    recursive: true

commands:
  - cmd: ["go", "test", "./pkg/myfeature/..."]
```

## 🚀 Performance Optimization

### 1. Increase Debounce for Large Projects

```yaml
# Small project
debounce: "250ms"

# Medium project
debounce: "500ms"

# Large project / Slow builds
debounce: "1s"
```

### 2. Ignore More Aggressively

```yaml
ignore:
  # All dependencies
  - "**/node_modules/**"
  - "**/vendor/**"
  - "**/target/**"
  - "**/__pycache__/**"
  
  # All build outputs
  - "**/dist/**"
  - "**/build/**"
  - "**/bin/**"
  
  # All lock files
  - "**/package-lock.json"
  - "**/Cargo.lock"
  - "**/poetry.lock"
```

### 3. Watch Only What You Need

```yaml
# ❌ Bad - watches everything
watch:
  - path: "./"

# ✅ Good - specific directories
watch:
  - path: "./src"
  - path: "./lib"
```

### 4. Use Parallel Execution Wisely

```yaml
# Independent commands can run in parallel
max_concurrency: 3

commands:
  - cmd: ["cargo", "test", "--lib"]
  - cmd: ["cargo", "test", "--bins"]
  - cmd: ["cargo", "test", "--examples"]
```

## 🐛 Troubleshooting

### Issue: Commands running for wrong project

**Problem:**
```
Watching Python project
Running: go test ./...  ❌
```

**Solution:**
```bash
cd your-python-project
rm gowatch.yaml
gowatch init  # Will auto-detect Python
```

### Issue: Too many events

**Problem:**
```
04:52:35 [WATCH] REMOVE file1
04:52:35 [WATCH] REMOVE file2
04:52:35 [WATCH] REMOVE file3
# 50 events in 1 second!
```

**Solution:**
```yaml
# Increase debounce
debounce: "2s"

# Add more ignores
ignore:
  - "**/build/**"
  - "**/dist/**"
```

### Issue: Commands timing out

**Problem:**
```
✗ Failed: cargo test (exit: timeout)
```

**Solution:**
```yaml
commands:
  - cmd: ["cargo", "test"]
    timeout: "300s"  # Increase timeout
```

### Issue: Wrong directory

**Problem:**
```
Path: C:/Users/name/Projects/  ❌
```

**Solution:**
```yaml
# Use relative paths
watch:
  - path: "./"  # Current directory
    recursive: true
```

## 📚 Example Projects

Check the `examples/` directory for complete configurations:

- `gowatch-python.yaml` - Python projects
- `gowatch-rust.yaml` - Rust projects  
- `gowatch-web.yaml` - Web projects
- `gowatch-go.yaml` - Go projects
- `gowatch-multi-language.yaml` - Multi-language projects

## 🎯 Best Practices

1. **Always use relative paths** in config
2. **Run GoWatch from project root**
3. **Test with `--dry-run` first**
4. **Use `gowatch init` for auto-detection**
5. **Increase debounce for large projects**
6. **Ignore build outputs and dependencies**
7. **Match commands to your project type**
8. **Use sequential mode for dependent commands**

## 💡 Quick Commands Reference

```bash
# Initialize with auto-detection
gowatch init

# Test your config
gowatch test-config

# Run with custom debounce
gowatch run --debounce 1s

# Dry run to see what would execute
gowatch run --dry-run

# Verbose output for debugging
gowatch run --verbose

# Use specific config file
gowatch run --config custom-config.yaml
```

---

**Happy watching!** 🕵️‍♂️

For more help, see:
- `README.md` - Main documentation
- `WINDOWS.md` - Windows-specific guide
- `QUICKSTART.md` - Quick start guide