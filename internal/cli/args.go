// Package cli implements the command-line interface for mdp.
package cli

import (
	"errors"
	"flag"
	"io"
)

var errHelp = errors.New("help requested")

type parsedArgs struct {
	configPath string
	filePath   string
}

func parseArgs(args []string) (*parsedArgs, error) {
	fs := flag.NewFlagSet("mdp", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	configPath := fs.String("config", "", "path to config file")
	showHelp := fs.Bool("help", false, "show help message")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, errHelp
		}
		return nil, err
	}

	if *showHelp {
		return nil, errHelp
	}

	if fs.NArg() == 0 {
		return nil, errors.New("markdown file is required")
	}

	return &parsedArgs{
		configPath: *configPath,
		filePath:   fs.Arg(0),
	}, nil
}
