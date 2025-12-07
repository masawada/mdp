package converter

import (
	"strings"
	"testing"
)

func TestConvert_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "heading",
			input:    "# Hello",
			contains: []string{"<h1>Hello</h1>"},
		},
		{
			name:     "paragraph",
			input:    "This is a paragraph.",
			contains: []string{"<p>This is a paragraph.</p>"},
		},
		{
			name:     "unordered list",
			input:    "- item1\n- item2",
			contains: []string{"<ul>", "<li>item1</li>", "<li>item2</li>", "</ul>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert([]byte(tt.input))
			if err != nil {
				t.Fatalf("Convert() error: %v", err)
			}
			for _, want := range tt.contains {
				if !strings.Contains(string(result), want) {
					t.Errorf("Convert() result does not contain %q\ngot: %s", want, result)
				}
			}
		})
	}
}

func TestConvert_GFM(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "table",
			input:    "| a | b |\n|---|---|\n| 1 | 2 |",
			contains: []string{"<table>", "<th>a</th>", "<td>1</td>", "</table>"},
		},
		{
			name:     "strikethrough",
			input:    "~~deleted~~",
			contains: []string{"<del>deleted</del>"},
		},
		{
			name:     "task list",
			input:    "- [ ] todo\n- [x] done",
			contains: []string{"<input", "type=\"checkbox\"", "checked"},
		},
		{
			name:     "autolink",
			input:    "Visit https://example.com for more info.",
			contains: []string{"<a href=\"https://example.com\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert([]byte(tt.input))
			if err != nil {
				t.Fatalf("Convert() error: %v", err)
			}
			for _, want := range tt.contains {
				if !strings.Contains(string(result), want) {
					t.Errorf("Convert() result does not contain %q\ngot: %s", want, result)
				}
			}
		})
	}
}
