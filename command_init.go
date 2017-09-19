package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// InputInitStruct contains all options for command "init"
type InputInitStruct struct {
	Help, HelpAlias bool
}

// ShowHelp Print command usage
func (t InputInitStruct) ShowHelp() bool {
	return t.Help || t.HelpAlias
}

// InputInitCommand options for init command stored here
var InputInitCommand InputInitStruct

func printInitUsage() {
	fmt.Printf("Usage: %s [OPTIONS] init [COMMAND_OPTIONS] \n", InputCommon.CommandName)
	initCommand.PrintDefaults()
	fmt.Println()
	printInitExamples()
}

func printInitExamples() {
	cmdName := InputCommon.CommandName

	fmt.Println("Examples:")
	fmt.Printf("    %s init      Init project under current working directory.\n", cmdName)
	fmt.Printf("    %s init -h   Show command usage help.\n", cmdName)
}

func initInitCommand() {
	initCommand = flag.NewFlagSet(`init`, flag.ExitOnError)
	initCommand.BoolVar(&InputInitCommand.Help, `h`, false, `Print command usage suboptions [Optional].`)
	initCommand.BoolVar(&InputInitCommand.HelpAlias, `help`, false, `Print command usage suboptions [Optional].`)
}

func doInit() {
	path, err := filepath.Abs(`./templates`)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		// TODO: check file exists

		err = os.Mkdir(path, 0755)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if verbose {
			log.Println(`Directory created:`, path)
		}
	} else {
		if verbose {
			log.Println(`Directory already created:`, path)
		}
	}

}
