package main

import (
	"flag"
	"fmt"
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

	if flag.NArg() > 1 {
		switch flag.Arg(0) {
		case "list":
			printCommands()
			os.Exit(0)
		case "build":
			index := CommandIndexArg(`build`)
			if index == -1 {
				log.Fatal(`Not found command position`)
				os.Exit(1)
			}

			// if verbose {
			log.Println(`Index for command build is:`, index)
			// }

			runCommand.Parse(os.Args[index+1:])

			if InputRunCommand.ShowHelp() {
				printRunUsage()
				os.Exit(0)
			}

			doBuild()
			os.Exit(0)
		default:
			fmt.Println(`HEREEEE`)
			printUsage()
			os.Exit(1)
		}
	}

	printUsage()
	os.Exit(0)
	// log.Fatal(`Missed`)

	// command = string(flag.NArg())
	// log.Fatal(command, flag.NArg())

	// cwd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
