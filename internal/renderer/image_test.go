package renderer

import (
	"testing"
)

func TestResolveImageSources(t *testing.T) {
	const baseDir = "/Users/you/docs"

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "resolves relative path",
			in:   `<p><img src="images/a.png" alt="a"></p>`,
			want: `<p><img src="file:///Users/you/docs/images/a.png" alt="a"></p>`,
		},
		{
			name: "resolves dot-slash path",
			in:   `<img src="./a.png" alt="">`,
			want: `<img src="file:///Users/you/docs/a.png" alt="">`,
		},
		{
			name: "resolves parent directory reference",
			in:   `<img src="../shared/a.png" alt="">`,
			want: `<img src="file:///Users/you/shared/a.png" alt="">`,
		},
		{
			name: "resolves root-relative path",
			in:   `<img src="/tmp/a.png" alt="">`,
			want: `<img src="file:///tmp/a.png" alt="">`,
		},
		{
			name: "keeps percent-encoded characters",
			in:   `<img src="my%20image.png" alt="">`,
			want: `<img src="file:///Users/you/docs/my%20image.png" alt="">`,
		},
		{
			name: "keeps query and fragment",
			in:   `<img src="a.png?v=1#top" alt="">`,
			want: `<img src="file:///Users/you/docs/a.png?v=1#top" alt="">`,
		},
		{
			name: "leaves http URL untouched",
			in:   `<img src="http://example.com/a.png" alt="">`,
			want: `<img src="http://example.com/a.png" alt="">`,
		},
		{
			name: "leaves https URL untouched",
			in:   `<img src="https://example.com/a.png" alt="">`,
			want: `<img src="https://example.com/a.png" alt="">`,
		},
		{
			name: "leaves data URL untouched",
			in:   `<img src="data:image/png;base64,AAAA" alt="">`,
			want: `<img src="data:image/png;base64,AAAA" alt="">`,
		},
		{
			name: "leaves file URL untouched",
			in:   `<img src="file:///other/a.png" alt="">`,
			want: `<img src="file:///other/a.png" alt="">`,
		},
		{
			name: "leaves protocol-relative URL untouched",
			in:   `<img src="//example.com/a.png" alt="">`,
			want: `<img src="//example.com/a.png" alt="">`,
		},
		{
			name: "leaves empty src untouched",
			in:   `<img src="" alt="">`,
			want: `<img src="" alt="">`,
		},
		{
			name: "leaves img without src untouched",
			in:   `<img alt="">`,
			want: `<img alt="">`,
		},
		{
			name: "trims surrounding whitespace like browsers do",
			in:   `<img src=" a.png " alt="">`,
			want: `<img src="file:///Users/you/docs/a.png" alt="">`,
		},
		{
			name: "leaves whitespace-only src untouched",
			in:   `<img src="  " alt="">`,
			want: `<img src="  " alt="">`,
		},
		{
			name: "leaves unparsable src untouched",
			in:   `<img src="http://[::1" alt="">`,
			want: `<img src="http://[::1" alt="">`,
		},
		{
			name: "rewrites raw html img with extra attributes",
			in:   `<img src="a.png" width="300" />`,
			want: `<img src="file:///Users/you/docs/a.png" width="300" />`,
		},
		{
			name: "rewrites src regardless of attribute order",
			in:   `<img alt="x" src="a.png">`,
			want: `<img alt="x" src="file:///Users/you/docs/a.png">`,
		},
		{
			name: "does not touch other elements",
			in:   `<a href="a.png">link</a><script src="a.js"></script>`,
			want: `<a href="a.png">link</a><script src="a.js"></script>`,
		},
		{
			name: "does not touch escaped img in code block",
			in:   `<pre><code>&lt;img src=&quot;a.png&quot;&gt;</code></pre>`,
			want: `<pre><code>&lt;img src=&quot;a.png&quot;&gt;</code></pre>`,
		},
		{
			name: "preserves surrounding content",
			in:   "<h1>Title</h1>\n<p>text <img src=\"a.png\" alt=\"\"> more</p>\n",
			want: "<h1>Title</h1>\n<p>text <img src=\"file:///Users/you/docs/a.png\" alt=\"\"> more</p>\n",
		},
		{
			name: "escapes special characters in attributes",
			in:   `<img src="a&amp;b.png" alt="&lt;x&gt;">`,
			want: `<img src="file:///Users/you/docs/a&amp;b.png" alt="&lt;x&gt;">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveImageSources([]byte(tt.in), baseDir)
			if string(got) != tt.want {
				t.Errorf("resolveImageSources() =\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}

	t.Run("encodes special characters in base directory", func(t *testing.T) {
		got := resolveImageSources([]byte(`<img src="a.png" alt="">`), "/Users/you/my docs")
		want := `<img src="file:///Users/you/my%20docs/a.png" alt="">`
		if string(got) != want {
			t.Errorf("resolveImageSources() = %s, want %s", got, want)
		}
	})
}
