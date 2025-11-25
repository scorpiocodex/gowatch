package runner

import (
	"context"
	"testing"
	"time"

	"gowatch/internal/config"
	"gowatch/internal/logger"
)

func TestRunner_ReplacePlaceholders(t *testing.T) {
	cfg := &config.Config{}
	log := logger.New(logger.LevelInfo, false)
	r := New(cfg, log, false, false)

	tests := []struct {
		name     string
		cmd      []string
		path     string
		event    string
		expected []string
	}{
		{
			name:     "no placeholders",
			cmd:      []string{"echo", "hello"},
			path:     "/tmp/test.go",
			event:    "WRITE",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "path placeholder",
			cmd:      []string{"echo", "{path}"},
			path:     "/tmp/test.go",
			event:    "WRITE",
			expected: []string{"echo", "/tmp/test.go"},
		},
		{
			name:     "event placeholder",
			cmd:      []string{"echo", "Event: {event}"},
			path:     "/tmp/test.go",
			event:    "WRITE",
			expected: []string{"echo", "Event: WRITE"},
		},
		{
			name:     "multiple placeholders",
			cmd:      []string{"process", "{path}", "--event={event}"},
			path:     "/tmp/test.go",
			event:    "CREATE",
			expected: []string{"process", "/tmp/test.go", "--event=CREATE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.replacePlaceholders(tt.cmd, tt.path, tt.event)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d parts, got %d", len(tt.expected), len(result))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("part %d: expected %q, got %q", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestRunner_DryRun(t *testing.T) {
	cfg := &config.Config{
		OnChange: config.OnChange{
			Commands: []config.Command{
				{Cmd: []string{"echo", "test"}},
			},
		},
		MaxConcurrency: 1,
	}
	log := logger.New(logger.LevelInfo, false)
	r := New(cfg, log, false, true)

	ctx := context.Background()
	results := r.Run(ctx, "/tmp/test.go", "WRITE")

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].ExitCode != 0 {
		t.Errorf("expected exit code 0 in dry run, got %d", results[0].ExitCode)
	}

	if results[0].Error != nil {
		t.Errorf("expected no error in dry run, got %v", results[0].Error)
	}
}

func TestRunner_ExecuteCommand(t *testing.T) {
	cfg := &config.Config{
		MaxConcurrency: 1,
	}
	log := logger.New(logger.LevelInfo, false)
	r := New(cfg, log, false, false)

	ctx := context.Background()

	// Test successful command
	cmd := config.Command{
		Cmd:     []string{"echo", "hello"},
		Timeout: "5s",
	}

	result := r.executeCommand(ctx, cmd, "/tmp/test.go", "WRITE")

	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}

	if result.Error != nil {
		t.Errorf("expected no error, got %v", result.Error)
	}

	if result.Duration == 0 {
		t.Error("expected non-zero duration")
	}
}

func TestRunner_ExecuteCommand_Timeout(t *testing.T) {
	cfg := &config.Config{
		MaxConcurrency: 1,
	}
	log := logger.New(logger.LevelInfo, false)
	r := New(cfg, log, false, false)

	ctx := context.Background()

	// Test command that times out
	cmd := config.Command{
		Cmd:     []string{"sleep", "10"},
		Timeout: "100ms",
	}

	result := r.executeCommand(ctx, cmd, "/tmp/test.go", "WRITE")

	if result.ExitCode == 0 {
		t.Error("expected non-zero exit code for timeout")
	}

	if result.Error == nil {
		t.Error("expected error for timeout")
	}
}

func TestRunner_Sequential(t *testing.T) {
	cfg := &config.Config{
		OnChange: config.OnChange{
			Commands: []config.Command{
				{Cmd: []string{"echo", "first"}},
				{Cmd: []string{"echo", "second"}},
			},
		},
		MaxConcurrency: 1,
	}
	log := logger.New(logger.LevelInfo, false)
	r := New(cfg, log, true, false)

	ctx := context.Background()
	results := r.Run(ctx, "/tmp/test.go", "WRITE")

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	for i, result := range results {
		if result.ExitCode != 0 {
			t.Errorf("command %d: expected exit code 0, got %d", i, result.ExitCode)
		}
	}
}

func TestRunner_Parallel(t *testing.T) {
	cfg := &config.Config{
		OnChange: config.OnChange{
			Commands: []config.Command{
				{Cmd: []string{"echo", "first"}},
				{Cmd: []string{"echo", "second"}},
				{Cmd: []string{"echo", "third"}},
			},
		},
		MaxConcurrency: 2,
	}
	log := logger.New(logger.LevelInfo, false)
	r := New(cfg, log, false, false)

	ctx := context.Background()
	start := time.Now()
	results := r.Run(ctx, "/tmp/test.go", "WRITE")
	duration := time.Since(start)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	for i, result := range results {
		if result.ExitCode != 0 {
			t.Errorf("command %d: expected exit code 0, got %d", i, result.ExitCode)
		}
	}

	// Parallel execution should be faster than sequential
	// This is a rough check - in reality, echo commands are very fast
	if duration > 5*time.Second {
		t.Errorf("parallel execution took too long: %v", duration)
	}
}
