package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Watch          []WatchPath `mapstructure:"watch"`
	OnChange       OnChange    `mapstructure:"on_change"`
	Debounce       string      `mapstructure:"debounce"`
	MaxConcurrency int         `mapstructure:"max_concurrency"`
}

type WatchPath struct {
	Path      string   `mapstructure:"path"`
	Recursive bool     `mapstructure:"recursive"`
	Ignore    []string `mapstructure:"ignore"`
}

type OnChange struct {
	Commands []Command `mapstructure:"commands"`
}

type Command struct {
	Cmd     []string `mapstructure:"cmd"`
	Run     string   `mapstructure:"run"`
	Timeout string   `mapstructure:"timeout"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("gowatch")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.config/gowatch")
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set defaults
	if cfg.Debounce == "" {
		cfg.Debounce = "250ms"
	}
	if cfg.MaxConcurrency == 0 {
		cfg.MaxConcurrency = 2
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if len(c.Watch) == 0 {
		return fmt.Errorf("at least one watch path is required")
	}

	// Validate debounce duration
	if _, err := time.ParseDuration(c.Debounce); err != nil {
		return fmt.Errorf("invalid debounce duration: %w", err)
	}

	// Validate watch paths exist
	for i, w := range c.Watch {
		if w.Path == "" {
			return fmt.Errorf("watch path %d: path is empty", i)
		}
		absPath, err := filepath.Abs(w.Path)
		if err != nil {
			return fmt.Errorf("watch path %d: invalid path %s: %w", i, w.Path, err)
		}
		if _, err := os.Stat(absPath); err != nil {
			return fmt.Errorf("watch path %d does not exist: %s", i, absPath)
		}
	}

	// Validate commands
	if len(c.OnChange.Commands) == 0 {
		return fmt.Errorf("at least one command is required")
	}

	for i, cmd := range c.OnChange.Commands {
		if len(cmd.Cmd) == 0 {
			return fmt.Errorf("command %d: cmd is empty", i)
		}
		if cmd.Timeout != "" {
			if _, err := time.ParseDuration(cmd.Timeout); err != nil {
				return fmt.Errorf("command %d: invalid timeout: %w", i, err)
			}
		}
	}

	// Validate max concurrency
	if c.MaxConcurrency < 1 {
		return fmt.Errorf("max_concurrency must be at least 1")
	}

	return nil
}

func (c *Config) GetDebounceDuration() time.Duration {
	d, _ := time.ParseDuration(c.Debounce)
	return d
}

func WriteExample(path string) error {
	exampleConfig := `# GoWatch Configuration Example
# Watch paths and patterns
watch:
  - path: "./"
    recursive: true
    ignore:
      - "vendor/**"
      - ".git/**"
      - "**/*.tmp"
      - "**/*.log"
      - "**/node_modules/**"

# Commands to run when files change
on_change:
  commands:
    - cmd: ["go", "test", "./..."]
      run: sequential
      timeout: "60s"
    
    - cmd: ["go", "build", "./..."]
      run: sequential
      timeout: "30s"

# Debounce interval (prevent rapid retriggering)
debounce: "250ms"

# Maximum concurrent command executions
max_concurrency: 2
`

	if err := os.WriteFile(path, []byte(exampleConfig), 0644); err != nil {
		return fmt.Errorf("failed to write example config: %w", err)
	}

	return nil
}

func WriteExampleIgnore(path string) error {
	exampleIgnore := `# GoWatch Ignore File
# Similar to .gitignore - patterns to exclude from watching

# Dependencies
vendor/
node_modules/

# Version control
.git/
.svn/

# Temp files
*.tmp
*.log
*.swp
*~

# Build outputs
*.exe
*.dll
*.so
*.dylib
dist/
build/

# IDE
.idea/
.vscode/
*.iml

# OS
.DS_Store
Thumbs.db
Desktop.ini
`

	// Platform-specific additions
	if runtime.GOOS == "windows" {
		exampleIgnore += `
# Windows-specific
~$*
*.lnk
$RECYCLE.BIN/
System Volume Information/
`
	}

	if err := os.WriteFile(path, []byte(exampleIgnore), 0644); err != nil {
		return fmt.Errorf("failed to write example ignore file: %w", err)
	}

	return nil
}

func (c *Config) ShouldIgnore(path string) bool {
	// Normalize path separators for cross-platform compatibility
	path = filepath.ToSlash(path)
	base := filepath.Base(path)

	for _, w := range c.Watch {
		for _, pattern := range w.Ignore {
			// Normalize pattern as well
			pattern = filepath.ToSlash(pattern)

			// Check base name match
			matched, err := filepath.Match(pattern, base)
			if err == nil && matched {
				return true
			}

			// Check full path match
			matched, err = filepath.Match(pattern, path)
			if err == nil && matched {
				return true
			}

			// Check glob pattern with ** (recursive)
			if strings.Contains(pattern, "**") {
				// Convert ** glob to simple contains check
				parts := strings.Split(pattern, "**")
				if len(parts) == 2 {
					prefix := strings.TrimSuffix(parts[0], "/")
					suffix := strings.TrimPrefix(parts[1], "/")

					if (prefix == "" || strings.HasPrefix(path, prefix)) &&
						(suffix == "" || strings.HasSuffix(path, suffix) || strings.Contains(path, suffix)) {
						return true
					}
				}
			}

			// Additional check for directory patterns
			// e.g., "vendor/" should match "vendor" directory
			if strings.HasSuffix(pattern, "/") {
				dirPattern := strings.TrimSuffix(pattern, "/")
				if strings.HasPrefix(path, dirPattern+"/") || path == dirPattern {
					return true
				}
			}
		}
	}
	return false
}
