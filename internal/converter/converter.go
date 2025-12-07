package converter

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func Convert(markdown []byte) ([]byte, error) {
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
	if err := md.Convert(markdown, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
