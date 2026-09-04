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

	"github.com/masawada/mdp/internal/config"
	"github.com/masawada/mdp/internal/output"
	"github.com/masawada/mdp/internal/renderer"
)

func TestRun_FileNotFound(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter: &stdout,
		errWriter: &stderr,
	}

	cfg := &config.Config{}
	exitCode := c.run([]string{"/nonexistent/file.md"}, false, cfg)
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

	cfg, err := config.Load(configFile)
	if err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter: &stdout,
		errWriter: &stderr,
	}

	exitCode := c.run([]string{mdFile}, false, cfg)
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

// expectedOutputPath returns the path where run() writes the HTML for mdFile.
func expectedOutputPath(t *testing.T, outputDir, mdFile string) string {
	t.Helper()
	absPath, err := filepath.Abs(mdFile)
	if err != nil {
		t.Fatal(err)
	}
	pathWithoutExt := strings.TrimSuffix(absPath, ".md")
	relativePath := strings.TrimPrefix(pathWithoutExt, "/")
	return filepath.Join(outputDir, relativePath, "index.html")
}

// loadTestConfig writes a config pointing at outputDir with a no-op browser and loads it.
func loadTestConfig(t *testing.T, tmpDir, outputDir string) *config.Config {
	t.Helper()
	configFile := filepath.Join(tmpDir, "config.yaml")
	configContent := fmt.Sprintf("output_dir: %s\nbrowser_command: echo\n", outputDir)
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil { //nolint:gosec // G306: test file
		t.Fatal(err)
	}
	cfg, err := config.Load(configFile)
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestRun_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	mdFiles := []string{filepath.Join(tmpDir, "a.md"), filepath.Join(tmpDir, "b.md")}
	for _, f := range mdFiles {
		if err := os.WriteFile(f, []byte("# Hello"), 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}
	}

	outputDir := filepath.Join(tmpDir, "output")
	cfg := loadTestConfig(t, tmpDir, outputDir)

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter: &stdout,
		errWriter: &stderr,
	}

	exitCode := c.run(mdFiles, false, cfg)
	if exitCode != 0 {
		t.Errorf("run() exit code = %d, want 0\nstderr: %s", exitCode, stderr.String())
	}

	for _, f := range mdFiles {
		expectedHTML := expectedOutputPath(t, outputDir, f)
		if _, err := os.Stat(expectedHTML); os.IsNotExist(err) {
			t.Errorf("HTML file not created at %s", expectedHTML)
		}
		if !strings.Contains(stdout.String(), "Generated: "+expectedHTML+"\n") {
			t.Errorf("stdout should report %s, got: %s", expectedHTML, stdout.String())
		}
	}
	if got := strings.Count(stdout.String(), "Generated: "); got != 2 {
		t.Errorf("stdout should contain 2 Generated lines, got %d: %s", got, stdout.String())
	}
}

func TestRun_DuplicateFiles(t *testing.T) {
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "a.md")
	if err := os.WriteFile(mdFile, []byte("# Hello"), 0644); err != nil { //nolint:gosec // G306: test file
		t.Fatal(err)
	}

	outputDir := filepath.Join(tmpDir, "output")
	cfg := loadTestConfig(t, tmpDir, outputDir)

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter: &stdout,
		errWriter: &stderr,
	}

	exitCode := c.run([]string{mdFile, mdFile}, false, cfg)
	if exitCode != 0 {
		t.Errorf("run() exit code = %d, want 0\nstderr: %s", exitCode, stderr.String())
	}
	if got := strings.Count(stdout.String(), "Generated: "); got != 1 {
		t.Errorf("stdout should contain 1 Generated line, got %d: %s", got, stdout.String())
	}
}

