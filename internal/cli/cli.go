package cli

import (
	"fmt"
	"io"
)

type cli struct {
	outWriter, errWriter io.Writer
}

func (c *cli) run(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(c.errWriter, "usage: mdp <markdown-file>")
		return 1
	}
	return 0
}
