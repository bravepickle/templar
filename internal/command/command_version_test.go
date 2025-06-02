package command

import (
	"bytes"
	"testing"

	"github.com/bravepickle/templar/internal/core"
	"github.com/stretchr/testify/require"
)

func TestVersionCommand_Basic(t *testing.T) {
	must := require.New(t)
	cmd := &VersionCommand{}

	must.Equal("version", cmd.Name())
	must.PanicsWithError(ErrNoInit.Error(), func() {
		cmd.usage()
	})
	must.Contains(cmd.Summary(), "show application information")
	must.Error(ErrNoCommand, cmd.Init(nil, nil))
	must.False(cmd.IsNil())
	must.Error(ErrNoInit, cmd.Usage())
}

func initSubcommand(must *require.Assertions, targetCmd string, buf *bytes.Buffer) Subcommand {
	cmd := NewCommand(NewCommandOpts{
		Name:    "test-app",
		Args:    []string{targetCmd},
		Output:  buf,
		NoColor: true,
		App:     core.Application{},
	})

	must.NotNil(cmd)
	must.NoError(cmd.Init(), "init failed")

	for _, c := range cmd.commands {
		if c.Name() == targetCmd {
			return c
		}
	}

	return nil
}

func TestVersionCommand_Usage(t *testing.T) {
	must := require.New(t)
	targetCmd := SubCommandVersion
	buf := bytes.NewBuffer([]byte{})

	sub := initSubcommand(must, targetCmd, buf)

	must.NotNil(sub, "subcommand not found")
	must.NoError(sub.Usage(), "usage failed")

	output := buf.String()
	t.Log("output:", output)

	must.Contains(output, "Usage: test-app [OPTIONS] version", "text on usage output missing")
}

func TestVersionCommand_Run(t *testing.T) {
	must := require.New(t)
	targetCmd := SubCommandVersion
	buf := bytes.NewBuffer([]byte{})

	sub := initSubcommand(must, targetCmd, buf)

	must.NotNil(sub, "subcommand not found")
	must.NoError(sub.Run(), "run failed")

	output := buf.String()
	t.Log("output:", output)

	must.Contains(output, "test-app:")
	must.Contains(output, "Version:")
	must.Contains(output, "GIT commit:")
	must.Contains(output, "Working directory:")
}
