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
