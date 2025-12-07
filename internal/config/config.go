package config

import (
	"os"
	"path/filepath"
)

var userHomeDir = os.UserHomeDir

func DefaultOutputDir() string {
	homeDir, err := userHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".mdp")
}
