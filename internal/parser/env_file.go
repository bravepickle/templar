package parser

import (
	"github.com/joho/godotenv"
	//"os"
	//"strings"
)

// EnvParser parsing environment params from string
type EnvParser struct{}

// Parse parses key-values from string and puts it to struct
func (p EnvParser) Parse(in string) (Params, error) {
	out, err := godotenv.Unmarshal(in)

	if err != nil {
		return nil, err
	}

	par := Params{}
	for k, v := range out {
		par[k] = v
	}

	return par, nil
}

//
//// ParseFromEnv parses key-values from OS environment
//func (p EnvParser) ParseFromEnv() map[string]string {
//	// lines := strings.Split(data, "\n")
//	//
//	lines := os.Environ()
//	result := make(map[string]string)
//
//	for _, line := range lines {
//		split := strings.SplitN(line, `=`, 2)
//
//		if len(split) != 2 {
//			//if main.verbose {
//			//	var substr string
//			//	if len(line) > MaxLineSubstringView {
//			//		substr = line[:MaxLineSubstringView] + `...`
//			//	} else {
//			//		substr = line
//			//	}
//			//	log.Printf("Failed to read key-value from line: \"%s\"\n", substr)
//			//}
//		} else {
//			result[split[0]] = split[1]
//		}
//	}
//
//	return result
//}

// NewEnvParser creates parser and parses raw data string
func NewEnvParser() EnvParser {
	return EnvParser{}
}
