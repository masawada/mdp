package output

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestBuildOutputPath(t *testing.T) {
	w := NewWriter("/base/dir")

	tests := []struct {
		name     string
		srcPath  string
		expected string
	}{
		{
			name:     "absolute path with .md extension",
			srcPath:  "/Users/user/docs/readme.md",
			expected: "/base/dir/Users/user/docs/readme/index.html",
		},
		{
			name:     "absolute path with .markdown extension",
			srcPath:  "/Users/user/docs/guide.markdown",
			expected: "/base/dir/Users/user/docs/guide/index.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := w.BuildOutputPath(tt.srcPath)
			if actual != tt.expected {
				t.Errorf("BuildOutputPath(%q) = %q, want %q", tt.srcPath, actual, tt.expected)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tmpDir := t.TempDir()
	w := NewWriter(tmpDir)

	srcPath := "/Users/user/docs/readme.md"
	htmlContent := []byte("<h1>Hello</h1>")

	outputPath, err := w.Write(srcPath, htmlContent)
	if err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "Users/user/docs/readme/index.html")
	if outputPath != expectedPath {
		t.Errorf("Write() returned path = %q, want %q", outputPath, expectedPath)
	}

	content, err := os.ReadFile(outputPath) //nolint:gosec // G304: path is from test
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if string(content) != string(htmlContent) {
		t.Errorf("File content = %q, want %q", content, htmlContent)
	}
}

func TestListFiles_DirectoryNotExist(t *testing.T) {
	_, err := ListFiles("/non/existent/directory")
	if err == nil {
		t.Error("ListFiles() expected error for non-existent directory, got nil")
	}
}

func TestListFiles_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	files, err := ListFiles(tmpDir)
	if err != nil {
		t.Fatalf("ListFiles() error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("ListFiles() returned %d files, want 0", len(files))
	}
}

func TestListFiles_WithFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	dir1 := filepath.Join(tmpDir, "path/to/file1")
	dir2 := filepath.Join(tmpDir, "path/to/file2")
	if err := os.MkdirAll(dir1, 0755); err != nil { //nolint:gosec // G301: test directory
		t.Fatalf("Failed to create dir: %v", err)
	}
	if err := os.MkdirAll(dir2, 0755); err != nil { //nolint:gosec // G301: test directory
		t.Fatalf("Failed to create dir: %v", err)
	}

	file1 := filepath.Join(dir1, "index.html")
	file2 := filepath.Join(dir2, "index.html")
	if err := os.WriteFile(file1, []byte("<h1>Test1</h1>"), 0644); err != nil { //nolint:gosec // G306: test file
		t.Fatalf("Failed to write file: %v", err)
	}
	if err := os.WriteFile(file2, []byte("<h1>Test2</h1>"), 0644); err != nil { //nolint:gosec // G306: test file
		t.Fatalf("Failed to write file: %v", err)
	}

	files, err := ListFiles(tmpDir)
	if err != nil {
		t.Fatalf("ListFiles() error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("ListFiles() returned %d files, want 2", len(files))
	}

	// Sort for consistent comparison
	sort.Strings(files)
	expected := []string{file1, file2}
	sort.Strings(expected)

	for i, f := range files {
		if f != expected[i] {
			t.Errorf("ListFiles()[%d] = %q, want %q", i, f, expected[i])
		}
	}
}
