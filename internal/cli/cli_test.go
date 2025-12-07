package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_FileNotFound(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter: &stdout,
		errWriter: &stderr,
	}

	exitCode := c.run("/nonexistent/file.md")
	if exitCode != 1 {
		t.Errorf("run() exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "not found") && !strings.Contains(stderr.String(), "no such file") {
		t.Errorf("stderr should contain error message, got: %s", stderr.String())
	}
}

func TestRun_Success(t *testing.T) {
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(mdFile, []byte("# Hello"), 0644)

	outputDir := filepath.Join(tmpDir, "output")

	configFile := filepath.Join(tmpDir, "config.yaml")
	configContent := fmt.Sprintf("output_dir: %s\nbrowser_command: echo\n", outputDir)
	os.WriteFile(configFile, []byte(configContent), 0644)

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter:  &stdout,
		errWriter:  &stderr,
		configPath: configFile,
	}

	exitCode := c.run(mdFile)
	if exitCode != 0 {
		t.Errorf("run() exit code = %d, want 0\nstderr: %s", exitCode, stderr.String())
	}

	absPath, _ := filepath.Abs(mdFile)
	pathWithoutExt := strings.TrimSuffix(absPath, ".md")
	relativePath := strings.TrimPrefix(pathWithoutExt, "/")
	expectedHTML := filepath.Join(outputDir, relativePath, "index.html")
	if _, err := os.Stat(expectedHTML); os.IsNotExist(err) {
		t.Errorf("HTML file not created at %s", expectedHTML)
	}
}
