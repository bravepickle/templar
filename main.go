package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bravepickle/templar/internal/command"
	"github.com/bravepickle/templar/internal/core"
)

var AppName string
var AppVersion string
var GitCommitHash string
var WorkDir string

func main() {
	//fmt.Printf("AppName: %s\n", AppName)
	//fmt.Printf("AppVersion: %s\n", AppVersion)
	//fmt.Printf("GitCommitHash: %s\n", GitCommitHash)

	if err := RunCommand(AppName, os.Args[1:], os.Stdout, AppVersion, GitCommitHash, WorkDir); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func RunCommand(name string, args []string, w io.Writer, version string, commit string, workdir string) error {
	if name == "" {
		name = core.DefaultAppName
	}

	app := core.Application{
		Version:       version,
		GitCommitHash: commit,
	}

	app.Init()

	cmd := command.NewCommand(command.NewCommandOpts{
		Name:    name,
		Args:    args,
		Output:  w,
		WorkDir: workdir,
		App:     app,
	})

	if err := cmd.Init(); err != nil {
		if cmd.Debug {
			err = fmt.Errorf("%s init: %w", cmd.Name, err)
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		if err = cmd.Usage(); err != nil {
			err = fmt.Errorf("%s usage: %w", cmd.Name, err)
			if cmd.Debug {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}

		return err
	}

	if err := cmd.Run(); err != nil {
		if cmd.Debug {
			err = fmt.Errorf("%s run: %w", cmd.Name, err)
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		if err = cmd.Usage(); err != nil {
			err = fmt.Errorf("%s usage: %w", cmd.Name, err)

			if cmd.Debug {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}

		return err
	}

	return nil
}
