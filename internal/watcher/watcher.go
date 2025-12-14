package watcher

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches for file changes
type Watcher struct {
	fsWatcher *fsnotify.Watcher
	filePath  string
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

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	if err := fsWatcher.Add(filePath); err != nil {
		fsWatcher.Close()
		return nil, fmt.Errorf("failed to watch file: %w", err)
	}

	w := &Watcher{
		fsWatcher: fsWatcher,
		filePath:  filePath,
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
