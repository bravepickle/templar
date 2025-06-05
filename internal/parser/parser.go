package parser

import (
	"github.com/bravepickle/templar/internal/core"
)

// Parser is a common interface for data parsers
type Parser interface {
	core.Nillable

	// Parse parses input string
	Parse(in string) (core.Params, error)
}
