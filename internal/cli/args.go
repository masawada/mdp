// Package cli implements the command-line interface for mdp.
package cli

import (
	"errors"
	"flag"
	"io"
)

// version is set via ldflags at build time.
var version = "dev"

var errHelp = errors.New("help requested")

type parsedArgs struct {
	configPath  string
	filePath    string
	showList    bool
	showVersion bool
	watchMode   bool
}

func parseArgs(args []string) (*parsedArgs, error) {
	fs := flag.NewFlagSet("mdp", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	configPath := fs.String("config", "", "path to config file")
	showHelp := fs.Bool("help", false, "show help message")
	showList := fs.Bool("list", false, "list generated files")
	showVersion := fs.Bool("version", false, "show version")
	watchMode := fs.Bool("watch", false, "watch for file changes")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, errHelp
		}
		return nil, err
	}

	if *showHelp {
		return nil, errHelp
	}

	if *showVersion {
		return &parsedArgs{
			showVersion: true,
		}, nil
	}

	if *showList {
		return &parsedArgs{
			configPath: *configPath,
			showList:   true,
		}, nil
	}

	if fs.NArg() == 0 {
		return nil, errors.New("markdown file is required")
	}

	return &parsedArgs{
		configPath: *configPath,
		filePath:   fs.Arg(0),
		watchMode:  *watchMode,
	}, nil
}
