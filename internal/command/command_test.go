package command

import (
	"bytes"
	"testing"

	"github.com/bravepickle/templar/v2/internal/core"
	"github.com/stretchr/testify/require"
)

func initTestSubcommand(must *require.Assertions, targetCmd string, buf *bytes.Buffer) (Subcommand, *Command) {
	cmd := NewCommand(NewCommandOpts{
		Name:    "test-app",
		Args:    []string{targetCmd},
		Output:  buf,
		NoColor: true,
		App:     core.Application{},
	})

	cmd.App.Init()
	must.NotNil(cmd)
	must.NoError(cmd.Init(), "init failed")

	for _, c := range cmd.commands {
		if c.Name() == targetCmd {
			return c, cmd
		}
	}

	return nil, cmd
}

func TestCommand_Usage(t *testing.T) {
	must := require.New(t)
	buf := bytes.NewBuffer([]byte{})

	cmd := NewCommand(NewCommandOpts{
		Name:    "test-app",
		Args:    []string{"--no-color", "help"},
		Output:  buf,
		NoColor: true,
		App:     core.Application{},
	})

	must.NotNil(cmd)
	must.NoError(cmd.Init(), "init failed")

	must.NoError(cmd.Run(), "run failed")

	output := buf.String()
	//t.Log("output:", output)

	must.Contains(output, "generate template contents with provided variables", "text on usage output missing")
	must.Contains(output, "Usage: test-app [OPTIONS] COMMAND [COMMAND_ARGS]", "text on usage output missing")

	cmd = NewCommand(NewCommandOpts{
		Name:    "test-app",
		Args:    []string{},
		Output:  buf,
		NoColor: true,
		App:     core.Application{},
	})

	must.NotNil(cmd)
	must.NoError(cmd.Init(), "init failed")

	must.NoError(cmd.Run(), "run failed")

	output = buf.String()
	//t.Log("output:", output)

	must.Contains(output, "generate template contents with provided variables", "text on usage output missing")
}
