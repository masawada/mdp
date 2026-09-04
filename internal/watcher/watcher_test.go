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
	if err := os.WriteFile(tmpFile, []byte("# Test"), 0644); err != nil { //nolint:gosec // G306: test file in temp dir
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

func TestNew_NoFiles(t *testing.T) {
	_, err := New()
	if err == nil {
		t.Fatal("New() should return error when no files are given")
	}
}

func TestNew_OneOfMultipleFilesNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(tmpFile, []byte("# Test"), 0644); err != nil { //nolint:gosec // G306: test file in temp dir
		t.Fatal(err)
	}

	_, err := New(tmpFile, "/nonexistent/file.md")
	if err == nil {
		t.Fatal("New() should return error when any file does not exist")
	}
}

func TestWatchMultipleFiles(t *testing.T) {
	// Two files in the same directory and one in another directory
	dirA := t.TempDir()
	dirB := t.TempDir()
	fileA1 := filepath.Join(dirA, "a1.md")
	fileA2 := filepath.Join(dirA, "a2.md")
	fileB := filepath.Join(dirB, "b.md")
	for _, f := range []string{fileA1, fileA2, fileB} {
		if err := os.WriteFile(f, []byte("# Test"), 0644); err != nil { //nolint:gosec // G306: test file in temp dir
			t.Fatal(err)
		}
	}

	w, err := New(fileA1, fileA2, fileB)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	defer func() { _ = w.Close() }()

	w.Start()

	// Modify each file in turn and expect an event naming that file
	for _, want := range []string{fileA2, fileB, fileA1} {
		go func() {
			time.Sleep(100 * time.Millisecond)
			_ = os.WriteFile(want, []byte("# Updated"), 0600)
		}()

		select {
		case got := <-w.Events():
			if got != want {
				t.Errorf("Events() = %q, want %q", got, want)
			}
		case err := <-w.Errors():
			t.Fatalf("Errors() returned: %v", err)
		case <-time.After(2 * time.Second):
			t.Fatalf("timeout waiting for event for %s", want)
		}
	}
}

func TestWatchIgnoresUnwatchedFile(t *testing.T) {
	tmpDir := t.TempDir()
	watched := filepath.Join(tmpDir, "watched.md")
	other := filepath.Join(tmpDir, "other.md")
	for _, f := range []string{watched, other} {
		if err := os.WriteFile(f, []byte("# Test"), 0644); err != nil { //nolint:gosec // G306: test file in temp dir
			t.Fatal(err)
		}
	}

	w, err := New(watched)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	defer func() { _ = w.Close() }()

	w.Start()

	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = os.WriteFile(other, []byte("# Updated"), 0600)
	}()

	select {
	case got := <-w.Events():
		t.Fatalf("Events() returned %q for unwatched file change", got)
	case err := <-w.Errors():
		t.Fatalf("Errors() returned: %v", err)
	case <-time.After(500 * time.Millisecond):
		// Success: no event
	}
}

func TestWatchFileChange(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(tmpFile, []byte("# Test"), 0644); err != nil { //nolint:gosec // G306: test file in temp dir
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
	case got := <-w.Events():
		if got != tmpFile {
			t.Errorf("Events() = %q, want %q", got, tmpFile)
		}
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
	if err := os.WriteFile(tmpFile, []byte("# Test"), 0644); err != nil { //nolint:gosec // G306: test file in temp dir
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
	if err := os.WriteFile(tmpFile, []byte("# Test"), 0644); err != nil { //nolint:gosec // G306: test file in temp dir
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
