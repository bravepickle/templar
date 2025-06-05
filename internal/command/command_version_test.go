package command

import (
	"bytes"
	"testing"

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

func TestVersionCommand_Usage(t *testing.T) {
	must := require.New(t)
	targetCmd := SubCommandVersion
	buf := bytes.NewBuffer([]byte{})

	sub, cmd := initTestSubcommand(must, targetCmd, buf)

	must.NoError(sub.Init(cmd, []string{}))
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

	sub, cmd := initTestSubcommand(must, targetCmd, buf)
	must.NotNil(sub, "subcommand not found")
	must.NotNil(cmd, "command not found")

	must.Error(ErrNoCommand, sub.Init(nil, []string{}))
	must.NoError(sub.Init(cmd, []string{}))
	must.NoError(sub.Run(), "run failed")

	output := buf.String()
	t.Log("output:", output)

	must.Contains(output, "test-app:")
	must.Contains(output, "Version:")
	must.Contains(output, "GIT commit:")
	must.Contains(output, "Working directory:")
}
