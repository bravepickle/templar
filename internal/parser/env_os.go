package parser

import (
	"fmt"
	"os"
	"strings"
	//"os"
	//"strings"
)

// EnvOsParser parses OS environment params
type EnvOsParser struct{}

func (p *EnvOsParser) IsNil() bool {
	return p == nil
}

// Parse parses key-values from string and puts it to struct
func (p *EnvOsParser) Parse(_ string) (Params, error) {
	lines := os.Environ()
	result := Params{}

	for _, line := range lines {
		split := strings.SplitN(line, `=`, 2)

		if len(split) != 2 {
			return nil, fmt.Errorf(`invalid environment variable format: %s`, line)
		}

		result[split[0]] = split[1]
	}

	return result, nil
}

func NewEnvOsParser() *EnvOsParser {
	return &EnvOsParser{}
}
