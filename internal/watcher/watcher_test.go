package watcher

import (
	"os"
	"path/filepath"
	"testing"
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
	defer w.Close()

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
