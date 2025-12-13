package cli

import (
	"errors"
	"fmt"
	"os"
)

const usageMessage = `usage: mdp [options] <markdown-file>

Options:
  --config <config-file>  path to config file
  --list                  list generated files
  --version               show version
  --help                  show this help message`

// Run executes the mdp command and returns the exit code.
func Run() int {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		if errors.Is(err, errHelp) {
			_, _ = fmt.Fprintln(os.Stdout, usageMessage)
			return 0
		}
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		_, _ = fmt.Fprintln(os.Stderr, "usage: mdp [--config <config-file>] [--help] <markdown-file>")
		return 1
	}

	if args.showVersion {
		_, _ = fmt.Fprintf(os.Stdout, "mdp version %s\n", version)
		return 0
	}

	c := &cli{
		outWriter:  os.Stdout,
		errWriter:  os.Stderr,
		configPath: args.configPath,
	}

	if args.showList {
		return c.listFiles()
	}

	return c.run(args.filePath)
}
