package output

import (
	"os"
	"path/filepath"
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
