package main

import (
	"flag"
	"log"
	"os"
)

// var inputFilePaths []string

// Parser interface is a common parser interface for all input formats
type Parser interface {
	GetParams() map[string]string
	ParseFromString(data string) (result map[string]string)
}

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
				os.Exit(1)
			}

			listCommand.Parse(os.Args[index+1:])

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
				os.Exit(1)
			}

			runCommand.Parse(os.Args[index+1:])

			if InputRunCommand.ShowHelp() {
				printRunUsage()
				os.Exit(0)
			}

			doBuild()
			os.Exit(0)
		case "init":
			index := CommandIndexArg(`init`)
			if index == -1 {
				log.Fatal(`Not found command position`)
				os.Exit(1)
			}

			initCommand.Parse(os.Args[index+1:])

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
