package main

import (
	"bufio"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func prepareBuildVars() {
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

	InputRunCommand.InputFile = inputFile
	InputRunCommand.OutputFile = outputFile
	InputRunCommand.TemplateFile = templateFile
	// inputFormat := InputRunCommand.InputFormat

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
}

func readContentsFromFile(filepath string) string {
	if verbose {
		log.Println(`Attempting to read template contents from file...`)
	}

	if stat, err := os.Stat(filepath); err != nil || stat.Mode().IsRegular() == false {
		log.Fatal(`File is invalid for reading: `, filepath)
	}

	rawFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	return string(rawFile)
}

func readContentsFromStdIn() string {
	var tplContents string
	if verbose {
		log.Println(`Attempting to read template contents from STDIN...`)
	}

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

		tplContents += line
	}

	return tplContents
}

func doBuild() {
	var err error
	var inputFile, outputFile, templateFile, inputFormat string
	prepareBuildVars()

	inputFile = InputRunCommand.InputFile
	outputFile = InputRunCommand.OutputFile
	templateFile = InputRunCommand.TemplateFile
	inputFormat = InputRunCommand.InputFormat

	if inputFile != `` {
		if stat, err := os.Stat(inputFile); err != nil || stat.Mode().IsRegular() == false {
			log.Fatal(`Input file is invalid for reading: `, inputFile)
		}
	}

	var tplContents string
	if !InputRunCommand.UseStdIn() {
		tplContents = readContentsFromFile(templateFile)
	} else {
		tplContents = readContentsFromStdIn()
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

	switch inputFormat {
	case `env`:
		parser = NewEnvParser(string(contents))
	default:
		log.Fatal(`Format not supported: `, inputFormat)
	}

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
