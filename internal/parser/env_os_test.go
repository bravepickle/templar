package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvOsParser(t *testing.T) {
	must := require.New(t)
	expected := map[string]any{"foo": "bar"}

	t.Setenv("foo", "bar")

	var actual map[string]any
	parser := NewEnvOsParser()
	actual, err := parser.Parse(``)

	must.NoError(err)
	must.Subset(actual, expected)
	must.False(parser.IsNil())
}
