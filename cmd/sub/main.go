package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bravepickle/templar/internal/command"
)

var AppVersion string
var GitCommitHash string
var WorkDir string

func main() {
	if err := RunCommand(`subCmd`, os.Args[1:], os.Stdout, AppVersion, GitCommitHash, WorkDir); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	os.Exit(0)
}

func RunCommand(name string, args []string, w io.Writer, version string, commit string, workdir string) error {
	fmt.Println("RunCommand:", name, args)

	app := command.Application{
		Version:       version,
		GitCommitHash: commit,
		WorkDir:       workdir,
	}

	app.Init()

	cmd := command.NewCommand(command.NewCommandOpts{
		Name:   name,
		Args:   args,
		Output: w,
		App:    app,
	})

	if err := cmd.Init(); err != nil {
		return fmt.Errorf("%s init: %w", cmd.Name, err)
	}

	return cmd.Run()
}
