package renderer

import (
	"bytes"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// resolveImageSources rewrites the src attribute of every <img> in the
// rendered HTML. A src without a scheme becomes a file:// URL under baseDir.
// Every other token is copied through unchanged.
func resolveImageSources(rendered []byte, baseDir string) []byte {
	var out bytes.Buffer
	z := html.NewTokenizer(bytes.NewReader(rendered))

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			// io.EOF is the only error a byte reader can produce
			out.Write(z.Raw())
			return out.Bytes()
		}

		if tt != html.StartTagToken && tt != html.SelfClosingTagToken {
			out.Write(z.Raw())
			continue
		}

		raw := z.Raw()
		token := z.Token()
		if token.Data != "img" {
			out.Write(raw)
			continue
		}

		changed := false
		for i, attr := range token.Attr {
			if attr.Key != "src" {
				continue
			}
			resolved, ok := resolveImageSource(attr.Val, baseDir)
			if ok {
				token.Attr[i].Val = resolved
				changed = true
			}
		}

		if !changed {
			out.Write(raw)
			continue
		}
		writeTag(&out, token, tt == html.SelfClosingTagToken)
	}
}

// resolveImageSource returns the src resolved against baseDir as a file:// URL.
// It reports false when src should stay as is: empty, unparsable, or
// already having a scheme or host.
func resolveImageSource(src, baseDir string) (string, bool) {
	// Browsers strip surrounding whitespace before resolving the URL
	src = strings.TrimSpace(src)
	if src == "" {
		return src, false
	}
	ref, err := url.Parse(src)
	if err != nil || ref.Scheme != "" || ref.Host != "" {
		return src, false
	}
	base := &url.URL{Scheme: "file", Path: strings.TrimSuffix(baseDir, "/") + "/"}
	return base.ResolveReference(ref).String(), true
}

func writeTag(out *bytes.Buffer, token html.Token, selfClosing bool) {
	out.WriteString("<")
	out.WriteString(token.Data)
	for _, attr := range token.Attr {
		out.WriteString(" ")
		out.WriteString(attr.Key)
		out.WriteString(`="`)
		out.WriteString(html.EscapeString(attr.Val))
		out.WriteString(`"`)
	}
	if selfClosing {
		out.WriteString(" />")
		return
	}
	out.WriteString(">")
}
