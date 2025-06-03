package parser

import (
	"github.com/bravepickle/templar/internal/core"
)

// Params parser params list
type Params map[string]any

// Parser is a common interface for data parsers
type Parser interface {
	core.Nillable

	// Parse parses input string
	Parse(in string) (Params, error)
}
