package parser

import (
	"github.com/bravepickle/templar/v2/internal/core"
	"github.com/joho/godotenv"
)

// EnvParser parsing environment params from string
type EnvParser struct {
	// WithOsEnv defines if OS environment variables should be checked.
	// If enabled, OS env will have higher priority
	//WithOsEnv bool
}

func (p *EnvParser) IsNil() bool {
	return p == nil
}

// Parse parses key-values from string and puts it to struct
func (p *EnvParser) Parse(in string) (core.Params, error) {
	out, err := godotenv.Unmarshal(in)

	if err != nil {
		return nil, err
	}

	par := core.Params{}
	for k, v := range out {
		par[k] = v
	}

	return par, nil
}

// NewEnvParser creates parser and parses raw data string
func NewEnvParser() *EnvParser {
	return &EnvParser{}
}
