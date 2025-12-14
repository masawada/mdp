package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(tmpFile, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create watcher
	w, err := New(tmpFile)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	defer func() { _ = w.Close() }()

	// Verify watcher is created
	if w == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNew_FileNotFound(t *testing.T) {
	_, err := New("/nonexistent/file.md")
	if err == nil {
		t.Fatal("New() should return error for nonexistent file")
	}
}

func TestWatchFileChange(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(tmpFile, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create watcher
	w, err := New(tmpFile)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	defer func() { _ = w.Close() }()

	w.Start()

	// Modify the file
	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = os.WriteFile(tmpFile, []byte("# Updated"), 0600)
	}()

	// Wait for event
	select {
	case <-w.Events():
		// Success
	case err := <-w.Errors():
		t.Fatalf("Errors() returned: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for event")
	}
}

func TestWatchFileChange_AtomicSave(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(tmpFile, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create watcher
	w, err := New(tmpFile)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	defer func() { _ = w.Close() }()

	w.Start()

	// Simulate atomic save (write to temp file, then rename)
	go func() {
		time.Sleep(100 * time.Millisecond)
		tmpNew := filepath.Join(tmpDir, "test.md.tmp")
		_ = os.WriteFile(tmpNew, []byte("# Updated"), 0600)
		_ = os.Rename(tmpNew, tmpFile)
	}()

	// Wait for event
	select {
	case <-w.Events():
		// Success
	case err := <-w.Errors():
		t.Fatalf("Errors() returned: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for event from atomic save")
	}
}

func TestClose(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(tmpFile, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create watcher
	w, err := New(tmpFile)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	// Close should not return error
	if err := w.Close(); err != nil {
		t.Fatalf("Close() returned error: %v", err)
	}
}
