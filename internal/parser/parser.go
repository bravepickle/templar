package parser

// Params parser params list
type Params map[string]any

// Parser is a common interface for data parsers
type Parser interface {
	// Parse parses input string
	Parse(in string) (Params, error)
}
