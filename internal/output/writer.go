package output

import (
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
