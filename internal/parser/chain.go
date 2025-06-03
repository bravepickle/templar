package parser

import (
	"errors"
	"fmt"
)

// ChainParser parsing environment params from parsers chain.
// Last parser will override all previously defined values with the same name.
type ChainParser struct {
	parsers []Parser
}

func (p *ChainParser) IsNil() bool {
	return p == nil
}

// Parse parses key-values from string and puts it to struct
func (p *ChainParser) Parse(in string) (Params, error) {
	if len(p.parsers) == 0 {
		return nil, errors.New("no parsers found")
	}

	par := Params{}
	var subPar Params
	var err error

	for _, parser := range p.parsers {
		subPar, err = parser.Parse(in)
		if err != nil {
			return nil, fmt.Errorf("failed to apply parser %T: %w", parser, err)
		}

		for k, v := range subPar {
			par[k] = v
		}
	}

	return par, nil
}

func NewChainParser(parser ...Parser) *ChainParser {
	return &ChainParser{parsers: parser}
}
