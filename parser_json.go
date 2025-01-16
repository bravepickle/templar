package main

import (
	"encoding/json"
	"log"
)

// JSONParser Parsing environment params
type JSONParser struct {
	Params any
}

// GetParams read all params from parser
func (p JSONParser) GetParams() interface{} {
	return p.Params
}

// ParseFromString parses key-values from string and puts it to struct
func (p JSONParser) ParseFromString(data string) interface{} {
	var result interface{}

	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Fatal(err)
	}

	// log.Println(`Data from JSON:`, result)

	// result := make(map[string]string)
	// for _, line := range lines {
	// 	line = strings.TrimSpace(line)
	// 	if line == `` || line[0:1] == `#` {
	// 		if verbose {
	// 			log.Printf(`Skipping line: %s`, line)
	// 		}
	//
	// 		continue
	// 	}
	//
	// 	split := strings.SplitN(line, `=`, 2)
	//
	// 	if len(split) != 2 {
	// 		if verbose {
	// 			var substr string
	// 			if len(line) > MaxLineSubstringView {
	// 				substr = line[:MaxLineSubstringView] + `...`
	// 			} else {
	// 				substr = line
	// 			}
	// 			log.Printf("Failed to read key-value from line: \"%s\"\n", substr)
	// 		}
	// 	} else {
	// 		result[strings.TrimSpace(split[0])] = strings.TrimSpace(split[1])
	// 	}
	// }

	return result
}

// NewJSONParser creates parser and parses raw data string
func NewJSONParser(rawData string) JSONParser {
	var parser JSONParser
	parser.Params = parser.ParseFromString(rawData)

	return parser
}
