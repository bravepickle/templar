package command

import (
	"bytes"
	"testing"

	"github.com/bravepickle/templar/v2/internal/core"
	"github.com/stretchr/testify/require"
)

func TestHelpCommand_Basic(t *testing.T) {
	must := require.New(t)
	cmd := NewCommand(NewCommandOpts{
		Name:    core.DefaultAppName,
		NoColor: true,
	})

	must.NoError(cmd.Init())
	sc := cmd.commands[SubCommandHelp]
	sub, ok := sc.(*HelpCommand)
	must.True(ok, "sub command must implement HelpCommand")

	must.NotEmpty(sub)
	must.Equal("help", sub.Name())

	sub.fs = nil
	must.PanicsWithError(ErrNoInit.Error(), func() {
		sub.usage()
	})

	must.Contains(sub.Summary(), "show help information on command or subcommand usage")
	must.Error(ErrNoInit, sub.Usage())
	must.Error(ErrNoInit, sub.Run())
	must.False(sub.IsNil())
}

func TestHelpCommand_Usage(t *testing.T) {
	must := require.New(t)
	targetCmd := SubCommandHelp
	buf := bytes.NewBuffer([]byte{})

	sub, cmd := initTestSubcommand(must, targetCmd, buf)

	must.NoError(sub.Init(cmd, []string{}))
	must.NotNil(sub, "subcommand not found")
	must.NoError(sub.Usage(), "usage failed")

	output := buf.String()
	t.Log("output:", output)

	must.Contains(output, "Usage: test-app [OPTIONS] help [COMMAND]", "text on usage output missing")
}

func TestHelpCommand_Run(t *testing.T) {
	must := require.New(t)
	targetCmd := SubCommandHelp
	buf := bytes.NewBuffer([]byte{})

	sub, cmd := initTestSubcommand(must, targetCmd, buf)
	must.NotNil(sub, "subcommand not found")
	must.NotNil(cmd, "command not found")

	must.Error(ErrNoCommand, sub.Init(nil, []string{}))
	must.NoError(sub.Init(cmd, []string{}))
	must.NoError(sub.Run(), "run failed")

	output := buf.String()
	//t.Log("output:", output)

	must.Contains(output, "generate template contents with provided variables")

	// help command info
	buf.Reset()
	must.NoError(sub.Init(cmd, []string{"help"}))
	must.NoError(sub.Run(), "run failed")

	output = buf.String()
	//t.Log("output:", output)

	must.Contains(output, "show help information on command or subcommand usage")

	// help init info
	buf.Reset()
	must.NoError(sub.Init(cmd, []string{"init"}))
	must.NoError(sub.Run(), "run failed")

	output = buf.String()
	//t.Log("output:", output)

	must.Contains(output, "init default files structure for building templates")

	// help unknown info
	buf.Reset()
	must.NoError(sub.Init(cmd, []string{"unknown"}))
	must.NoError(sub.Run(), "run failed")

	output = buf.String()
	//t.Log("output:", output)

	must.Contains(output, "generate template contents")
}
