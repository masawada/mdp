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
		actual, err := DefaultOutputDir()
		if err != nil {
			t.Fatalf("DefaultOutputDir() returned error: %v", err)
		}
		if actual != expected {
			t.Errorf("DefaultOutputDir() = %q, want %q", actual, expected)
		}
	})

	t.Run("returns error when UserHomeDir fails", func(t *testing.T) {
		original := userHomeDir
		defer func() { userHomeDir = original }()

		userHomeDir = func() (string, error) {
			return "", errors.New("$HOME is not defined")
		}

		_, err := DefaultOutputDir()
		if err == nil {
			t.Error("DefaultOutputDir() should return error when UserHomeDir fails")
		}
	})
}

func TestDefaultBrowserCommand(t *testing.T) {
	t.Run("returns open on darwin", func(t *testing.T) {
		original := goos
		defer func() { goos = original }()

		goos = "darwin"
		cmd, err := DefaultBrowserCommand()
		if err != nil {
			t.Fatalf("DefaultBrowserCommand() returned error: %v", err)
		}
		if cmd != "open" {
			t.Errorf("DefaultBrowserCommand() = %q, want \"open\"", cmd)
		}
	})

	t.Run("returns xdg-open on linux", func(t *testing.T) {
		original := goos
		defer func() { goos = original }()

		goos = "linux"
		cmd, err := DefaultBrowserCommand()
		if err != nil {
			t.Fatalf("DefaultBrowserCommand() returned error: %v", err)
		}
		if cmd != "xdg-open" {
			t.Errorf("DefaultBrowserCommand() = %q, want \"xdg-open\"", cmd)
		}
	})

	t.Run("returns error on unsupported platform", func(t *testing.T) {
		original := goos
		defer func() { goos = original }()

		goos = "windows"

		_, err := DefaultBrowserCommand()
		if err == nil {
			t.Error("DefaultBrowserCommand() should return error on unsupported platform")
		}
	})
}

func TestConfigPathCandidates(t *testing.T) {
	t.Run("returns 4 candidate paths in correct order", func(t *testing.T) {
		configDir, _ := os.UserConfigDir()
		homeDir, _ := os.UserHomeDir()

		candidates := configPathCandidates()

		expected := []string{
			filepath.Join(configDir, "mdp", "config.yaml"),
			filepath.Join(configDir, "mdp", "config.yml"),
			filepath.Join(homeDir, ".config", "mdp", "config.yaml"),
			filepath.Join(homeDir, ".config", "mdp", "config.yml"),
		}

		if len(candidates) != len(expected) {
			t.Fatalf("configPathCandidates() returned %d candidates, want %d", len(candidates), len(expected))
		}

		for i, path := range candidates {
			if path != expected[i] {
				t.Errorf("configPathCandidates()[%d] = %q, want %q", i, path, expected[i])
			}
		}
	})
}

