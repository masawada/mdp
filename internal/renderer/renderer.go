package renderer

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Renderer struct {
	tmpl *template.Template
}

type templateData struct {
	Content template.HTML
}

func NewRenderer(configDir string, themeName string) (*Renderer, error) {
	if themeName == "" {
		return &Renderer{tmpl: nil}, nil
	}

	themePath := filepath.Join(configDir, "themes", themeName+".html")
	content, err := os.ReadFile(themePath)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(themeName).Parse(string(content))
	if err != nil {
		return nil, err
	}

	return &Renderer{tmpl: tmpl}, nil
}

func (r *Renderer) Render(markdown []byte) ([]byte, error) {
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
	if err := md.Convert(markdown, &buf); err != nil {
		return nil, err
	}

	html := buf.Bytes()

	if r.tmpl == nil {
		return html, nil
	}

	var out bytes.Buffer
	data := templateData{Content: template.HTML(html)}
	if err := r.tmpl.Execute(&out, data); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
