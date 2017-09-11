package main

import (
	"bufio"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func doBuild() {
	var err error
	var inputFile, outputFile, templateFile string
	if InputRunCommand.InputFile != `` {
		inputFile, err = RealPath(InputRunCommand.InputFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !InputRunCommand.UseStdOut() {
		outputFile, err = RealPath(InputRunCommand.OutputFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !InputRunCommand.UseStdIn() {
		templateFile, err = RealPath(InputRunCommand.TemplateFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	// if InputRunCommand.InputFile[:1] != `/` { // is absolute path?
	// 	InputRunCommand.InputFile = cwd + `/` + strings.TrimLeft(InputRunCommand.InputFile, `./`)
	// }

	InputRunCommand.InputFile = inputFile
	InputRunCommand.OutputFile = outputFile
	InputRunCommand.TemplateFile = templateFile
	// outputFile := InputRunCommand.OutputFile
	// templateFile := InputRunCommand.TemplateFile
	inputFormat := InputRunCommand.InputFormat

	if verbose {
		log.Println(`Read params from:`, InputRunCommand.InputFile)
		log.Println(`Read input format:`, InputRunCommand.InputFormat)
		log.Println(`Output rendered data to:`, InputRunCommand.OutputFile)

		if InputRunCommand.UseStdIn() {
			log.Println(`Template to render: STDIN`)
		} else {
			log.Println(`Template to render:`, InputRunCommand.TemplateFile)
		}

	}

	if inputFile != `` {
		if stat, err := os.Stat(inputFile); err != nil || stat.Mode().IsRegular() == false {
			log.Fatal(`Input file is invalid for reading: `, inputFile)
		}
	}

	// TODO: fix path creation. Should support properly absolute paths

	var tplContents string

	log.Println(`TPL FILE:`, InputRunCommand.UseStdIn(), InputRunCommand.TemplateFile)
	if !InputRunCommand.UseStdIn() {
		if verbose {
			log.Println(`Attempting to read template contents from file...`)
		}
		// if templateFile[:1] != `/` { // is absolute path?
		// 	templateFile = cwd + `/` + strings.TrimLeft(templateFile, `./`)
		// }

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
					log.Println(`Found EOF. Finishing reading input...`)
				}

				break
			}

			if err != nil {
				log.Fatal(err)
			}

			log.Println(`Read string: `, line)

			tplContents += line

			log.Println(`Content: `, tplContents)
		}

		// envs := os.Environ()
		//
		// for k, v := range envs {
		// 	log.Println(`ENV:`, k, ` `, v)
		// }

		// io.Reader
		// cnt, err := io.ReadFull(os.Stdin, rawTemplate)
		// log.Println(`cnt:`, cnt)
		// rawTemplate := buf
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// tplContents = string(rawTemplate)
	}

	if tplContents == `` {
		log.Fatal(`Either template file or STDIN should contain template to render`)
	}

	var parser Parser
	var contents []byte

	if inputFile != `` {
		contents, err = ioutil.ReadFile(inputFile)
		if err != nil {
			log.Fatal(err)
		}

		if verbose && string(contents) == `` {
			log.Printf("File %s is empty", inputFile)
		}
	}

	log.Println(InputRunCommand)

	switch InputRunCommand.InputFormat {
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

	var output io.Writer
	var file *os.File

	if InputRunCommand.UseStdOut() {
		output = os.Stdout
	} else if verbose {
		err = createFile(outputFile)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		file, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		defer file.Close()

		output = io.MultiWriter(file, os.Stdout)
	} else {
		err = createFile(outputFile)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		file, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		defer file.Close()

		output = file
	}

	if err = tpl.Execute(output, parser.GetParams()); err != nil {
		log.Fatal(err, parser.GetParams())
	}
}

func createFile(path string) (err error) {
	var file *os.File
	if verbose {
		log.Println(`Attempting create file:`, path)
	}
	// detect if file exists
	_, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		if verbose {
			log.Println(`File created`)
		}
	}

	if verbose {
		log.Println(`File already exists. Rewriting...`)
		// TODO: backup file
	}

	return err
}
