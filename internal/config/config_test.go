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
