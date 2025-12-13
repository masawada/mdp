// Package output handles writing rendered HTML to the filesystem.
package output

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Writer writes rendered HTML files to a base directory.
type Writer struct {
	baseDir string
}

// NewWriter creates a new Writer with the specified base directory.
func NewWriter(baseDir string) *Writer {
	return &Writer{baseDir: baseDir}
}

// BuildOutputPath constructs the output path for a given source file path.
func (w *Writer) BuildOutputPath(srcPath string) string {
	ext := filepath.Ext(srcPath)
	pathWithoutExt := strings.TrimSuffix(srcPath, ext)
	relativePath := strings.TrimPrefix(pathWithoutExt, "/")
	return filepath.Join(w.baseDir, relativePath, "index.html")
}

func (w *Writer) Write(srcPath string, html []byte) (string, error) {
	outputPath := w.BuildOutputPath(srcPath)

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil { //nolint:gosec // G301: need world-readable for browser
		return "", err
	}

	if err := os.WriteFile(outputPath, html, 0644); err != nil { //nolint:gosec // G306: need world-readable for browser
		return "", err
	}

	return outputPath, nil
}

// ListFiles returns a list of generated HTML files in the specified directory.
func ListFiles(baseDir string) ([]string, error) {
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("output directory does not exist: %s", baseDir)
	}

	var files []string
	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "index.html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
