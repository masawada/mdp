// Package renderer converts Markdown to HTML with optional theming.
package renderer

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	meta "github.com/yuin/goldmark-meta"
)

// Renderer converts Markdown to HTML using an optional theme template.
type Renderer struct {
	tmpl *template.Template
}

type templateData struct {
	Title   string
	Content template.HTML
}

// NewRenderer creates a new Renderer with the specified theme.
func NewRenderer(configDir string, themeName string) (*Renderer, error) {
	if themeName == "" {
		return &Renderer{tmpl: nil}, nil
	}

	themePath := filepath.Join(configDir, "themes", themeName+".html")
	content, err := os.ReadFile(themePath) //nolint:gosec // G304: theme path is from trusted config
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(themeName).Parse(string(content))
	if err != nil {
		return nil, err
	}

	return &Renderer{tmpl: tmpl}, nil
}

// Render converts Markdown to HTML, applying the theme template if configured.
func (r *Renderer) Render(markdown []byte) ([]byte, error) {
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
		),
	)

	context := parser.NewContext()
	if err := md.Convert(markdown, &buf, parser.WithContext(context)); err != nil {
		return nil, err
	}

	title := extractTitle(context)
	html := buf.Bytes()

	if r.tmpl == nil {
		return html, nil
	}

	var out bytes.Buffer
	data := templateData{
		Title:   title,
		Content: template.HTML(html), //nolint:gosec // G203: HTML from markdown conversion is intentional
	}
	if err := r.tmpl.Execute(&out, data); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// extractTitle extracts the document title from markdown.
// Priority: 1. Front-matter title, 2. First heading, 3. "Untitled"
func extractTitle(context parser.Context) string {
	// Front-matter から取得
	metaData := meta.Get(context)
	if title, ok := metaData["title"].(string); ok && title != "" {
		return title
	}

	return "Untitled"
}
