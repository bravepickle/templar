package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var verbose bool // use as variable to easily refer to it
var runCommand, initCommand, listCommand *flag.FlagSet

// var command string // action selected when running app
// var cwd string // current working directory application was run from

func init() {
	index := strings.LastIndex(os.Args[0], `/`)
	if index != -1 {
		InputCommon.CommandName = os.Args[0][index+1:]
	}
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
	initListCommand()
	flag.Parse()

	verbose = InputCommon.IsVerbose()

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
