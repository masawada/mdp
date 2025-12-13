// Package main provides the entry point for the mdp command.
package main

import (
	"os"

	"github.com/masawada/mdp/internal/cli"
)

func main() {
	os.Exit(cli.Run())
}
