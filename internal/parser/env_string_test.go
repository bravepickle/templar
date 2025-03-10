package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvParser(t *testing.T) {
	must := require.New(t)
	expected := map[string]any{"foo": "bar"}

	in := `
# comment
foo = bar
`
	var actual map[string]any

	parser := NewEnvParser()
	actual, err := parser.Parse(in)
	//t.Logf("Parsed env: %s", actual)

	must.NoError(err)
	must.Equal(expected, actual)
}
