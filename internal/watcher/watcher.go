package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches for file changes
type Watcher struct {
	fsWatcher *fsnotify.Watcher
	filePath  string
	fileName  string
	events    chan struct{}
	errors    chan error
	done      chan struct{}
}

// New creates a new Watcher for the specified file
func New(filePath string) (*Watcher, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Get absolute path and directory
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	dir := filepath.Dir(absPath)
	fileName := filepath.Base(absPath)

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	// Watch the directory instead of the file
	if err := fsWatcher.Add(dir); err != nil {
		fsWatcher.Close()
		return nil, fmt.Errorf("failed to watch directory: %w", err)
	}

	w := &Watcher{
		fsWatcher: fsWatcher,
		filePath:  absPath,
		fileName:  fileName,
		events:    make(chan struct{}),
		errors:    make(chan error),
		done:      make(chan struct{}),
	}

	return w, nil
}

// Close stops the watcher and releases resources
func (w *Watcher) Close() error {
	close(w.done)
	return w.fsWatcher.Close()
}

// Start begins watching for file changes
func (w *Watcher) Start() {
	go w.loop()
}

// Events returns a channel that receives notifications when the file changes
func (w *Watcher) Events() <-chan struct{} {
	return w.events
}

// Errors returns a channel that receives watcher errors
func (w *Watcher) Errors() <-chan error {
	return w.errors
}

func (w *Watcher) loop() {
	// Debounce timer to coalesce rapid events
	var debounceTimer *time.Timer
	const debounceInterval = 100 * time.Millisecond

	for {
		select {
		case <-w.done:
			return
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}
			// Filter events by target file name
			if filepath.Base(event.Name) != w.fileName {
				continue
			}
			// Handle Write and Create events (Create handles atomic saves)
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				// Debounce: reset timer on each event
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(debounceInterval, func() {
					select {
					case w.events <- struct{}{}:
					case <-w.done:
					}
				})
			}
		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			select {
			case w.errors <- err:
			case <-w.done:
			}
		}
	}
}
