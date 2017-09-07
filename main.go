package main

import (
	"bufio"
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

	if inputFile[:1] != `/` { // is absolute path?
		inputFile = cwd + `/` + strings.TrimLeft(inputFile, `./`)
	}

	if verbose {
		log.Println(`Read params from:`, inputFile)
		log.Println(`Read input format:`, inputFormat)
		log.Println(`Output rendered data to:`, outputFile)
		log.Println(`Template to render:`, templateFile)
	}

	if stat, err := os.Stat(inputFile); err != nil || stat.Mode().IsRegular() == false {
		log.Fatal(`Input file is invalid for reading: `, inputFile)
	}

	// TODO: fix path creation. Should support properly absolute paths

	var tplContents string

	if templateFile != `` { // is defined?
		if verbose {
			log.Println(`Attempting to read template contents from file...`)
		}
		if templateFile[:1] != `/` { // is absolute path?
			templateFile = cwd + `/` + strings.TrimLeft(templateFile, `./`)
		}

		if stat, err := os.Stat(templateFile); err != nil || stat.Mode().IsRegular() == false {
			log.Fatal(`Template file is invalid for reading: `, templateFile)
		}

		rawTemplate, err := ioutil.ReadFile(templateFile)
		if err != nil {
			log.Fatal(err)
		}

		tplContents = string(rawTemplate)
	} else {
		if verbose {
			log.Println(`Attempting to read template contents from STDIN...`)
		}

		// rawTemplate, err := ioutil.ReadAll(os.Stdin)
		// var rawTemplate []byte

		reader := bufio.NewReader(os.Stdin)

		for {
			line, err := reader.ReadString('\n')

			if err != nil && err.Error() == `EOF` { // handling gracefully EOF
				if verbose {
					log.Println(`Found EOF for STDIN input data. Finishing...`)
				}

				break
			}

			if err != nil {
				log.Fatal(err)
			}

			log.Println(`Read string: `, line)

			tplContents += line

			log.Println(`Content: `, tplContents)
			// TODO: ctr+d to stop input
		}

		// io.Reader
		// cnt, err := io.ReadFull(os.Stdin, rawTemplate)
		// log.Println(`cnt:`, cnt)
		// rawTemplate := buf
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// tplContents = string(rawTemplate)
	}

	contents, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	if verbose && string(contents) == `` {
		log.Printf("File %s is empty", inputFile)
	}

	if tplContents == `` {
		log.Fatal(`Either template file or STDIN should contain template to render`)
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

	if verbose {
		log.Println(`===================== Parsed Values =====================`)

		for k, v := range parser.GetParams() {
			log.Println(k, `=`, v)
		}

		log.Println(`=================== Parsed Values End ===================`)
	}

	tpl, err := template.New(outputFile).Parse(tplContents)
	if err != nil {
		log.Fatal(err)
	}

	if err = tpl.Execute(os.Stdout, parser.GetParams()); err != nil {
		log.Fatal(err)
	}
}
