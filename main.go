package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bravepickle/templar/v2/internal/command"
	"github.com/bravepickle/templar/v2/internal/core"
)

var ErrCommandFailed = errors.New("command failed")

var AppName string
var AppVersion string
var GitCommitHash string
var WorkDir string

func main() {
	if err := RunCommand(AppName, os.Args[1:], os.Stdout, AppVersion, GitCommitHash, WorkDir); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func alertCommandFailed(cmd *command.Command, err error) {
	if cmd.Debug {
		_, _ = fmt.Fprintln(os.Stderr, err.Error()+"\n")
	} else {
		_, _ = fmt.Fprintln(os.Stderr, ErrCommandFailed.Error()+"\n")
	}
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
		alertCommandFailed(cmd, fmt.Errorf("init: %w", err))

		if err = cmd.Usage(); err != nil {
			alertCommandFailed(cmd, fmt.Errorf("usage: %w", err))

			return err
		}

		return err
	}

	if err := cmd.Run(); err != nil {
		alertCommandFailed(cmd, fmt.Errorf("run: %w", err))

		return err
	}

	return nil
}
