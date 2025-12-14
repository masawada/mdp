package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/masawada/mdp/internal/browser"
	"github.com/masawada/mdp/internal/config"
	"github.com/masawada/mdp/internal/output"
	"github.com/masawada/mdp/internal/renderer"
)

type cli struct {
	outWriter, errWriter io.Writer
	configPath           string
}

func (c *cli) run(filePath string) int {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(c.errWriter, "error: file not found: %s\n", filePath)
		return 1
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "error: %v\n", err)
		return 1
	}

	cfg, err := config.Load(c.configPath)
	if err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "error: failed to load config: %v\n", err)
		return 1
	}

	r, err := renderer.NewRenderer(cfg.ConfigDir, cfg.Theme)
	if err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "error: failed to initialize renderer: %v\n", err)
		return 1
	}

	markdown, err := os.ReadFile(absPath) //nolint:gosec // G304: path is user-specified input file
	if err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "error: failed to read file: %v\n", err)
		return 1
	}

	html, err := r.Render(markdown)
	if err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "error: failed to render: %v\n", err)
		return 1
	}

	writer := output.NewWriter(cfg.OutputDir)
	outputPath, err := writer.Write(absPath, html)
	if err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "error: failed to write html: %v\n", err)
		return 1
	}

	_, _ = fmt.Fprintf(c.outWriter, "Generated: %s\n", outputPath)

	opener := browser.NewOpener(cfg.BrowserCommand)
	if err := opener.Open(outputPath); err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "error: failed to open browser: %v\n", err)
		return 1
	}

	return 0
}

// reconvert reads the markdown file, renders it, and writes the output
func (c *cli) reconvert(filePath string, r *renderer.Renderer, w *output.Writer) (string, error) {
	// Read file
	markdown, err := os.ReadFile(filePath) //nolint:gosec // G304: path is user-specified input file
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Render markdown to HTML
	html, err := r.Render(markdown)
	if err != nil {
		return "", fmt.Errorf("failed to render: %w", err)
	}

	// Write output
	outputPath, err := w.Write(filePath, html)
	if err != nil {
		return "", fmt.Errorf("failed to write: %w", err)
	}

	return outputPath, nil
}

func (c *cli) listFiles() int {
	cfg, err := config.Load(c.configPath)
	if err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "error: failed to load config: %v\n", err)
		return 1
	}

	files, err := output.ListFiles(cfg.OutputDir)
	if err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "error: %v\n", err)
		return 1
	}

	for _, file := range files {
		_, _ = fmt.Fprintln(c.outWriter, file)
	}

	return 0
}
