package output

import (
	"os"
	"path/filepath"
	"strings"
)

type Writer struct {
	baseDir string
}

func NewWriter(baseDir string) *Writer {
	return &Writer{baseDir: baseDir}
}

func (w *Writer) BuildOutputPath(srcPath string) string {
	ext := filepath.Ext(srcPath)
	pathWithoutExt := strings.TrimSuffix(srcPath, ext)
	relativePath := strings.TrimPrefix(pathWithoutExt, "/")
	return filepath.Join(w.baseDir, relativePath, "index.html")
}

func (w *Writer) Write(srcPath string, html []byte) (string, error) {
	outputPath := w.BuildOutputPath(srcPath)

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(outputPath, html, 0644); err != nil {
		return "", err
	}

	return outputPath, nil
}
