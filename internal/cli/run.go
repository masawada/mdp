package cli

import (
	"fmt"
	"os"
)

func Run() int {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintln(os.Stderr, "usage: mdp [--config <config-file>] <markdown-file>")
		return 1
	}

	return (&cli{
		outWriter:  os.Stdout,
		errWriter:  os.Stderr,
		configPath: args.configPath,
	}).run(args.filePath)
}
