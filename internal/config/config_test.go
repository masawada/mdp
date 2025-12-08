package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultOutputDir(t *testing.T) {
	t.Run("returns home/.mdp on success", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Fatal(err)
		}
		expected := filepath.Join(homeDir, ".mdp")
		actual := DefaultOutputDir()
		if actual != expected {
			t.Errorf("DefaultOutputDir() = %q, want %q", actual, expected)
		}
	})

	t.Run("panics when UserHomeDir fails", func(t *testing.T) {
		original := userHomeDir
		defer func() { userHomeDir = original }()

		userHomeDir = func() (string, error) {
			return "", errors.New("$HOME is not defined")
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("DefaultOutputDir() should panic when UserHomeDir fails")
			}
		}()

		DefaultOutputDir()
	})
}

func TestDefaultBrowserCommand(t *testing.T) {
	t.Run("returns open on darwin", func(t *testing.T) {
		original := goos
		defer func() { goos = original }()

		goos = "darwin"
		cmd := DefaultBrowserCommand()
		if cmd != "open" {
			t.Errorf("DefaultBrowserCommand() = %q, want \"open\"", cmd)
		}
	})

	t.Run("returns xdg-open on linux", func(t *testing.T) {
		original := goos
		defer func() { goos = original }()

		goos = "linux"
		cmd := DefaultBrowserCommand()
		if cmd != "xdg-open" {
			t.Errorf("DefaultBrowserCommand() = %q, want \"xdg-open\"", cmd)
		}
	})

	t.Run("panics on unsupported platform", func(t *testing.T) {
		original := goos
		defer func() { goos = original }()

		goos = "windows"

		defer func() {
			if r := recover(); r == nil {
				t.Error("DefaultBrowserCommand() should panic on unsupported platform")
			}
		}()

		DefaultBrowserCommand()
	})
}

func TestConfigPath(t *testing.T) {
	t.Run("returns UserConfigDir/mdp/config.yaml", func(t *testing.T) {
		configDir, _ := os.UserConfigDir()
		expected := filepath.Join(configDir, "mdp", "config.yaml")
		actual := configPath()
		if actual != expected {
			t.Errorf("configPath() = %q, want %q", actual, expected)
		}
	})

	t.Run("panics when UserConfigDir fails", func(t *testing.T) {
		original := userConfigDir
		defer func() { userConfigDir = original }()

		userConfigDir = func() (string, error) {
			return "", errors.New("UserConfigDir not available")
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("configPath() should panic when UserConfigDir fails")
			}
		}()

		configPath()
	})
}

func TestLoad(t *testing.T) {
	t.Run("file not found returns default", func(t *testing.T) {
		cfg, err := Load("/nonexistent/path/config.yaml")
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.OutputDir != DefaultOutputDir() {
			t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, DefaultOutputDir())
		}
		if cfg.BrowserCommand != DefaultBrowserCommand() {
			t.Errorf("BrowserCommand = %q, want %q", cfg.BrowserCommand, DefaultBrowserCommand())
		}
	})

	t.Run("empty path uses default config path", func(t *testing.T) {
		cfg, err := Load("")
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.OutputDir != DefaultOutputDir() {
			t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, DefaultOutputDir())
		}
	})

	t.Run("valid config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("output_dir: /custom/output\nbrowser_command: firefox\n")
		os.WriteFile(configFile, content, 0644)

		cfg, err := Load(configFile)
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.OutputDir != "/custom/output" {
			t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, "/custom/output")
		}
		if cfg.BrowserCommand != "firefox" {
			t.Errorf("BrowserCommand = %q, want %q", cfg.BrowserCommand, "firefox")
		}
	})

	t.Run("invalid yaml returns error", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("invalid: yaml: content:\n")
		os.WriteFile(configFile, content, 0644)

		_, err := Load(configFile)
		if err == nil {
			t.Error("Load() should return error for invalid yaml")
		}
	})

	t.Run("theme field is loaded correctly", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("theme: my-theme\n")
		os.WriteFile(configFile, content, 0644)

		cfg, err := Load(configFile)
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.Theme != "my-theme" {
			t.Errorf("Theme = %q, want %q", cfg.Theme, "my-theme")
		}
	})

	t.Run("theme field defaults to empty string when omitted", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("output_dir: /custom/output\n")
		os.WriteFile(configFile, content, 0644)

		cfg, err := Load(configFile)
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.Theme != "" {
			t.Errorf("Theme = %q, want empty string", cfg.Theme)
		}
	})

	t.Run("configDir is set to config file directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("output_dir: /custom/output\n")
		os.WriteFile(configFile, content, 0644)

		cfg, err := Load(configFile)
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.ConfigDir != tmpDir {
			t.Errorf("ConfigDir = %q, want %q", cfg.ConfigDir, tmpDir)
		}
	})

	t.Run("configDir is set to default config directory when path is empty", func(t *testing.T) {
		cfg, err := Load("")
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		configDir, _ := os.UserConfigDir()
		expected := filepath.Join(configDir, "mdp")
		if cfg.ConfigDir != expected {
			t.Errorf("ConfigDir = %q, want %q", cfg.ConfigDir, expected)
		}
	})
}
