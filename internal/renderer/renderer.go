// Package renderer converts Markdown to HTML with optional theming.
package renderer

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// Renderer converts Markdown to HTML using an optional theme template.
type Renderer struct {
	tmpl *template.Template
	md   goldmark.Markdown
}

type templateData struct {
	Title   string
	Content template.HTML
}

// Options holds optional settings for the Renderer.
type Options struct {
	HardWraps bool
	Unsafe    bool
}

// NewRenderer creates a new Renderer with the specified theme.
func NewRenderer(configDir string, themeName string, opts Options) (*Renderer, error) {
	gmOpts := []goldmark.Option{
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			meta.Meta,
		),
	}
	if opts.HardWraps {
		gmOpts = append(gmOpts, goldmark.WithRendererOptions(html.WithHardWraps()))
	}
	if opts.Unsafe {
		gmOpts = append(gmOpts, goldmark.WithRendererOptions(html.WithUnsafe()))
	}
	md := goldmark.New(gmOpts...)

	if themeName == "" {
		return &Renderer{tmpl: nil, md: md}, nil
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

	return &Renderer{tmpl: tmpl, md: md}, nil
}

// Render converts Markdown to HTML, applying the theme template if configured.
// Relative image paths in the content resolve against baseDir,
// the directory of the source markdown file.
func (r *Renderer) Render(markdown []byte, baseDir string) ([]byte, error) {
	context := parser.NewContext()

	// Parse only once
	doc := r.md.Parser().Parse(text.NewReader(markdown), parser.WithContext(context))

	// Extract the title from the AST
	title, err := extractTitle(markdown, doc, context)
	if err != nil {
		return nil, err
	}

	// Render HTML from the same AST
	var buf bytes.Buffer
	if err := r.md.Renderer().Render(&buf, markdown, doc); err != nil {
		return nil, err
	}

	// Rewrite before applying the theme so the template's own img tags stay as they are
	html := resolveImageSources(buf.Bytes(), baseDir)

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
// Priority: 1. Front-matter title, 2. First heading, 3. "Untitled".
func extractTitle(source []byte, doc ast.Node, context parser.Context) (string, error) {
	// From front-matter
	metaData := meta.Get(context)
	if title, ok := metaData["title"].(string); ok && title != "" {
		return title, nil
	}

	// From the first heading
	heading, err := findFirstHeading(doc, source)
	if err != nil {
		return "", err
	}
	if heading != "" {
		return heading, nil
	}

	return "Untitled", nil
}

// findFirstHeading walks the AST and returns the text of the first heading.
func findFirstHeading(doc ast.Node, source []byte) (string, error) {
	var result string
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindHeading {
			text, err := extractNodeText(n, source)
			if err != nil {
				return ast.WalkStop, err
			}
			result = text
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

// extractNodeText extracts text content from a node and its children.
func extractNodeText(n ast.Node, source []byte) (string, error) {
	var buf bytes.Buffer
	err := ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && child.Kind() == ast.KindText {
			if textNode, ok := child.(*ast.Text); ok {
				buf.Write(textNode.Segment.Value(source))
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
