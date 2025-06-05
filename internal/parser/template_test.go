package parser

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTemplateBuilder(t *testing.T) {
	must := require.New(t)
	tpl := NewTemplate(
		"test",
		"Hello, {{ .target }}! I am {{ .source }} from {{ env \"PLACE\" }}",
		map[string]any{
			"target": "World",
			"source": "John",
		},
	)
	must.NotNil(tpl)

	must.NoError(os.Setenv("PLACE", "Mars"))

	buf := bytes.NewBuffer([]byte{})
	must.NoError(tpl.Build(buf))

	must.Equal("Hello, World! I am John from Mars", buf.String())
}
