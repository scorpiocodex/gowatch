package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"gowatch/internal/config"
	"gowatch/internal/logger"
	"gowatch/internal/runner"
	"gowatch/internal/watcher"

	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	watchPath  string
	command    string
	debounce   string
	dryRun     bool
	verbose    bool
	sequential bool
	noColor    bool
	timeout    string
	maxConcur  int
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gowatch",
	Short: "🕵️‍♂️ File watcher and auto-runner",
	Long: `GoWatch watches filesystem changes and automatically runs configured commands.
	
Perfect for development workflows, testing, and automation.`,
	Version: "1.0.0",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start watching files and running commands",
	Long: `Start the file watcher and execute commands when changes are detected.

Examples:
  # Watch current directory and run tests
  gowatch run --path . --cmd "go test ./..."

  # Use a config file
  gowatch run --config gowatch.yaml

  # Dry run to see what would execute
  gowatch run --config gowatch.yaml --dry-run`,
	RunE: runWatch,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create example configuration files",
	Long:  "Create example gowatch.yaml and .gowatchignore files in the current directory.",
	RunE:  initConfig,
}

var testConfigCmd = &cobra.Command{
	Use:   "test-config",
	Short: "Validate and display configuration",
	Long:  "Load and validate the configuration file, then display the parsed settings.",
	RunE:  testConfig,
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(testConfigCmd)

	// Run command flags
	runCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file path (default: gowatch.yaml)")
	runCmd.Flags().StringVarP(&watchPath, "path", "p", "", "path to watch")
	runCmd.Flags().StringVar(&command, "cmd", "", "command to run on change")
	runCmd.Flags().StringVarP(&debounce, "debounce", "d", "250ms", "debounce duration")
	runCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would run without executing")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")
	runCmd.Flags().BoolVar(&sequential, "sequential", false, "run commands sequentially")
	runCmd.Flags().BoolVar(&noColor, "no-color", false, "disable colored output")
	runCmd.Flags().StringVar(&timeout, "timeout", "60s", "command timeout")
	runCmd.Flags().IntVar(&maxConcur, "max-concurrency", 2, "maximum concurrent commands")

	// Test config flags
	testConfigCmd.Flags().StringVarP(&cfgFile, "config", "c", "gowatch.yaml", "config file path")
}

func runWatch(cmd *cobra.Command, args []string) error {
	// Setup logger
	logLevel := logger.LevelInfo
	if verbose {
		logLevel = logger.LevelDebug
	}
	log := logger.New(logLevel, !noColor)

	// Display banner
	log.Banner("GoWatch - File Watcher & Auto-Runner", "1.0.0")

	// Load or build config
	var cfg *config.Config
	var err error

	if cfgFile != "" || (watchPath == "" && command == "") {
		// Load from file
		if cfgFile == "" {
			cfgFile = "gowatch.yaml"
		}
		log.Section("Configuration")
		log.Info("Loading config from: %s", cfgFile)
		cfg, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		log.Success("Configuration loaded successfully")
	} else {
		// Build from flags
		if watchPath == "" {
			watchPath = "."
		}
		if command == "" {
			return fmt.Errorf("--cmd is required when not using a config file")
		}

		log.Section("Configuration")
		log.Info("Using CLI arguments")

		// Determine shell command based on OS
		shellCmd := []string{"sh", "-c", command}
		if runtime.GOOS == "windows" {
			shellCmd = []string{"cmd", "/C", command}
		}

		cfg = &config.Config{
			Watch: []config.WatchPath{
				{
					Path:      watchPath,
					Recursive: true,
				},
			},
			OnChange: config.OnChange{
				Commands: []config.Command{
					{
						Cmd:     shellCmd,
						Timeout: timeout,
					},
				},
			},
			Debounce:       debounce,
			MaxConcurrency: maxConcur,
		}

		// Validate CLI-based config
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}
		log.Success("Configuration validated")
	}

	// Display configuration summary
	log.Section("Watch Configuration")
	for i, w := range cfg.Watch {
		recursive := ""
		if w.Recursive {
			recursive = " (recursive)"
		}
		log.Info("Path %d: %s%s", i+1, w.Path, recursive)
		if len(w.Ignore) > 0 {
			log.Debug("  Ignoring: %v", w.Ignore)
		}
	}

	log.Section("Commands")
	for i, c := range cfg.OnChange.Commands {
		log.Info("Command %d: %v", i+1, c.Cmd)
		if c.Timeout != "" {
			log.Debug("  Timeout: %s", c.Timeout)
		}
	}

	log.Section("Settings")
	log.Info("Debounce: %s", cfg.Debounce)
	log.Info("Max Concurrency: %d", cfg.MaxConcurrency)
	log.Info("Sequential Mode: %v", sequential)
	if dryRun {
		log.Warn("DRY RUN MODE - Commands will not be executed")
	}

	// Create watcher
	w, err := watcher.New(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	// Create runner
	r := runner.New(cfg, log, sequential, dryRun)

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	// Windows-compatible signal handling
	if runtime.GOOS == "windows" {
		signal.Notify(sigCh, os.Interrupt)
	} else {
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	}

	go func() {
		sig := <-sigCh
		log.Info("")
		log.Warn("Received signal: %v", sig)
		log.Info("Shutting down gracefully...")
		cancel()
	}()

	// Start watching
	log.Section("Starting Watcher")
	events, err := w.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start watcher: %w", err)
	}

	log.Success("Watcher started successfully")
	log.Info("Watching for file changes... (Press Ctrl+C to stop)")
	log.Separator()

	// Process events
	eventCount := 0
	for {
		select {
		case <-ctx.Done():
			log.Info("")
			log.Section("Shutdown")
			log.Info("Events processed: %d", eventCount)
			log.Success("Shutdown complete")
			return nil

		case event, ok := <-events:
			if !ok {
				log.Info("Event channel closed")
				return nil
			}

			eventCount++

			// Run commands
			results := r.Run(ctx, event.Path, event.Op)

			// Check for failures
			hasFailure := false
			for _, result := range results {
				if result.ExitCode != 0 {
					hasFailure = true
				}
			}

			if hasFailure && !dryRun {
				log.Error("Execution completed with errors")
			}
		}
	}
}

