package watcher

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gowatch/internal/config"
	"gowatch/internal/logger"
)

func TestDebouncer(t *testing.T) {
	d := NewDebouncer(100 * time.Millisecond)

	called := make(chan string, 10)

	// Add multiple events for the same key rapidly
	for i := 0; i < 5; i++ {
		d.Add("test-key", func() {
			called <- "executed"
		})
		time.Sleep(20 * time.Millisecond)
	}

	// Wait for debounce period
	time.Sleep(150 * time.Millisecond)

	// Should only be called once
	select {
	case <-called:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("debounced function was not called")
	}

	// Should not be called again
	select {
	case <-called:
		t.Fatal("debounced function was called multiple times")
	case <-time.After(100 * time.Millisecond):
		// Expected
	}
}

func TestDebouncer_MultipleKeys(t *testing.T) {
	d := NewDebouncer(50 * time.Millisecond)

	called := make(chan string, 10)

	// Add events for different keys
	d.Add("key1", func() {
		called <- "key1"
	})

	d.Add("key2", func() {
		called <- "key2"
	})

	// Wait for debounce period
	time.Sleep(100 * time.Millisecond)

	// Both should be called
	results := make(map[string]bool)
	for i := 0; i < 2; i++ {
		select {
		case key := <-called:
			results[key] = true
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timeout waiting for debounced functions")
		}
	}

	if !results["key1"] || !results["key2"] {
		t.Errorf("not all keys were called: %v", results)
	}
}

func TestWatcher_ShouldIgnore(t *testing.T) {
	cfg := &config.Config{
		Watch: []config.WatchPath{
			{
				Path: "/tmp",
				Ignore: []string{
					"*.tmp",
					"vendor/**",
					".git/**",
				},
			},
		},
		Debounce: "100ms",
	}

	log := logger.New(logger.LevelInfo, false)
	w, err := New(cfg, log)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer w.Stop()

	tests := []struct {
		path   string
		ignore bool
	}{
		{"/tmp/test.go", false},
		{"/tmp/test.tmp", true},
		{"/tmp/.hidden", true},
		{"/tmp/vendor/pkg", true},
		{"/tmp/.git/config", true},
		{"/tmp/normal/file.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := w.shouldIgnore(tt.path)
			if result != tt.ignore {
				t.Errorf("shouldIgnore(%s) = %v, want %v", tt.path, result, tt.ignore)
			}
		})
	}
}

func TestWatcher_Integration(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "gowatch-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		Watch: []config.WatchPath{
			{
				Path:      tmpDir,
				Recursive: true,
			},
		},
		Debounce:       "50ms",
		MaxConcurrency: 1,
	}

	log := logger.New(logger.LevelInfo, false)
	w, err := New(cfg, log)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer w.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	events, err := w.Start(ctx)
	if err != nil {
		t.Fatalf("failed to start watcher: %v", err)
	}

	// Give watcher time to initialize
	time.Sleep(100 * time.Millisecond)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Wait for event
	select {
	case event := <-events:
		if event.Path != testFile {
			t.Errorf("expected event for %s, got %s", testFile, event.Path)
		}
		if event.Op == "" {
			t.Error("expected non-empty operation")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for file event")
	}
}

func TestWatcher_RecursiveWatch(t *testing.T) {
	// Create temporary directory structure
	tmpDir, err := os.MkdirTemp("", "gowatch-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	cfg := &config.Config{
		Watch: []config.WatchPath{
			{
				Path:      tmpDir,
				Recursive: true,
			},
		},
		Debounce:       "50ms",
		MaxConcurrency: 1,
	}

	log := logger.New(logger.LevelInfo, false)
	w, err := New(cfg, log)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer w.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	events, err := w.Start(ctx)
	if err != nil {
		t.Fatalf("failed to start watcher: %v", err)
	}

	// Give watcher time to initialize
	time.Sleep(100 * time.Millisecond)

	// Create a test file in subdirectory
	testFile := filepath.Join(subDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Wait for event
	select {
	case event := <-events:
		if event.Path != testFile {
			t.Errorf("expected event for %s, got %s", testFile, event.Path)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for recursive watch event")
	}
}
