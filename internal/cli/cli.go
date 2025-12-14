package cli

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/masawada/mdp/internal/browser"
	"github.com/masawada/mdp/internal/config"
	"github.com/masawada/mdp/internal/output"
	"github.com/masawada/mdp/internal/renderer"
	"github.com/masawada/mdp/internal/watcher"
)

type cli struct {
	outWriter, errWriter io.Writer
	configPath           string
}

func (c *cli) run(filePath string, watchMode bool) int {
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

	// If watch mode is enabled, start the watch loop
	if watchMode {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		return c.runWatchLoop(absPath, r, writer, sigChan)
	}

	return 0
}

// runWatchLoop watches for file changes and regenerates HTML
func (c *cli) runWatchLoop(filePath string, r *renderer.Renderer, w *output.Writer, sigChan <-chan os.Signal) int {
	// Create watcher
	fileWatcher, err := watcher.New(filePath)
	if err != nil {
		_, _ = fmt.Fprintf(c.errWriter, "error: failed to start watcher: %v\n", err)
		return 1
	}
	defer fileWatcher.Close()

	fileWatcher.Start()
	_, _ = fmt.Fprintln(c.outWriter, "Watching for changes... (Ctrl+C to stop)")

	for {
		select {
		case <-fileWatcher.Events():
			outputPath, err := c.reconvert(filePath, r, w)
			if err != nil {
				_, _ = fmt.Fprintf(c.errWriter, "error: %v\n", err)
				continue
			}
			_, _ = fmt.Fprintf(c.outWriter, "Regenerated: %s\n", outputPath)
		case err := <-fileWatcher.Errors():
			_, _ = fmt.Fprintf(c.errWriter, "watcher error: %v\n", err)
		case <-sigChan:
			_, _ = fmt.Fprintln(c.outWriter, "\nStopping watcher...")
			return 0
		}
	}
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
