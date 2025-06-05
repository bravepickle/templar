package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONParser(t *testing.T) {
	must := require.New(t)

	in := `{"faz": "baz"}`
	expected := map[string]any{"faz": "baz"}
	var actual map[string]any

	parser := NewJSONParser()
	actual, err := parser.Parse(in)
	must.NoError(err)
	must.Equal(expected, actual)
	must.False(parser.IsNil())
}
