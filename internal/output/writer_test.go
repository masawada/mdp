package output

import (
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
