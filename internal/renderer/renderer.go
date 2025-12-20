// Package renderer converts Markdown to HTML with optional theming.
package renderer

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
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
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
		),
	)

	context := parser.NewContext()

	// パースは1回だけ
	doc := md.Parser().Parse(text.NewReader(markdown), parser.WithContext(context))

	// AST からタイトルを抽出
	title := extractTitle(markdown, doc, context)

	// 同じ AST から HTML に変換
	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, markdown, doc); err != nil {
		return nil, err
	}

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
func extractTitle(source []byte, doc ast.Node, context parser.Context) string {
	// Front-matter から取得
	metaData := meta.Get(context)
	if title, ok := metaData["title"].(string); ok && title != "" {
		return title
	}

	// 最初の heading から取得
	if heading := findFirstHeading(doc, source); heading != "" {
		return heading
	}

	return "Untitled"
}

// findFirstHeading walks the AST and returns the text of the first heading.
func findFirstHeading(doc ast.Node, source []byte) string {
	var result string
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindHeading {
			result = extractNodeText(n, source)
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	return result
}

// extractNodeText extracts text content from a node and its children.
func extractNodeText(n ast.Node, source []byte) string {
	var buf bytes.Buffer
	ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && child.Kind() == ast.KindText {
			buf.Write(child.(*ast.Text).Segment.Value(source))
		}
		return ast.WalkContinue, nil
	})
	return buf.String()
}
