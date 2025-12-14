package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/masawada/mdp/internal/output"
	"github.com/masawada/mdp/internal/renderer"
)

func TestRun_FileNotFound(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter: &stdout,
		errWriter: &stderr,
	}

	exitCode := c.run("/nonexistent/file.md", false)
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
	if err := os.WriteFile(mdFile, []byte("# Hello"), 0644); err != nil { //nolint:gosec // G306: test file
		t.Fatal(err)
	}

	outputDir := filepath.Join(tmpDir, "output")

	configFile := filepath.Join(tmpDir, "config.yaml")
	configContent := fmt.Sprintf("output_dir: %s\nbrowser_command: echo\n", outputDir)
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil { //nolint:gosec // G306: test file
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter:  &stdout,
		errWriter:  &stderr,
		configPath: configFile,
	}

	exitCode := c.run(mdFile, false)
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

func TestListFiles_WithFiles(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")

	// Create test files
	dir1 := filepath.Join(outputDir, "path/to/file1")
	if err := os.MkdirAll(dir1, 0755); err != nil { //nolint:gosec // G301: test directory
		t.Fatal(err)
	}
	file1 := filepath.Join(dir1, "index.html")
	if err := os.WriteFile(file1, []byte("<h1>Test</h1>"), 0644); err != nil { //nolint:gosec // G306: test file
		t.Fatal(err)
	}

	configFile := filepath.Join(tmpDir, "config.yaml")
	configContent := fmt.Sprintf("output_dir: %s\n", outputDir)
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil { //nolint:gosec // G306: test file
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter:  &stdout,
		errWriter:  &stderr,
		configPath: configFile,
	}

	exitCode := c.listFiles()
	if exitCode != 0 {
		t.Errorf("listFiles() exit code = %d, want 0\nstderr: %s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), file1) {
		t.Errorf("stdout should contain file path, got: %s", stdout.String())
	}
}

func TestListFiles_NoFiles(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil { //nolint:gosec // G301: test directory
		t.Fatal(err)
	}

	configFile := filepath.Join(tmpDir, "config.yaml")
	configContent := fmt.Sprintf("output_dir: %s\n", outputDir)
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil { //nolint:gosec // G306: test file
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter:  &stdout,
		errWriter:  &stderr,
		configPath: configFile,
	}

	exitCode := c.listFiles()
	if exitCode != 0 {
		t.Errorf("listFiles() exit code = %d, want 0", exitCode)
	}
	if stdout.String() != "" {
		t.Errorf("stdout should be empty, got: %s", stdout.String())
	}
}

func TestListFiles_DirectoryNotExist(t *testing.T) {
	tmpDir := t.TempDir()

	configFile := filepath.Join(tmpDir, "config.yaml")
	configContent := "output_dir: /nonexistent/directory\n"
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil { //nolint:gosec // G306: test file
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter:  &stdout,
		errWriter:  &stderr,
		configPath: configFile,
	}

	exitCode := c.listFiles()
	if exitCode != 1 {
		t.Errorf("listFiles() exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "output directory does not exist") {
		t.Errorf("stderr should contain error message, got: %s", stderr.String())
	}
}

func TestReconvert(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(mdFile, []byte("# Hello"), 0644); err != nil {
		t.Fatal(err)
	}

	outDir := filepath.Join(tmpDir, "output")
	var outBuf, errBuf bytes.Buffer

	c := &cli{
		outWriter: &outBuf,
		errWriter: &errBuf,
	}

	r, err := renderer.NewRenderer("", "")
	if err != nil {
		t.Fatal(err)
	}
	w := output.NewWriter(outDir)

	outputPath, err := c.reconvert(mdFile, r, w)
	if err != nil {
		t.Fatalf("reconvert() returned error: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("output file not found: %v", err)
	}
}

func TestRunWatchLoop_SignalHandling(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(mdFile, []byte("# Hello"), 0644); err != nil {
		t.Fatal(err)
	}

	outDir := filepath.Join(tmpDir, "output")
	var outBuf, errBuf bytes.Buffer

	c := &cli{
		outWriter: &outBuf,
		errWriter: &errBuf,
	}

	r, _ := renderer.NewRenderer("", "")
	w := output.NewWriter(outDir)

	// Create channel for signal injection
	sigChan := make(chan os.Signal, 1)

	// Send signal in separate goroutine
	go func() {
		time.Sleep(100 * time.Millisecond)
		sigChan <- syscall.SIGINT
	}()

	exitCode := c.runWatchLoop(mdFile, r, w, sigChan)
	if exitCode != 0 {
		t.Errorf("runWatchLoop() returned %d, want 0", exitCode)
	}
}
