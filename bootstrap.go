package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// InputRunCommandStruct contains variables with all options for input of build command
type InputRunCommandStruct struct {
	InputFile, OutputFile, InputFormat, TemplateFile string
	Help, HelpAlias                                  bool
}

// UseStdIn Use STDIN for template data input
func (t InputRunCommandStruct) UseStdIn() bool {
	return t.TemplateFile == ``
}

// UseStdOut Use STDOUT for rendered template data output
func (t InputRunCommandStruct) UseStdOut() bool {
	return t.OutputFile == ``
}

// ShowHelp Print command usage
func (t InputRunCommandStruct) ShowHelp() bool {
	return t.Help || t.HelpAlias
}

// InputCommonStruct contains all basic options for running application
type InputCommonStruct struct {
	CommandName                                      string
	InputFile, OutputFile, InputFormat, TemplateFile string
	Help, HelpAlias, Verbose, VerboseAlias           bool
}

// ShowHelp Print command usage
func (t InputCommonStruct) ShowHelp() bool {
	return t.Help || t.HelpAlias
}

// IsVerbose Run in verbose mode
func (t InputCommonStruct) IsVerbose() bool {
	// return true
	return t.Verbose || t.VerboseAlias
}

// RealPath creates real path to file or folder from relevant
func RealPath(path string) (string, error) {
	return filepath.Abs(path)
}

// InputRunCommand options for run command stored here
var InputRunCommand InputRunCommandStruct

// InputCommon basic options for running application
var InputCommon InputCommonStruct

var verbose bool // use as variable to easily refer to it
var runCommand, initCommand *flag.FlagSet

// var command string // action selected when running app
// var cwd string // current working directory application was run from

func init() {
	index := strings.LastIndex(os.Args[0], `/`)
	if index != -1 {
		InputCommon.CommandName = os.Args[0][index+1:]
	}
}

func printUsage() {
	cmdName := InputCommon.CommandName
	fmt.Printf("Usage: %s [OPTIONS] [COMMAND] [COMMAND_OPTIONS] \n", cmdName)
	flag.PrintDefaults()
	fmt.Println()
	printCommands()
	fmt.Println()
	printExamples()
}

func printExamples() {
	cmdName := InputCommon.CommandName

	fmt.Println("Examples:")
	fmt.Printf("    %s -h      See help for using this command\n", cmdName)
	fmt.Printf("    %s init    Initialize current working directory as new project\n", cmdName)
	fmt.Printf("    %s --verbose\n\tinit Initialize new project in verbose mode\n", cmdName)
	fmt.Printf("    %s build --format=env --input=./data.env --template=./templates/test.tpl --output=./out.txt"+
		"\n\t Create out.txt file from test.tpl and environment parameters found in data.env file\n", cmdName)
}

func printCommands() {
	fmt.Println("Commands:")
	fmt.Printf("    list    List all available commands.\n")
	fmt.Printf("    init    Initialize project for templated within current dir.\n")
	fmt.Printf("    build   Build file from template.\n")
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
}

func initRunCommand() {
	runCommand = flag.NewFlagSet(`build`, flag.ExitOnError)
	runCommand.StringVar(&InputRunCommand.InputFormat, `format`, `json`, `Input file format [Optional]. Supported values: env, json, key-value.`) // help := flag.Bool(`h`, value, usage)
	runCommand.BoolVar(&InputRunCommand.Help, `h`, false, `Print command usage suboptions [Optional].`)
	runCommand.BoolVar(&InputRunCommand.HelpAlias, `help`, false, `Print command usage suboptions [Optional].`)
	runCommand.StringVar(&InputRunCommand.InputFile, `input`, ``, `Input file to read params from [Optional].`)
	runCommand.StringVar(&InputRunCommand.OutputFile, `output`, ``, `Output file to render to [Optional].`)
	runCommand.StringVar(&InputRunCommand.TemplateFile, `template`, ``, `Template file to render [Optional].`)
}

func checkVerbosity() {
	if verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.SetPrefix(`DEBUG: `)
		log.Println(`Running in verbose mode.`)
	} else {
		log.SetFlags(log.LstdFlags)
	}
}

func initCommands() {
	flag.BoolVar(&InputCommon.Verbose, `verbose`, false, `Run in verbose mode [Optional].`)
	flag.BoolVar(&InputCommon.VerboseAlias, `v`, false, `Run in verbose mode [Optional].`)
	flag.BoolVar(&InputCommon.Help, `h`, false, `Print command usage options [Optional].`)
	flag.BoolVar(&InputCommon.HelpAlias, `help`, false, `Print command usage options [Optional].`)
	initRunCommand()
	initInitCommand()
	flag.Parse()

	verbose = InputCommon.IsVerbose()
	// verbose = true

	if InputCommon.Help || InputCommon.HelpAlias {
		printUsage()
		os.Exit(0)
	}
}

// CommandIndexArg get index key for specified argument
// return -1 if not found
func CommandIndexArg(argument string) int {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == argument {
			return i
		}
	}

	return -1
}
