package main

import (
	"flag"
	"fmt"
)

// InputListStruct contains all options for command "list"
type InputListStruct struct {
	Help, HelpAlias bool
}

// ShowHelp Print command usage
func (t InputListStruct) ShowHelp() bool {
	return t.Help || t.HelpAlias
}

// InputListCommand options for list command stored here
var InputListCommand InputListStruct

func printListUsage() {
	fmt.Printf("Usage: %s [OPTIONS] list [COMMAND_OPTIONS] \n", InputCommon.CommandName)
	listCommand.PrintDefaults()
	fmt.Println()
	printListExamples()
}

func printListExamples() {
	cmdName := InputCommon.CommandName

	fmt.Println("Examples:")
	fmt.Printf("    %s list      List all available commands.\n", cmdName)
	fmt.Printf("    %s list -h   Show command usage help.\n", cmdName)
}

func initListCommand() {
	listCommand = flag.NewFlagSet(`list`, flag.ExitOnError)
	listCommand.BoolVar(&InputListCommand.Help, `h`, false, `Print command usage suboptions [Optional].`)
	listCommand.BoolVar(&InputListCommand.HelpAlias, `help`, false, `Print command usage suboptions [Optional].`)
}

func printCommands() {
	fmt.Println("Commands:")
	fmt.Printf("    list    List all available commands.\n")
	fmt.Printf("    init    Initialize project for templated within current dir.\n")
	fmt.Printf("    build   Build file from template.\n")
}

func doList() {
	printCommands()
}
