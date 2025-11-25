package watcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"gowatch/internal/config"
	"gowatch/internal/logger"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	cfg       *config.Config
	log       *logger.Logger
	fsWatcher *fsnotify.Watcher
	debouncer *Debouncer
	mu        sync.Mutex
	watched   map[string]bool
}

type Event struct {
	Path      string
	Op        string
	Timestamp time.Time
}

func New(cfg *config.Config, log *logger.Logger) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	debouncer := NewDebouncer(cfg.GetDebounceDuration())

	return &Watcher{
		cfg:       cfg,
		log:       log,
		fsWatcher: fsw,
		debouncer: debouncer,
		watched:   make(map[string]bool),
	}, nil
}

func (w *Watcher) Start(ctx context.Context) (<-chan Event, error) {
	events := make(chan Event, 100)

	// Add watch paths
	for _, wp := range w.cfg.Watch {
		if err := w.addPath(wp); err != nil {
			return nil, err
		}
	}

	// Start event processing
	go w.processEvents(ctx, events)

	w.log.Watch("Started watching %d path(s)", len(w.cfg.Watch))
	return events, nil
}

func (w *Watcher) addPath(wp config.WatchPath) error {
	absPath, err := filepath.Abs(wp.Path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to stat path %s: %w", absPath, err)
	}

	if info.IsDir() {
		if wp.Recursive {
			return w.addRecursive(absPath)
		}
		return w.addSingle(absPath)
	}

	return w.addSingle(absPath)
}

func (w *Watcher) addSingle(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.watched[path] {
		return nil
	}

	if err := w.fsWatcher.Add(path); err != nil {
		return fmt.Errorf("failed to watch %s: %w", path, err)
	}

	w.watched[path] = true
	w.log.Debug("Watching: %s", path)
	return nil
}

func (w *Watcher) addRecursive(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		// Check ignore patterns
		if w.shouldIgnore(path) {
			w.log.Debug("Ignoring: %s", path)
			return filepath.SkipDir
		}

		return w.addSingle(path)
	})
}

func (w *Watcher) shouldIgnore(path string) bool {
	base := filepath.Base(path)

	// Common ignore patterns
	if strings.HasPrefix(base, ".") && base != "." {
		return true
	}

	// Windows-specific ignores
	if runtime.GOOS == "windows" {
		// Ignore Windows temp files
		if strings.HasPrefix(base, "~$") {
			return true
		}
		// Ignore Windows shortcuts
		if strings.HasSuffix(base, ".lnk") {
			return true
		}
		// Ignore system folders
		systemFolders := []string{
			"$RECYCLE.BIN",
			"System Volume Information",
			"$Recycle.Bin",
		}
		for _, folder := range systemFolders {
			if strings.Contains(path, folder) {
				return true
			}
		}
	}

	if w.cfg.ShouldIgnore(path) {
		return true
	}

	// Check .gowatchignore file
	ignoreFile := filepath.Join(filepath.Dir(path), ".gowatchignore")
	if _, err := os.Stat(ignoreFile); err == nil {
		// File exists, could parse it here
		// For simplicity, we rely on config ignore patterns
	}

	return false
}

func (w *Watcher) processEvents(ctx context.Context, output chan<- Event) {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			w.log.Watch("Stopping watcher")
			w.fsWatcher.Close()
			return

		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				w.log.Debug("Event channel closed")
				return
			}

			// Filter out ignored paths
			if w.shouldIgnore(event.Name) {
				w.log.Debug("Ignored: %s", event.Name)
				continue
			}

			// Filter out CHMOD events if not needed
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				w.log.Debug("Skipping CHMOD event: %s", event.Name)
				continue
			}

			w.log.Debug("Raw event: %s %s", event.Op, event.Name)

			// Handle directory creation (add to watch list)
			if event.Op&fsnotify.Create == fsnotify.Create {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					for _, wp := range w.cfg.Watch {
						absWatchPath, _ := filepath.Abs(wp.Path)
						absEventPath, _ := filepath.Abs(event.Name)

						// Normalize paths for comparison (important on Windows)
						absWatchPath = filepath.Clean(absWatchPath)
						absEventPath = filepath.Clean(absEventPath)

						if wp.Recursive && strings.HasPrefix(absEventPath, absWatchPath) {
							if err := w.addSingle(event.Name); err != nil {
								w.log.Error("Failed to watch new directory: %v", err)
							} else {
								w.log.Debug("Added watch for new directory: %s", event.Name)
							}
						}
					}
				}
			}

			// Debounce the event
			w.debouncer.Add(event.Name, func() {
				ev := Event{
					Path:      event.Name,
					Op:        event.Op.String(),
					Timestamp: time.Now(),
				}

				select {
				case output <- ev:
					w.log.Watch("%s → %s", ev.Op, ev.Path)
				case <-ctx.Done():
					return
				}
			})

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				w.log.Debug("Error channel closed")
				return
			}
			w.log.Error("Watcher error: %v", err)
		}
	}
}

func (w *Watcher) Stop() {
	w.fsWatcher.Close()
}

// Debouncer prevents rapid-fire events
type Debouncer struct {
	delay   time.Duration
	mu      sync.Mutex
	timers  map[string]*time.Timer
	pending map[string]func()
}

func NewDebouncer(delay time.Duration) *Debouncer {
	return &Debouncer{
		delay:   delay,
		timers:  make(map[string]*time.Timer),
		pending: make(map[string]func()),
	}
}

func (d *Debouncer) Add(key string, fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Cancel existing timer for this key
	if timer, exists := d.timers[key]; exists {
		timer.Stop()
	}

	// Store the function
	d.pending[key] = fn

	// Create new timer
	d.timers[key] = time.AfterFunc(d.delay, func() {
		d.mu.Lock()
		fn := d.pending[key]
		delete(d.pending, key)
		delete(d.timers, key)
		d.mu.Unlock()

		if fn != nil {
			fn()
		}
	})
}
