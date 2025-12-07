package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
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

type Config struct {
	OutputDir      string `yaml:"output_dir"`
	BrowserCommand string `yaml:"browser_command"`
}

func Load() (*Config, error) {
	cfg := &Config{
		OutputDir:      DefaultOutputDir(),
		BrowserCommand: DefaultBrowserCommand(),
	}

	data, err := os.ReadFile(configPath())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.OutputDir == "" {
		cfg.OutputDir = DefaultOutputDir()
	}
	if cfg.BrowserCommand == "" {
		cfg.BrowserCommand = DefaultBrowserCommand()
	}

	return cfg, nil
}
