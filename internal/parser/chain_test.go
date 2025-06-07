package parser

import (
	"testing"

	"github.com/bravepickle/templar/internal/core"
	"github.com/stretchr/testify/require"
)

func TestChainParser(t *testing.T) {
	must := require.New(t)

	var err error
	var actual core.Params

	// case 1 - no parsers
	parser := NewChainParser()
	_, err = parser.Parse(`{}`)
	must.ErrorContains(err, "no parsers found")

	// case 2 - env parsers
	expected := map[string]any{"foo": "bar", "foo2": "bar2"}
	parser = NewChainParser(NewEnvOsParser(), NewEnvParser())

	t.Setenv("foo", "bar")
	t.Setenv("foo2", "ban")

	actual, err = parser.Parse(`
foo2=bar2
# extra
x=555
`)

	must.NoError(err)
	must.Subset(actual, expected)

	// case 3 - env parser + json
	expected = map[string]any{"foo": "bar", "foo2": "baz"}
	parser = NewChainParser(NewEnvOsParser(), NewJSONParser())

	t.Setenv("foo", "bar")
	t.Setenv("foo2", "ban")

	actual, err = parser.Parse(`{"x": "y", "foo2": "baz"}`)
	must.NoError(err)
	must.Subset(actual, expected)
	must.False(parser.IsNil())
}
