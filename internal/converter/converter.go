package converter

import (
	"bytes"

	"github.com/yuin/goldmark"
)

func Convert(markdown []byte) ([]byte, error) {
	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert(markdown, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
