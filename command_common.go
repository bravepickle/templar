package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

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

// InputCommon basic options for running application
var InputCommon InputCommonStruct
