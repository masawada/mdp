package renderer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testTitleTemplate = `<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>{{.Content}}</body>
</html>`

func TestNewRenderer(t *testing.T) {
	t.Run("returns renderer without template when themeName is empty", func(t *testing.T) {
		r, err := NewRenderer("", "")
		if err != nil {
			t.Fatalf("NewRenderer() returned error: %v", err)
		}
		if r == nil {
			t.Fatal("NewRenderer() returned nil")
		}
	})

	t.Run("returns renderer with template when theme file exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		themesDir := filepath.Join(tmpDir, "themes")
		if err := os.MkdirAll(themesDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		themeFile := filepath.Join(themesDir, "test-theme.html")
		//nolint:gosec // G306: test file
		if err := os.WriteFile(themeFile, []byte("<html>{{.Content}}</html>"), 0644); err != nil {
			t.Fatal(err)
		}

		r, err := NewRenderer(tmpDir, "test-theme")
		if err != nil {
			t.Fatalf("NewRenderer() returned error: %v", err)
		}
		if r == nil {
			t.Fatal("NewRenderer() returned nil")
		}
	})

	t.Run("returns error when theme file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()

		_, err := NewRenderer(tmpDir, "nonexistent-theme")
		if err == nil {
			t.Error("NewRenderer() should return error when theme file does not exist")
		}
	})

	t.Run("returns error when template is invalid", func(t *testing.T) {
		tmpDir := t.TempDir()
		themesDir := filepath.Join(tmpDir, "themes")
		if err := os.MkdirAll(themesDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		themeFile := filepath.Join(themesDir, "invalid-theme.html")
		//nolint:gosec // G306: test file
		if err := os.WriteFile(themeFile, []byte("<html>{{.Invalid}</html>"), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := NewRenderer(tmpDir, "invalid-theme")
		if err == nil {
			t.Error("NewRenderer() should return error when template is invalid")
		}
	})
}

func TestRender(t *testing.T) {
	t.Run("converts markdown to html without theme", func(t *testing.T) {
		r, err := NewRenderer("", "")
		if err != nil {
			t.Fatalf("NewRenderer() returned error: %v", err)
		}

		html, err := r.Render([]byte("# Hello"))
		if err != nil {
			t.Fatalf("Render() returned error: %v", err)
		}

		expected := "<h1>Hello</h1>"
		if !strings.Contains(string(html), expected) {
			t.Errorf("Render() = %q, want to contain %q", string(html), expected)
		}
	})

	t.Run("applies theme template when theme is set", func(t *testing.T) {
		tmpDir := t.TempDir()
		themesDir := filepath.Join(tmpDir, "themes")
		if err := os.MkdirAll(themesDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		themeFile := filepath.Join(themesDir, "test-theme.html")
		//nolint:gosec // G306: test file
		if err := os.WriteFile(themeFile, []byte("<!DOCTYPE html><html><body>{{.Content}}</body></html>"), 0644); err != nil {
			t.Fatal(err)
		}

		r, err := NewRenderer(tmpDir, "test-theme")
		if err != nil {
			t.Fatalf("NewRenderer() returned error: %v", err)
		}

		html, err := r.Render([]byte("# Hello"))
		if err != nil {
			t.Fatalf("Render() returned error: %v", err)
		}

		result := string(html)
		if !strings.Contains(result, "<!DOCTYPE html>") {
			t.Errorf("Render() should contain DOCTYPE, got %q", result)
		}
		if !strings.Contains(result, "<h1>Hello</h1>") {
			t.Errorf("Render() should contain converted markdown, got %q", result)
		}
	})

	t.Run("does not escape html tags in content", func(t *testing.T) {
		tmpDir := t.TempDir()
		themesDir := filepath.Join(tmpDir, "themes")
		if err := os.MkdirAll(themesDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		themeFile := filepath.Join(themesDir, "test-theme.html")
		//nolint:gosec // G306: test file
		if err := os.WriteFile(themeFile, []byte("<div>{{.Content}}</div>"), 0644); err != nil {
			t.Fatal(err)
		}

		r, err := NewRenderer(tmpDir, "test-theme")
		if err != nil {
			t.Fatalf("NewRenderer() returned error: %v", err)
		}

		html, err := r.Render([]byte("**bold**"))
		if err != nil {
			t.Fatalf("Render() returned error: %v", err)
		}

		result := string(html)
		// HTMLタグがエスケープされていないことを確認
		if strings.Contains(result, "&lt;strong&gt;") {
			t.Errorf("Render() should not escape HTML tags, got %q", result)
		}
		if !strings.Contains(result, "<strong>bold</strong>") {
			t.Errorf("Render() should contain unescaped HTML tags, got %q", result)
		}
	})

	t.Run("extracts title from front-matter", func(t *testing.T) {
		tmpDir := t.TempDir()
		themesDir := filepath.Join(tmpDir, "themes")
		if err := os.MkdirAll(themesDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		themeFile := filepath.Join(themesDir, "test-theme.html")
		//nolint:gosec // G306: test file
		if err := os.WriteFile(themeFile, []byte(testTitleTemplate), 0644); err != nil {
			t.Fatal(err)
		}

		r, err := NewRenderer(tmpDir, "test-theme")
		if err != nil {
			t.Fatalf("NewRenderer() returned error: %v", err)
		}

		markdown := []byte(`---
title: Front-matter Title
---

# Heading Title

Content here.
`)

		html, err := r.Render(markdown)
		if err != nil {
			t.Fatalf("Render() returned error: %v", err)
		}

		result := string(html)
		if !strings.Contains(result, "<title>Front-matter Title</title>") {
			t.Errorf("Expected title 'Front-matter Title', got: %s", result)
		}
	})

	t.Run("extracts title from first heading when no front-matter", func(t *testing.T) {
		tmpDir := t.TempDir()
		themesDir := filepath.Join(tmpDir, "themes")
		if err := os.MkdirAll(themesDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		themeFile := filepath.Join(themesDir, "test-theme.html")
		//nolint:gosec // G306: test file
		if err := os.WriteFile(themeFile, []byte(testTitleTemplate), 0644); err != nil {
			t.Fatal(err)
		}

		r, err := NewRenderer(tmpDir, "test-theme")
		if err != nil {
			t.Fatalf("NewRenderer() returned error: %v", err)
		}

		markdown := []byte(`# First Heading

Some content.

## Second Heading

More content.
`)

		html, err := r.Render(markdown)
		if err != nil {
			t.Fatalf("Render() returned error: %v", err)
		}

		result := string(html)
		if !strings.Contains(result, "<title>First Heading</title>") {
			t.Errorf("Expected title 'First Heading', got: %s", result)
		}
	})

	t.Run("returns Untitled when no front-matter and no heading", func(t *testing.T) {
		tmpDir := t.TempDir()
		themesDir := filepath.Join(tmpDir, "themes")
		if err := os.MkdirAll(themesDir, 0755); err != nil { //nolint:gosec // G301: test directory
			t.Fatal(err)
		}
		themeFile := filepath.Join(themesDir, "test-theme.html")
		//nolint:gosec // G306: test file
		if err := os.WriteFile(themeFile, []byte(testTitleTemplate), 0644); err != nil {
			t.Fatal(err)
		}

		r, err := NewRenderer(tmpDir, "test-theme")
		if err != nil {
			t.Fatalf("NewRenderer() returned error: %v", err)
		}

		markdown := []byte(`Just some text without any heading.

More text here.
`)

		html, err := r.Render(markdown)
		if err != nil {
			t.Fatalf("Render() returned error: %v", err)
		}

		result := string(html)
		if !strings.Contains(result, "<title>Untitled</title>") {
			t.Errorf("Expected title 'Untitled', got: %s", result)
		}
	})
}
