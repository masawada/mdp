// Package config handles configuration loading and defaults.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

var userHomeDir = os.UserHomeDir
var userConfigDir = os.UserConfigDir
var goos = runtime.GOOS

// DefaultOutputDir returns the default output directory path.
func DefaultOutputDir() string {
	homeDir, err := userHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".mdp")
}

// DefaultBrowserCommand returns the default browser command for the current OS.
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

func expandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~/") && path != "~" {
		return path, nil
	}
	home, err := userHomeDir()
	if err != nil {
		return "", err
	}
	return home + path[1:], nil
}

func configPathCandidates() []string {
	var candidates []string

	if configDir, err := userConfigDir(); err == nil {
		dir := filepath.Join(configDir, "mdp")
		candidates = append(candidates, filepath.Join(dir, "config.yaml"))
		candidates = append(candidates, filepath.Join(dir, "config.yml"))
	}

	if homeDir, err := userHomeDir(); err == nil {
		dir := filepath.Join(homeDir, ".config", "mdp")
		candidates = append(candidates, filepath.Join(dir, "config.yaml"))
		candidates = append(candidates, filepath.Join(dir, "config.yml"))
	}

	return candidates
}

func resolveConfigPath() string {
	for _, path := range configPathCandidates() {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// Config holds the application configuration.
type Config struct {
	OutputDir      string `yaml:"output_dir"`
	BrowserCommand string `yaml:"browser_command"`
	Theme          string `yaml:"theme"`
	ConfigDir      string `yaml:"-"`
}

// Load loads the configuration from the specified path or the default location.
func Load(path string) (*Config, error) {
	cfg := &Config{
		OutputDir:      DefaultOutputDir(),
		BrowserCommand: DefaultBrowserCommand(),
	}

	if path == "" {
		path = resolveConfigPath()
		if path == "" {
			// No config file found, use default ConfigDir
			cfg.ConfigDir = filepath.Dir(configPathCandidates()[0])
			return cfg, nil
		}
	}

	cfg.ConfigDir = filepath.Dir(path)

	data, err := os.ReadFile(path) //nolint:gosec // G304: path is user-specified config file
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
	} else {
		expanded, err := expandTilde(cfg.OutputDir)
		if err != nil {
			return nil, fmt.Errorf("failed to expand output_dir: %w", err)
		}
		cfg.OutputDir = expanded
	}
	if cfg.BrowserCommand == "" {
		cfg.BrowserCommand = DefaultBrowserCommand()
	}

	return cfg, nil
}
