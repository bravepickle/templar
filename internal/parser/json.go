package parser

import (
	"encoding/json"
)

// JSONParser parses environment params
type JSONParser struct{}

func (p *JSONParser) IsNil() bool {
	return p == nil
}

func (p *JSONParser) Parse(in string) (Params, error) {
	var out Params
	err := json.Unmarshal([]byte(in), &out)

	return out, err
}

// NewJSONParser creates parser and parses raw data string
func NewJSONParser() *JSONParser {
	return &JSONParser{}
}
