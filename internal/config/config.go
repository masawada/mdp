package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var userHomeDir = os.UserHomeDir
var userConfigDir = os.UserConfigDir
var goos = runtime.GOOS

func DefaultOutputDir() string {
	homeDir, err := userHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".mdp")
}

func DefaultBrowserCommand() string {
	switch goos {
	case "darwin":
		return "open"
	case "linux":
		return "xdg-open"
	default:
		panic(fmt.Sprintf("unsupported platform: %s", goos))
	}
}

func configPath() string {
	configDir, err := userConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(configDir, "mdp", "config.yaml")
}