func TestRun_OneFileNotFound_GeneratesNothing(t *testing.T) {
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "a.md")
	if err := os.WriteFile(mdFile, []byte("# Hello"), 0644); err != nil { //nolint:gosec // G306: test file
		t.Fatal(err)
	}
	missing := filepath.Join(tmpDir, "missing.md")

	outputDir := filepath.Join(tmpDir, "output")
	cfg := loadTestConfig(t, tmpDir, outputDir)

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter: &stdout,
		errWriter: &stderr,
	}

	exitCode := c.run([]string{mdFile, missing}, false, cfg)
	if exitCode != 1 {
		t.Errorf("run() exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "file not found: "+missing) {
		t.Errorf("stderr should name the missing file, got: %s", stderr.String())
	}
	if stdout.String() != "" {
		t.Errorf("stdout should be empty, got: %s", stdout.String())
	}
	if _, err := os.Stat(expectedOutputPath(t, outputDir, mdFile)); err == nil {
		t.Errorf("HTML for existing file should not be generated when another file is missing")
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

	cfg, err := config.Load(configFile)
	if err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter: &stdout,
		errWriter: &stderr,
	}

	exitCode := c.listFiles(cfg)
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

	cfg, err := config.Load(configFile)
	if err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter: &stdout,
		errWriter: &stderr,
	}

	exitCode := c.listFiles(cfg)
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

	cfg, err := config.Load(configFile)
	if err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := &cli{
		outWriter: &stdout,
		errWriter: &stderr,
	}

	exitCode := c.listFiles(cfg)
	if exitCode != 1 {
		t.Errorf("listFiles() exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "output directory does not exist") {
		t.Errorf("stderr should contain error message, got: %s", stderr.String())
	}
}

func TestConvert(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(mdFile, []byte("# Hello"), 0644); err != nil { //nolint:gosec // G306: test file in temp dir
		t.Fatal(err)
	}

	outDir := filepath.Join(tmpDir, "output")
	var outBuf, errBuf bytes.Buffer

	c := &cli{
		outWriter: &outBuf,
		errWriter: &errBuf,
	}

	r, err := renderer.NewRenderer("", "", renderer.Options{})
	if err != nil {
		t.Fatal(err)
	}
	w := output.NewWriter(outDir)

	outputPath, err := c.convert(mdFile, r, w)
	if err != nil {
		t.Fatalf("convert() returned error: %v", err)
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
	if err := os.WriteFile(mdFile, []byte("# Hello"), 0644); err != nil { //nolint:gosec // G306: test file in temp dir
		t.Fatal(err)
	}

	outDir := filepath.Join(tmpDir, "output")
	var outBuf, errBuf bytes.Buffer

	c := &cli{
		outWriter: &outBuf,
		errWriter: &errBuf,
	}

	r, _ := renderer.NewRenderer("", "", renderer.Options{})
	w := output.NewWriter(outDir)

	// Create channel for signal injection
	sigChan := make(chan os.Signal, 1)

	// Send signal in separate goroutine
	go func() {
		time.Sleep(100 * time.Millisecond)
		sigChan <- syscall.SIGINT
	}()

	exitCode := c.runWatchLoop([]string{mdFile}, r, w, sigChan)
	if exitCode != 0 {
		t.Errorf("runWatchLoop() returned %d, want 0", exitCode)
	}
}

func TestRunWatchLoop_RegeneratesChangedFile(t *testing.T) {
	tmpDir := t.TempDir()
	fileA := filepath.Join(tmpDir, "a.md")
	fileB := filepath.Join(tmpDir, "b.md")
	for _, f := range []string{fileA, fileB} {
		if err := os.WriteFile(f, []byte("# Hello"), 0644); err != nil { //nolint:gosec // G306: test file in temp dir
			t.Fatal(err)
		}
	}

	outDir := filepath.Join(tmpDir, "output")
	var outBuf, errBuf bytes.Buffer

	c := &cli{
		outWriter: &outBuf,
		errWriter: &errBuf,
	}

	r, _ := renderer.NewRenderer("", "", renderer.Options{})
	w := output.NewWriter(outDir)

	sigChan := make(chan os.Signal, 1)

	go func() {
		time.Sleep(200 * time.Millisecond)
		_ = os.WriteFile(fileB, []byte("# Updated"), 0600)
		time.Sleep(500 * time.Millisecond)
		sigChan <- syscall.SIGINT
	}()

	exitCode := c.runWatchLoop([]string{fileA, fileB}, r, w, sigChan)
	if exitCode != 0 {
		t.Errorf("runWatchLoop() returned %d, want 0", exitCode)
	}

	outputB := expectedOutputPath(t, outDir, fileB)
	if !strings.Contains(outBuf.String(), "Regenerated: "+outputB+"\n") {
		t.Errorf("stdout should report regeneration of %s, got: %s", outputB, outBuf.String())
	}
	if _, err := os.Stat(outputB); err != nil {
		t.Errorf("regenerated HTML not found: %v", err)
	}
	outputA := expectedOutputPath(t, outDir, fileA)
	if strings.Contains(outBuf.String(), outputA) {
		t.Errorf("unchanged file should not be regenerated, got: %s", outBuf.String())
	}
}
