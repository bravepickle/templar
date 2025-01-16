package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	// "html/template"
	"text/template"
)

// InputRunCommandStruct contains variables with all options for input of build command
type InputRunCommandStruct struct {
	InputFile, OutputFile, InputFormat, TemplateFile, BatchInputFile, WorkingDirectory string
	Help, HelpAlias                                                                    bool
}

// UseBatchInput returns flag if batch input was used
func (t InputRunCommandStruct) UseBatchInput() bool {
	return t.BatchInputFile != ``
}

// HasWorkDir checks if working directory is set
func (t InputRunCommandStruct) HasWorkDir() bool {
	return t.WorkingDirectory != ``
}

// UseStdIn Use STDIN for template data input
func (t InputRunCommandStruct) UseStdIn() bool {
	return t.TemplateFile == `` && !t.UseBatchInput()
}

// UseStdOut Use STDOUT for rendered template data output
func (t InputRunCommandStruct) UseStdOut() bool {
	return t.OutputFile == `` && !t.UseBatchInput()
}

// ShowHelp Print command usage
func (t InputRunCommandStruct) ShowHelp() bool {
	return t.Help || t.HelpAlias
}

func prepareBuildVars() {
	var err error
	var inputFile, outputFile, templateFile, batchInputFile string

	if InputRunCommand.InputFile != `` {
		inputFile, err = RealPath(InputRunCommand.InputFile, InputRunCommand.WorkingDirectory)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !InputRunCommand.UseStdOut() {
		outputFile, err = RealPath(InputRunCommand.OutputFile, InputRunCommand.WorkingDirectory)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !InputRunCommand.UseStdIn() {
		templateFile, err = RealPath(InputRunCommand.TemplateFile, InputRunCommand.WorkingDirectory)
		if err != nil {
			log.Fatal(err)
		}
	}

	if InputRunCommand.UseBatchInput() {
		batchInputFile, err = RealPath(InputRunCommand.BatchInputFile, InputRunCommand.WorkingDirectory)
		if err != nil {
			log.Fatal(err)
		}
	}

	InputRunCommand.InputFile = inputFile
	InputRunCommand.OutputFile = outputFile
	InputRunCommand.TemplateFile = templateFile
	InputRunCommand.BatchInputFile = batchInputFile

	if verbose {
		if InputRunCommand.UseBatchInput() {
			log.Println(`Templates to build from batch file:`, InputRunCommand.BatchInputFile)
		} else {
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
}

func readContentsFromFile(filepath string) string {
	if verbose {
		log.Printf("Attempting to read template contents from file %s...\n", filepath)
	}

	if stat, err := os.Stat(filepath); err != nil || stat.Mode().IsRegular() == false {
		log.Fatal(`File is invalid for reading: `, filepath)
	}

	rawFile, err := os.ReadFile(filepath)
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

func readTplContents(templateFile string) string {
	var tplContents string
	if !InputRunCommand.UseStdIn() {
		tplContents = readContentsFromFile(templateFile)
	} else {
		tplContents = readContentsFromStdIn()
	}

	if tplContents == `` {
		log.Fatal(`Either template file or STDIN should contain template to render`)
	}

	return tplContents
}

func assertFileReadable(filename string) {
	if stat, err := os.Stat(filename); err != nil || stat.Mode().IsRegular() == false {
		log.Fatal(`Input file is invalid for reading: `, filename)
	}
}

func dumpParsedValues(parser Parser) {
	log.Println(`===================== Parsed Values =====================`)

	params, ok := parser.GetParams().(map[string]string)
	if ok {
		for k, v := range params {
			log.Println(k, `=`, v)
		}
	} else {
		params, ok := parser.GetParams().(map[string]interface{})
		if ok {
			for k, v := range params {
				log.Println(k, `=`, v)
			}
		} else {
			log.Println(`Raw data: `, parser.GetParams())
			// fmt.Println(reflect.TypeOf(parser.GetParams()))
		}
	}

	log.Println(`=================== Parsed Values End ===================`)
}

func openOutputWriter(outputFile string) (output io.Writer, file *os.File) {
	var err error

	if InputRunCommand.UseStdOut() {
		fmt.Println("Use STDOUT ", outputFile)

		output = os.Stdout
	} else if verbose {
		fmt.Println("Creating file in verbose mode ", outputFile)
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

		output = io.MultiWriter(file, os.Stdout)
	} else {
		fmt.Println("Creating file ", outputFile)

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

		output = file
	}

	return output, file
}

func readInputFileContents(inputFile string) (result string) {
	assertFileReadable(inputFile)

	contents, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	result = string(contents)

	if verbose && result == `` {
		log.Printf("File %s is empty", inputFile)
	}

	return result
}

type batchItem struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Data   any    `json:"data"`
}

func doBuild() {
	var err error
	var inputFile, outputFile, templateFile, inputFormat string
	prepareBuildVars()

	if InputRunCommand.UseBatchInput() {
		contents := readInputFileContents(InputRunCommand.BatchInputFile)

		fmt.Println("Batch contents:", contents)

		var batchData []batchItem

		if err = json.Unmarshal([]byte(contents), &batchData); err != nil {
			log.Fatalf("Failed reading batch file: %v", err)
		}

		fmt.Printf("Batch parsed contents:\n%v\n", batchData)

		for _, item := range batchData {
			//if _, err = os.Stat(item.Input); os.IsNotExist(err) {
			//	log.Fatalf("Input file %s does not exist", item.Input)
			//}

			buildTemplate(item.Input, item.Output, item.Data)
		}

		//inputFile = InputRunCommand.InputFile
		//outputFile = InputRunCommand.OutputFile
		//templateFile = InputRunCommand.TemplateFile
		//inputFormat = InputRunCommand.InputFormat
		//
		//buildTemplate(templateFile, outputFile, inputFile, inputFormat)

		return
	}

	inputFile = InputRunCommand.InputFile
	outputFile = InputRunCommand.OutputFile
	templateFile = InputRunCommand.TemplateFile
	inputFormat = InputRunCommand.InputFormat

	var params any
	if inputFile != `` {
		var parser Parser
		contents := readInputFileContents(inputFile)

		switch inputFormat {
		case `env`:
			parser = NewEnvParser(contents)
		case `json`:
			parser = NewJSONParser(contents)
		default:
			log.Fatal(`Format not supported: `, inputFormat)
		}

		if verbose {
			dumpParsedValues(parser)
		}

		params = parser.GetParams()
	}

	buildTemplate(templateFile, outputFile, params)
}

func buildTemplate(
	templateFile string,
	outputFile string,
	params any,
) {
	var err error
	//var contents string

	tplContents := readTplContents(templateFile)

	tpl, err := template.New(outputFile).Funcs(funcMap).Parse(tplContents)
	if err != nil {
		log.Fatal(err)
	}

	output, file := openOutputWriter(outputFile)

	if file != nil {
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	fmt.Println("Building Template", tplContents)
	fmt.Printf("Output file: %s\nParams: %v\nOutput: %v\n", outputFile, params, output)

	//buf := new(bytes.Buffer)
	if err = tpl.Execute(output, params); err != nil {
		//if err = tpl.Execute(buf, params); err != nil {
		log.Fatal(err, params)
	}

	//fmt.Println("Written output: ", buf.String())
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

func printRunUsage() {
	fmt.Printf("Usage: %s [OPTIONS] build [COMMAND_OPTIONS] \n", InputCommon.CommandName)
	runCommand.PrintDefaults()
	fmt.Println()
	printRunExamples()
}

func printRunExamples() {
	cmdName := InputCommon.CommandName

	fmt.Println("Examples:")
	fmt.Printf("    %s build --format=env --input=./data.env --template=./templates/test.tpl --output=./out.txt"+
		"\n\tCreate out.txt file from test.tpl and environment parameters found in data.env file\n", cmdName)
	fmt.Printf("    echo 'Buy me {{ .ApplesCount }}.' | %s build --format=env --input=./data.env --output=./out.txt"+
		"\n\tCreate out.txt file from tamplate passed through STDIN aka piping."+
		"\n\tIf no STDIN is passed, then text can be typed directly and finished with Ctrl+D keystroke."+
		"\n\tDefault behavior when template is not specified.\n", cmdName)
	fmt.Printf("    echo 'Buy me {{ .ApplesCount }}.' | %s build --format=env --input=./data.env"+
		"\n\tOutputs rendered template from STDIN to STDOUT.\n", cmdName)
	fmt.Printf("    %s build --format=env --output=./out.txt"+
		"\n\tOutputs rendered template from STDIN to ./out.txt and sets template params from OS ENV.\n", cmdName)
}

func initRunCommand() {
	runCommand = flag.NewFlagSet(`build`, flag.ExitOnError)
	runCommand.StringVar(&InputRunCommand.InputFormat, `format`, `json`, `Input file format [Optional]. Supported values: env, json, key-value.`) // help := flag.Bool(`h`, value, usage)
	runCommand.BoolVar(&InputRunCommand.Help, `h`, false, `Print command usage sub-options [Optional].`)
	runCommand.BoolVar(&InputRunCommand.HelpAlias, `help`, false, `Print command usage sub-options [Optional].`)
	runCommand.StringVar(&InputRunCommand.InputFile, `input`, ``, `Input file to read params from [Optional].`)
	runCommand.StringVar(&InputRunCommand.OutputFile, `output`, ``, `Output file to render to [Optional].`)
	runCommand.StringVar(&InputRunCommand.TemplateFile, `template`, ``, `Template file to render [Optional].`)
	runCommand.StringVar(&InputRunCommand.WorkingDirectory, `d`, ``, `Working directory for files [Optional].`)
	runCommand.StringVar(&InputRunCommand.BatchInputFile, `batch`, ``, `File path for batch build of templates found in JSON file. May affect some other params. Format: [{"input":"","output":"","data":[]},{"input":"","output":""}] [Optional].`)
}

// InputRunCommand options for run command stored here
var InputRunCommand InputRunCommandStruct
