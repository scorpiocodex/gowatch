package runner

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"gowatch/internal/config"
	"gowatch/internal/logger"

	"golang.org/x/sync/errgroup"
)

type Runner struct {
	cfg        *config.Config
	log        *logger.Logger
	sequential bool
	dryRun     bool
	mu         sync.Mutex
	running    int
}

type RunResult struct {
	Command  []string
	ExitCode int
	Duration time.Duration
	Error    error
}

func New(cfg *config.Config, log *logger.Logger, sequential, dryRun bool) *Runner {
	return &Runner{
		cfg:        cfg,
		log:        log,
		sequential: sequential,
		dryRun:     dryRun,
	}
}

func (r *Runner) Run(ctx context.Context, eventPath, eventType string) []RunResult {
	commands := r.cfg.OnChange.Commands
	if len(commands) == 0 {
		r.log.Warn("No commands configured to run")
		return nil
	}

	r.log.Separator()
	r.log.Runner("File change detected")
	r.log.Info("  Path:  %s", eventPath)
	r.log.Info("  Event: %s", eventType)
	r.log.Separator()

	results := make([]RunResult, 0, len(commands))

	if r.sequential {
		for i, cmd := range commands {
			r.log.Info("Command %d/%d", i+1, len(commands))
			result := r.executeCommand(ctx, cmd, eventPath, eventType)
			results = append(results, result)
			if result.Error != nil && result.ExitCode != 0 {
				r.log.Error("Command failed, stopping execution chain")
				break
			}
		}
	} else {
		results = r.executeParallel(ctx, commands, eventPath, eventType)
	}

	// Summary
	r.log.Separator()
	successCount := 0
	for _, result := range results {
		if result.ExitCode == 0 {
			successCount++
		}
	}

	if successCount == len(results) {
		r.log.Success("All commands completed successfully (%d/%d)", successCount, len(results))
	} else {
		r.log.Error("Some commands failed (%d/%d succeeded)", successCount, len(results))
	}
	r.log.Separator()

	return results
}

func (r *Runner) executeParallel(ctx context.Context, commands []config.Command, eventPath, eventType string) []RunResult {
	results := make([]RunResult, len(commands))
	g, gctx := errgroup.WithContext(ctx)

	// Limit concurrency
	sem := make(chan struct{}, r.cfg.MaxConcurrency)

	for i, cmd := range commands {
		i, cmd := i, cmd
		g.Go(func() error {
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-gctx.Done():
				return gctx.Err()
			}

			r.log.Info("Command %d/%d (parallel)", i+1, len(commands))
			results[i] = r.executeCommand(gctx, cmd, eventPath, eventType)
			return nil
		})
	}

	g.Wait()
	return results
}

func (r *Runner) executeCommand(ctx context.Context, cmd config.Command, eventPath, eventType string) RunResult {
	cmdWithPlaceholders := r.replacePlaceholders(cmd.Cmd, eventPath, eventType)
	cmdString := strings.Join(cmdWithPlaceholders, " ")

	if r.dryRun {
		r.log.Info("[DRY-RUN] Would execute: %s", cmdString)
		return RunResult{
			Command:  cmdWithPlaceholders,
			ExitCode: 0,
		}
	}

	// Parse timeout
	timeout := 60 * time.Second
	if cmd.Timeout != "" {
		if d, err := time.ParseDuration(cmd.Timeout); err == nil {
			timeout = d
		}
	}

	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	r.log.CommandStart(cmdString)

	// Validate command
	if len(cmdWithPlaceholders) == 0 {
		return RunResult{
			Command:  cmdWithPlaceholders,
			ExitCode: -1,
			Duration: time.Since(start),
			Error:    fmt.Errorf("empty command"),
		}
	}

	// Prepare command - handle shell commands on Windows
	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		// On Windows, check if we need cmd.exe
		if needsShell(cmdWithPlaceholders) {
			// Use cmd.exe /C for shell commands
			shellCmd := strings.Join(cmdWithPlaceholders, " ")
			command = exec.CommandContext(cmdCtx, "cmd.exe", "/C", shellCmd)
		} else {
			command = exec.CommandContext(cmdCtx, cmdWithPlaceholders[0], cmdWithPlaceholders[1:]...)
		}
	} else {
		command = exec.CommandContext(cmdCtx, cmdWithPlaceholders[0], cmdWithPlaceholders[1:]...)
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		r.log.Error("Failed to get stdout pipe: %v", err)
		return RunResult{
			Command:  cmdWithPlaceholders,
			ExitCode: -1,
			Duration: time.Since(start),
			Error:    fmt.Errorf("failed to get stdout pipe: %w", err),
		}
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		r.log.Error("Failed to get stderr pipe: %v", err)
		return RunResult{
			Command:  cmdWithPlaceholders,
			ExitCode: -1,
			Duration: time.Since(start),
			Error:    fmt.Errorf("failed to get stderr pipe: %w", err),
		}
	}

	if err := command.Start(); err != nil {
		r.log.Error("Failed to start command: %v", err)
		return RunResult{
			Command:  cmdWithPlaceholders,
			ExitCode: -1,
			Duration: time.Since(start),
			Error:    fmt.Errorf("failed to start command: %w", err),
		}
	}

	// Stream output
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			r.log.CommandOutput(scanner.Text(), false)
		}
	}()

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			r.log.CommandOutput(scanner.Text(), true)
		}
	}()

	wg.Wait()

	err = command.Wait()
	duration := time.Since(start)

	result := RunResult{
		Command:  cmdWithPlaceholders,
		Duration: duration,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		result.Error = err
		r.log.CommandEnd(cmdString, result.ExitCode, duration)
	} else {
		result.ExitCode = 0
		r.log.CommandEnd(cmdString, 0, duration)
	}

	return result
}

// needsShell determines if a command needs shell interpretation on Windows
func needsShell(cmd []string) bool {
	if len(cmd) == 0 {
		return false
	}

	// Commands that need shell
	shellCommands := map[string]bool{
		"sh":         true,
		"bash":       true,
		"powershell": true,
		"pwsh":       true,
	}

	firstCmd := strings.ToLower(cmd[0])

	// If explicitly using a shell
	if shellCommands[firstCmd] {
		return true
	}

	// If command contains shell operators
	fullCmd := strings.Join(cmd, " ")
	shellOperators := []string{"|", ">", "<", ">>", "&&", "||", "&"}
	for _, op := range shellOperators {
		if strings.Contains(fullCmd, op) {
			return true
		}
	}

	return false
}

func (r *Runner) replacePlaceholders(cmd []string, path, event string) []string {
	result := make([]string, len(cmd))
	for i, part := range cmd {
		part = strings.ReplaceAll(part, "{path}", path)
		part = strings.ReplaceAll(part, "{event}", event)
		result[i] = part
	}
	return result
}
