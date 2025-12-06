package cli

import (
	"fmt"
	"io"
)

type cli struct {
	outWriter, errWriter io.Writer
}

func (c *cli) run(args []string) int {
	fmt.Fprintf(c.outWriter, "")
	return 0
}
