package cli

import (
	"fmt"
	"io"
	"os"
)

type cli struct {
	outWriter, errWriter io.Writer
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

	return 0
}
