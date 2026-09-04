// Package watcher provides file watching functionality.
package watcher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches for file changes.
type Watcher struct {
	fsWatcher *fsnotify.Watcher
	// watched holds the absolute paths of the files to watch.
	watched map[string]bool
	events  chan string
	errors  chan error
	done    chan struct{}
}

// New creates a new Watcher for the specified files.
func New(filePaths ...string) (*Watcher, error) {
	if len(filePaths) == 0 {
		return nil, errors.New("no files to watch")
	}

	watched := make(map[string]bool, len(filePaths))
	dirs := make(map[string]bool)
	for _, filePath := range filePaths {
		// Check if file exists
		if _, err := os.Stat(filePath); err != nil {
			return nil, fmt.Errorf("file not found: %w", err)
		}

		// Get absolute path and directory
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path: %w", err)
		}
		watched[absPath] = true
		dirs[filepath.Dir(absPath)] = true
	}

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	// Watch the directories instead of the files
	for dir := range dirs {
		if err := fsWatcher.Add(dir); err != nil {
			_ = fsWatcher.Close()
			return nil, fmt.Errorf("failed to watch directory: %w", err)
		}
	}

	w := &Watcher{
		fsWatcher: fsWatcher,
		watched:   watched,
		events:    make(chan string),
		errors:    make(chan error),
		done:      make(chan struct{}),
	}

	return w, nil
}

// Close stops the watcher and releases resources.
func (w *Watcher) Close() error {
	close(w.done)
	return w.fsWatcher.Close()
}

// Start begins watching for file changes.
func (w *Watcher) Start() {
	go w.loop()
}

// Events returns a channel that receives the absolute path of a file when it changes.
func (w *Watcher) Events() <-chan string {
	return w.events
}

// Errors returns a channel that receives watcher errors.
func (w *Watcher) Errors() <-chan error {
	return w.errors
}

//nolint:cyclop // complexity is acceptable for event loop with debouncing
func (w *Watcher) loop() {
	// Debounce timers per file to coalesce rapid events
	debounceTimers := make(map[string]*time.Timer)
	const debounceInterval = 100 * time.Millisecond

	for {
		select {
		case <-w.done:
			return
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}
			// Filter events by watched files
			if !w.watched[event.Name] {
				continue
			}
			// Handle Write and Create events (Create handles atomic saves)
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				// Debounce: reset timer on each event
				if timer := debounceTimers[event.Name]; timer != nil {
					timer.Stop()
				}
				debounceTimers[event.Name] = time.AfterFunc(debounceInterval, func() {
					select {
					case w.events <- event.Name:
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
