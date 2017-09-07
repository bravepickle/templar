package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

var inputFilePaths []string

// Parser interface is a common parser interface for all input formats
type Parser interface {
	GetParams() map[string]string
	ParseFromString(data string) (result map[string]string)
}

var verbose bool

func main() {
	var inputFile, outputFile, inputFormat, templateFile string

	flag.StringVar(&inputFile, `input`, ``, `Input file to read params from [Optional].`)
	flag.StringVar(&inputFormat, `format`, `json`, `Input file format [Optional]. Supported values: env, json, key-value.`) // help := flag.Bool(`h`, value, usage)
	flag.StringVar(&outputFile, `output`, ``, `Output file to render to [Optional].`)
	flag.StringVar(&templateFile, `template`, ``, `Template file to render [Optional].`)
	flag.BoolVar(&verbose, `verbose`, false, `Run in verbose mode [Optional].`)
	flag.Parse()

	if verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.SetPrefix(`DEBUG: `)
		log.Println(`Running in verbose mode.`)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	inputFile = cwd + `/` + strings.TrimLeft(inputFile, `./`)

	if verbose {
		log.Println(`Read params from:`, inputFile)
		log.Println(`Read input format:`, inputFormat)
		log.Println(`Output rendered data to:`, outputFile)
		log.Println(`Template to render:`, templateFile)
	}

	if stat, err := os.Stat(inputFile); err != nil || stat.Mode().IsRegular() == false {
		log.Fatal(`Input file is invalid for reading: `, inputFile)
	}

	contents, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	if verbose && string(contents) == `` {
		log.Printf("File %s is empty", inputFile)
	}

	var parser Parser

	switch inputFormat {
	case `env`:
		parser = NewEnvParser(string(contents))
	default:
		log.Fatal(`Format not supported: `, inputFormat)

	}

	// fmt.Println(`Hello, world`, inputFile, inputFormat, outputFile, cwd, string(contents))
	// fmt.Println(parser)

	// log.Println(`=====================`)
	//
	// for k, v := range parser.GetParams() {
	// 	log.Println(`Parsed`, k, `=`, v)
	// }
	//
	// log.Println(`=====================`)

	tpl, err := template.New(outputFile).Parse(`This is a tryoute for "{{ .TEST }}" value`)
	if err != nil {
		log.Fatal(err)
	}

	if err = tpl.Execute(os.Stdout, parser.GetParams()); err != nil {
		log.Fatal(err)
	}
}
