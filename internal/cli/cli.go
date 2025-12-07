package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/masawada/mdp/internal/browser"
	"github.com/masawada/mdp/internal/config"
	"github.com/masawada/mdp/internal/converter"
	"github.com/masawada/mdp/internal/output"
)

type cli struct {
	outWriter, errWriter io.Writer
	configPath           string
}

func (c *cli) run(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(c.errWriter, "usage: mdp <markdown-file>")
		return 1
	}

	filePath := args[0]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Fprintf(c.errWriter, "error: file not found: %s\n", filePath)
		return 1
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Fprintf(c.errWriter, "error: %v\n", err)
		return 1
	}

	cfg, err := config.Load(c.configPath)
	if err != nil {
		fmt.Fprintf(c.errWriter, "error: failed to load config: %v\n", err)
		return 1
	}

	markdown, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Fprintf(c.errWriter, "error: failed to read file: %v\n", err)
		return 1
	}

	html, err := converter.Convert(markdown)
	if err != nil {
		fmt.Fprintf(c.errWriter, "error: failed to convert markdown: %v\n", err)
		return 1
	}

	writer := output.NewWriter(cfg.OutputDir)
	outputPath, err := writer.Write(absPath, html)
	if err != nil {
		fmt.Fprintf(c.errWriter, "error: failed to write html: %v\n", err)
		return 1
	}

	fmt.Fprintf(c.outWriter, "Generated: %s\n", outputPath)

	opener := browser.NewOpener(cfg.BrowserCommand)
	if err := opener.Open(outputPath); err != nil {
		fmt.Fprintf(c.errWriter, "error: failed to open browser: %v\n", err)
		return 1
	}

	return 0
}