func TestResolveConfigPath(t *testing.T) {
	t.Run("returns first candidate (config.yaml) when it exists", func(t *testing.T) {
		tmpDir := t.TempDir()

		originalConfigDir := userConfigDir
		defer func() { userConfigDir = originalConfigDir }()
		userConfigDir = func() (string, error) { return tmpDir, nil }

		configDir := filepath.Join(tmpDir, "mdp")
		if err := os.MkdirAll(configDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		configFile := filepath.Join(configDir, "config.yaml")
		if err := os.WriteFile(configFile, []byte(""), 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

		result := resolveConfigPath()
		if result != configFile {
			t.Errorf("resolveConfigPath() = %q, want %q", result, configFile)
		}
	})

	t.Run("returns config.yml when config.yaml does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()

		originalConfigDir := userConfigDir
		defer func() { userConfigDir = originalConfigDir }()
		userConfigDir = func() (string, error) { return tmpDir, nil }

		configDir := filepath.Join(tmpDir, "mdp")
		if err := os.MkdirAll(configDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		configFile := filepath.Join(configDir, "config.yml")
		if err := os.WriteFile(configFile, []byte(""), 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

		result := resolveConfigPath()
		if result != configFile {
			t.Errorf("resolveConfigPath() = %q, want %q", result, configFile)
		}
	})

	t.Run("returns fallback path when primary does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		primaryDir := filepath.Join(tmpDir, "primary")
		fallbackDir := filepath.Join(tmpDir, "fallback")

		originalConfigDir := userConfigDir
		originalHomeDir := userHomeDir
		defer func() {
			userConfigDir = originalConfigDir
			userHomeDir = originalHomeDir
		}()
		userConfigDir = func() (string, error) { return primaryDir, nil }
		userHomeDir = func() (string, error) { return fallbackDir, nil }

		// Create config only in fallback location
		configDir := filepath.Join(fallbackDir, ".config", "mdp")
		if err := os.MkdirAll(configDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		configFile := filepath.Join(configDir, "config.yaml")
		if err := os.WriteFile(configFile, []byte(""), 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

		result := resolveConfigPath()
		if result != configFile {
			t.Errorf("resolveConfigPath() = %q, want %q", result, configFile)
		}
	})

	t.Run("returns empty string when no config file exists", func(t *testing.T) {
		tmpDir := t.TempDir()

		originalConfigDir := userConfigDir
		originalHomeDir := userHomeDir
		defer func() {
			userConfigDir = originalConfigDir
			userHomeDir = originalHomeDir
		}()
		userConfigDir = func() (string, error) { return filepath.Join(tmpDir, "config"), nil }
		userHomeDir = func() (string, error) { return filepath.Join(tmpDir, "home"), nil }

		result := resolveConfigPath()
		if result != "" {
			t.Errorf("resolveConfigPath() = %q, want empty string", result)
		}
	})
}

func TestLoad(t *testing.T) {
	t.Run("file not found returns default", func(t *testing.T) {
		cfg, err := Load("/nonexistent/path/config.yaml")
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		expectedOutputDir, _ := DefaultOutputDir()
		if cfg.OutputDir != expectedOutputDir {
			t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, expectedOutputDir)
		}
		expectedBrowserCommand, _ := DefaultBrowserCommand()
		if cfg.BrowserCommand != expectedBrowserCommand {
			t.Errorf("BrowserCommand = %q, want %q", cfg.BrowserCommand, expectedBrowserCommand)
		}
	})

	t.Run("empty path uses default config path", func(t *testing.T) {
		cfg, err := Load("")
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		expectedOutputDir, _ := DefaultOutputDir()
		if cfg.OutputDir != expectedOutputDir {
			t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, expectedOutputDir)
		}
	})

	t.Run("valid config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("output_dir: /custom/output\nbrowser_command: firefox\n")
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

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
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

		_, err := Load(configFile)
		if err == nil {
			t.Error("Load() should return error for invalid yaml")
		}
	})

	t.Run("theme field is loaded correctly", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("theme: my-theme\n")
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

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
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

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
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

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

	t.Run("loads config.yml when config.yaml does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()

		originalConfigDir := userConfigDir
		defer func() { userConfigDir = originalConfigDir }()
		userConfigDir = func() (string, error) { return tmpDir, nil }

		configDir := filepath.Join(tmpDir, "mdp")
		if err := os.MkdirAll(configDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		configFile := filepath.Join(configDir, "config.yml")
		content := []byte("output_dir: /yml/output\n")
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

		cfg, err := Load("")
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.OutputDir != "/yml/output" {
			t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, "/yml/output")
		}
		if cfg.ConfigDir != configDir {
			t.Errorf("ConfigDir = %q, want %q", cfg.ConfigDir, configDir)
		}
	})

	t.Run("loads config from fallback path when primary does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		primaryDir := filepath.Join(tmpDir, "primary")
		fallbackDir := filepath.Join(tmpDir, "fallback")

		originalConfigDir := userConfigDir
		originalHomeDir := userHomeDir
		defer func() {
			userConfigDir = originalConfigDir
			userHomeDir = originalHomeDir
		}()
		userConfigDir = func() (string, error) { return primaryDir, nil }
		userHomeDir = func() (string, error) { return fallbackDir, nil }

		// Create config only in fallback location
		configDir := filepath.Join(fallbackDir, ".config", "mdp")
		if err := os.MkdirAll(configDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		configFile := filepath.Join(configDir, "config.yaml")
		content := []byte("output_dir: /fallback/output\n")
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

		cfg, err := Load("")
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.OutputDir != "/fallback/output" {
			t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, "/fallback/output")
		}
		if cfg.ConfigDir != configDir {
			t.Errorf("ConfigDir = %q, want %q", cfg.ConfigDir, configDir)
		}
	})

	t.Run("uses first candidate dir for ConfigDir when no config file exists", func(t *testing.T) {
		tmpDir := t.TempDir()

		originalConfigDir := userConfigDir
		originalHomeDir := userHomeDir
		defer func() {
			userConfigDir = originalConfigDir
			userHomeDir = originalHomeDir
		}()
		userConfigDir = func() (string, error) { return filepath.Join(tmpDir, "config"), nil }
		userHomeDir = func() (string, error) { return filepath.Join(tmpDir, "home"), nil }

		cfg, err := Load("")
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		expected := filepath.Join(tmpDir, "config", "mdp")
		if cfg.ConfigDir != expected {
			t.Errorf("ConfigDir = %q, want %q", cfg.ConfigDir, expected)
		}
	})

	t.Run("hard_wraps field is loaded correctly", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("hard_wraps: true\n")
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

		cfg, err := Load(configFile)
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if !cfg.HardWraps {
			t.Error("HardWraps should be true")
		}
	})

	t.Run("hard_wraps defaults to false when omitted", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("output_dir: /custom/output\n")
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

		cfg, err := Load(configFile)
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.HardWraps {
			t.Error("HardWraps should default to false")
		}
	})

	t.Run("unsafe field is loaded correctly", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("unsafe: true\n")
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

		cfg, err := Load(configFile)
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if !cfg.Unsafe {
			t.Error("Unsafe should be true")
		}
	})

	t.Run("unsafe defaults to false when omitted", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("output_dir: /custom/output\n")
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

		cfg, err := Load(configFile)
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		if cfg.Unsafe {
			t.Error("Unsafe should default to false")
		}
	})

	t.Run("expands tilde in output_dir", func(t *testing.T) {
		tmpDir := t.TempDir()
		homeDir := filepath.Join(tmpDir, "home")

		originalHomeDir := userHomeDir
		defer func() { userHomeDir = originalHomeDir }()
		userHomeDir = func() (string, error) { return homeDir, nil }

		configFile := filepath.Join(tmpDir, "config.yaml")
		content := []byte("output_dir: ~/.mdp-custom\n")
		if err := os.WriteFile(configFile, content, 0644); err != nil { //nolint:gosec // G306: test file
			t.Fatal(err)
		}

		cfg, err := Load(configFile)
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}
		expected := filepath.Join(homeDir, ".mdp-custom")
		if cfg.OutputDir != expected {
			t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, expected)
		}
	})
}

