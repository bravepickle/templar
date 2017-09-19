package main

import (
	"log"
	"strings"
)

// MaxLineSubstringView specifies max lenght of line to show in logs
const MaxLineSubstringView = 20

// EnvParser Parsing environment params
type EnvParser struct {
	Params map[string]string
}

// GetParams read all params from parser
func (p EnvParser) GetParams() interface{} {
	return p.Params
}

// ParseFromString parses key-values from string and puts it to struct
func (p EnvParser) ParseFromString(data string) map[string]string {
	lines := strings.Split(data, "\n")

	result := make(map[string]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == `` || line[0:1] == `#` {
			if verbose {
				log.Printf(`Skipping line: %s`, line)
			}

			continue
		}

		split := strings.SplitN(line, `=`, 2)

		if len(split) != 2 {
			if verbose {
				var substr string
				if len(line) > MaxLineSubstringView {
					substr = line[:MaxLineSubstringView] + `...`
				} else {
					substr = line
				}
				log.Printf("Failed to read key-value from line: \"%s\"\n", substr)
			}
		} else {
			result[strings.TrimSpace(split[0])] = strings.TrimSpace(split[1])
		}
	}

	return result
}

// NewEnvParser creates parser and parses raw data string
func NewEnvParser(rawData string) EnvParser {
	var parser EnvParser
	parser.Params = parser.ParseFromString(rawData)

	return parser
}
