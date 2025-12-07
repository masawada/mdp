package cli

import (
	"errors"
	"flag"
	"io"
)

type parsedArgs struct {
	configPath string
	filePath   string
}

func parseArgs(args []string) (*parsedArgs, error) {
	fs := flag.NewFlagSet("mdp", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	configPath := fs.String("config", "", "path to config file")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() == 0 {
		return nil, errors.New("markdown file is required")
	}

	return &parsedArgs{
		configPath: *configPath,
		filePath:   fs.Arg(0),
	}, nil
}
