package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProjectType represents the detected project type
type ProjectType string

const (
	ProjectGo         ProjectType = "go"
	ProjectPython     ProjectType = "python"
	ProjectRust       ProjectType = "rust"
	ProjectNode       ProjectType = "node"
	ProjectTypeScript ProjectType = "typescript"
	ProjectUnknown    ProjectType = "unknown"
)

// DetectProjectType attempts to detect the project type based on files
func DetectProjectType(path string) ProjectType {
	// Check for Go project
	if fileExists(filepath.Join(path, "go.mod")) {
		return ProjectGo
	}

	// Check for Rust project
	if fileExists(filepath.Join(path, "Cargo.toml")) {
		return ProjectRust
	}

	// Check for Python project
	if fileExists(filepath.Join(path, "pyproject.toml")) ||
		fileExists(filepath.Join(path, "setup.py")) ||
		fileExists(filepath.Join(path, "requirements.txt")) ||
		fileExists(filepath.Join(path, "Pipfile")) ||
		fileExists(filepath.Join(path, "poetry.lock")) {
		return ProjectPython
	}

	// Check for Node/TypeScript project
	if fileExists(filepath.Join(path, "package.json")) {
		// Check if TypeScript
		if fileExists(filepath.Join(path, "tsconfig.json")) {
			return ProjectTypeScript
		}
		return ProjectNode
	}

	return ProjectUnknown
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetTemplateForType returns the appropriate config template for a project type
func GetTemplateForType(projectType ProjectType) string {
	switch projectType {
	case ProjectGo:
		return goTemplate
	case ProjectPython:
		return pythonTemplate
	case ProjectRust:
		return rustTemplate
	case ProjectNode, ProjectTypeScript:
		return nodeTemplate
	default:
		return defaultTemplate
	}
}

// Template configurations
const goTemplate = `# GoWatch Configuration for Go Project
watch:
  - path: "./"
    recursive: true
    ignore:
      - "**/vendor/**"
      - "**/*.exe"
      - "**/bin/**"
      - ".git/**"

on_change:
  commands:
    - cmd: ["go", "fmt", "./..."]
      timeout: "30s"
    - cmd: ["go", "test", "-v", "./..."]
      timeout: "120s"
    - cmd: ["go", "build", "./..."]
      timeout: "90s"

debounce: "500ms"
max_concurrency: 1
`

const pythonTemplate = `# GoWatch Configuration for Python Project
watch:
  - path: "./"
    recursive: true
    ignore:
      - "**/__pycache__/**"
      - "**/venv/**"
      - "**/env/**"
      - "**/.pytest_cache/**"
      - ".git/**"

on_change:
  commands:
    - cmd: ["python", "-m", "pytest", "-v"]
      timeout: "120s"

debounce: "500ms"
max_concurrency: 1
`

const rustTemplate = `# GoWatch Configuration for Rust Project
watch:
  - path: "./src"
    recursive: true
  - path: "Cargo.toml"

ignore:
  - "**/target/**"
  - ".git/**"

on_change:
  commands:
    - cmd: ["cargo", "check"]
      timeout: "60s"
    - cmd: ["cargo", "test"]
      timeout: "180s"
    - cmd: ["cargo", "build"]
      timeout: "120s"

debounce: "750ms"
max_concurrency: 1
`

const nodeTemplate = `# GoWatch Configuration for Node.js/TypeScript Project
watch:
  - path: "./src"
    recursive: true

ignore:
  - "**/node_modules/**"
  - "**/dist/**"
  - "**/build/**"
  - ".git/**"

on_change:
  commands:
    - cmd: ["npm", "run", "lint"]
      timeout: "60s"
    - cmd: ["npm", "test"]
      timeout: "120s"
    - cmd: ["npm", "run", "build"]
      timeout: "180s"

debounce: "500ms"
max_concurrency: 1
`

const defaultTemplate = `# GoWatch Configuration
watch:
  - path: "./"
    recursive: true
    ignore:
      - "vendor/**"
      - ".git/**"
      - "**/*.tmp"

on_change:
  commands:
    - cmd: ["echo", "File changed: {path}"]

debounce: "250ms"
max_concurrency: 2
`

// WriteTemplateForProject writes a config template based on detected project type
func WriteTemplateForProject(path string) error {
	projectType := DetectProjectType(path)
	template := GetTemplateForType(projectType)

	configPath := filepath.Join(path, "gowatch.yaml")

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file already exists: %s", configPath)
	}

	if err := os.WriteFile(configPath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetProjectTypeName returns a human-readable name for the project type
func GetProjectTypeName(pt ProjectType) string {
	switch pt {
	case ProjectGo:
		return "Go"
	case ProjectPython:
		return "Python"
	case ProjectRust:
		return "Rust"
	case ProjectNode:
		return "Node.js"
	case ProjectTypeScript:
		return "TypeScript"
	default:
		return "Unknown"
	}
}
