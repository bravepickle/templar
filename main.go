package main

import (
	"flag"
	"log"
	"os"
)

var AppVersion string
var GitCommitHash string
var AppConfigsDir string

// Parser interface is a common parser interface for all input formats
//type Parser interface {
//	// Parse reads params from source and returns them as a result
//	Parse() interface{}
//}

func init() {
	initCommands()
}

func main() {
	checkVerbosity()

	if flag.NArg() > 0 {
		switch flag.Arg(0) {
		case "list":
			index := CommandIndexArg(`list`)
			if index == -1 {
				log.Fatal(`Not found command position`)
				//os.Exit(1)
			}

			if err := listCommand.Parse(os.Args[index+1:]); err != nil {
				log.Fatal(err)
			}

			if InputListCommand.ShowHelp() {
				printListUsage()
				os.Exit(0)
			}

			doList()
			os.Exit(0)
		case "build":
			index := CommandIndexArg(`build`)
			if index == -1 {
				log.Fatal(`Not found command position`)
			}

			if err := runCommand.Parse(os.Args[index+1:]); err != nil {
				log.Fatal(err)
			}

			if InputRunCommand.ShowHelp() {
				printRunUsage()
				os.Exit(0)
			}

			if err := doBuild(); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		case "init":
			index := CommandIndexArg(`init`)
			if index == -1 {
				log.Fatal(`Not found command position`)
			}

			if err := initCommand.Parse(os.Args[index+1:]); err != nil {
				log.Fatal(err)
			}

			if InputInitCommand.ShowHelp() {
				printInitUsage()
				os.Exit(0)
			}

			doInit()
			os.Exit(0)
		default:
			printUsage()
			os.Exit(1)
		}
	}

	printUsage()
	os.Exit(0)
}