func initConfig(cmd *cobra.Command, args []string) error {
	log := logger.New(logger.LevelInfo, !noColor)

	log.Banner("GoWatch Initialization", "1.0.0")

	// Detect project type
	cwd, _ := os.Getwd()
	projectType := config.DetectProjectType(cwd)

	log.Section("Project Detection")
	if projectType != config.ProjectUnknown {
		log.Success("Detected project type: %s", config.GetProjectTypeName(projectType))
	} else {
		log.Info("Could not detect project type, using default template")
	}

	log.Section("Creating Configuration Files")

	configPath := "gowatch.yaml"
	ignorePath := ".gowatchignore"

	// Check if files already exist
	configExists := false
	ignoreExists := false

	if _, err := os.Stat(configPath); err == nil {
		log.Warn("Config file already exists: %s", configPath)
		configExists = true
	}

	if _, err := os.Stat(ignorePath); err == nil {
		log.Warn("Ignore file already exists: %s", ignorePath)
		ignoreExists = true
	}

	// Create files
	if !configExists {
		// Use project-specific template
		if err := config.WriteTemplateForProject(cwd); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
		log.Success("Created: %s (optimized for %s)", configPath, config.GetProjectTypeName(projectType))
	}

	if !ignoreExists {
		if err := config.WriteExampleIgnore(ignorePath); err != nil {
			return fmt.Errorf("failed to write example ignore: %w", err)
		}
		log.Success("Created: %s", ignorePath)
	}

	if !configExists || !ignoreExists {
		log.Section("Next Steps")
		log.Info("1. Review %s (customized for your project type)", configPath)
		log.Info("2. Customize %s with your ignore patterns", ignorePath)
		log.Info("3. Test your config: gowatch test-config")
		log.Info("4. Start watching: gowatch run")

		if projectType != config.ProjectUnknown {
			log.Info("")
			log.Info("💡 Tip: The config has been optimized for %s projects!", config.GetProjectTypeName(projectType))
		}
	} else {
		log.Info("")
		log.Info("All files already exist. No changes made.")
	}

	return nil
}

func testConfig(cmd *cobra.Command, args []string) error {
	log := logger.New(logger.LevelInfo, !noColor)

	log.Banner("GoWatch Configuration Test", "1.0.0")
	log.Section("Loading Configuration")
	log.Info("Config file: %s", cfgFile)

	cfg, err := config.Load(cfgFile)
	if err != nil {
		log.Error("Failed to load config: %v", err)
		return err
	}

	log.Success("Configuration loaded successfully")

	log.Section("Watch Paths")
	for i, w := range cfg.Watch {
		recursive := ""
		if w.Recursive {
			recursive = " (recursive)"
		}
		log.Info("%d. %s%s", i+1, w.Path, recursive)
		if len(w.Ignore) > 0 {
			for _, pattern := range w.Ignore {
				log.Debug("   Ignore: %s", pattern)
			}
		}
	}

	log.Section("Commands")
	for i, c := range cfg.OnChange.Commands {
		log.Info("%d. %v", i+1, c.Cmd)
		if c.Timeout != "" {
			log.Debug("   Timeout: %s", c.Timeout)
		}
		if c.Run != "" {
			log.Debug("   Mode: %s", c.Run)
		}
	}

	log.Section("Settings")
	log.Info("Debounce: %s", cfg.Debounce)
	log.Info("Max Concurrency: %d", cfg.MaxConcurrency)

	log.Section("Validation")
	log.Success("All configuration checks passed!")

	return nil
}