func TestExpandTilde(t *testing.T) {
	t.Run("returns empty string unchanged", func(t *testing.T) {
		result, err := expandTilde("")
		if err != nil {
			t.Fatalf("expandTilde() returned error: %v", err)
		}
		if result != "" {
			t.Errorf("expandTilde(\"\") = %q, want \"\"", result)
		}
	})

	t.Run("returns path without tilde unchanged", func(t *testing.T) {
		result, err := expandTilde("/some/absolute/path")
		if err != nil {
			t.Fatalf("expandTilde() returned error: %v", err)
		}
		if result != "/some/absolute/path" {
			t.Errorf("expandTilde(\"/some/absolute/path\") = %q, want \"/some/absolute/path\"", result)
		}
	})

	t.Run("expands ~ to home directory", func(t *testing.T) {
		homeDir := "/home/testuser"
		originalHomeDir := userHomeDir
		defer func() { userHomeDir = originalHomeDir }()
		userHomeDir = func() (string, error) { return homeDir, nil }

		result, err := expandTilde("~")
		if err != nil {
			t.Fatalf("expandTilde() returned error: %v", err)
		}
		if result != homeDir {
			t.Errorf("expandTilde(\"~\") = %q, want %q", result, homeDir)
		}
	})

	t.Run("expands ~/ to home directory with path", func(t *testing.T) {
		homeDir := "/home/testuser"
		originalHomeDir := userHomeDir
		defer func() { userHomeDir = originalHomeDir }()
		userHomeDir = func() (string, error) { return homeDir, nil }

		result, err := expandTilde("~/.config/mdp")
		if err != nil {
			t.Fatalf("expandTilde() returned error: %v", err)
		}
		expected := homeDir + "/.config/mdp"
		if result != expected {
			t.Errorf("expandTilde(\"~/.config/mdp\") = %q, want %q", result, expected)
		}
	})

	t.Run("returns error when userHomeDir fails", func(t *testing.T) {
		originalHomeDir := userHomeDir
		defer func() { userHomeDir = originalHomeDir }()
		userHomeDir = func() (string, error) { return "", errors.New("home dir not found") }

		_, err := expandTilde("~/.config")
		if err == nil {
			t.Error("expandTilde() should return error when userHomeDir fails")
		}
	})

	t.Run("does not expand ~username", func(t *testing.T) {
		result, err := expandTilde("~otheruser/.config")
		if err != nil {
			t.Fatalf("expandTilde() returned error: %v", err)
		}
		if result != "~otheruser/.config" {
			t.Errorf("expandTilde(\"~otheruser/.config\") = %q, want \"~otheruser/.config\"", result)
		}
	})
}
