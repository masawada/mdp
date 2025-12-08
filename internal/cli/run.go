package cli

import (
	"errors"
	"fmt"
	"os"
)

const usageMessage = `usage: mdp [options] <markdown-file>

Options:
  --config <config-file>  path to config file
  --help                  show this help message`

func Run() int {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		if errors.Is(err, errHelp) {
			fmt.Fprintln(os.Stdout, usageMessage)
			return 0
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintln(os.Stderr, "usage: mdp [--config <config-file>] [--help] <markdown-file>")
		return 1
	}

	return (&cli{
		outWriter:  os.Stdout,
		errWriter:  os.Stderr,
		configPath: args.configPath,
	}).run(args.filePath)
}
