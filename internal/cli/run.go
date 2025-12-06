package cli

import "os"

func Run() int {
	return (&cli{
		outWriter: os.Stdout,
		errWriter: os.Stderr,
	}).run(os.Args[1:])
}
