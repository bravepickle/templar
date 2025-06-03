package command

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildCommand_Basic(t *testing.T) {
	must := require.New(t)
	cmd := &BuildCommand{}

	must.Equal("build", cmd.Name())
	must.PanicsWithError(ErrNoInit.Error(), func() {
		cmd.usage()
	})
	must.Contains(cmd.Summary(), "render template contents with provided variables")
	must.Error(ErrNoCommand, cmd.Init(nil, nil))
	must.False(cmd.IsNil())
	must.Error(ErrNoInit, cmd.Usage())
}

func TestBuildCommand_Usage(t *testing.T) {
	must := require.New(t)
	targetCmd := SubCommandBuild
	buf := bytes.NewBuffer([]byte{})

	sub, cmd := initTestSubcommand(must, targetCmd, buf)

	must.NoError(sub.Init(cmd, []string{}))
	must.NotNil(sub, "subcommand not found")
	must.NoError(sub.Usage(), "usage failed")

	output := buf.String()
	t.Log("output:", output)

	must.Contains(output, "render template contents with provided variables", "text on usage output missing")
	must.Contains(output, "Usage: test-app [OPTIONS] build [COMMAND_OPTIONS]", "text on usage output missing")
}

func TestBuildCommand_Run(t *testing.T) {
	must := require.New(t)
	targetCmd := SubCommandBuild

	datasets := []struct {
		name           string
		args           []string
		expectedErr    string
		expectedOutput []string
		beforeBuild    func(sub Subcommand, cmd *Command)
		afterBuild     func(sub Subcommand, cmd *Command)
	}{
		{
			name:           "no input",
			args:           nil,
			expectedErr:    "no template file specified",
			expectedOutput: nil,
			beforeBuild:    nil,
			afterBuild:     nil,
		},
		{
			name:           "invalid batch path",
			args:           []string{"--batch", "unknown.txt"},
			expectedErr:    "unknown.txt: no such file or directory",
			expectedOutput: nil,
			beforeBuild:    nil,
			afterBuild:     nil,
		},
	}
	for _, d := range datasets {
		t.Run(d.name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})

			sub, cmd := initTestSubcommand(must, targetCmd, buf)
			must.NotNil(sub, "subcommand not found")
			must.NotNil(cmd, "command not found")

			cmd.WorkDir = t.TempDir() // Automatically cleaned up after test

			t.Logf("workdir: %s", cmd.WorkDir)

			must.Error(ErrNoCommand, sub.Init(nil, d.args))
			must.NoError(sub.Init(cmd, d.args))

			if d.beforeBuild != nil {
				d.beforeBuild(sub, cmd)
			}

			err := sub.Run()
			if d.expectedErr != "" {
				must.ErrorContains(err, d.expectedErr, "unexpected error")
			} else {
				must.NoError(err, "unexpected error on run")
			}

			output := buf.String()
			t.Log("output:", output)

			if len(d.expectedOutput) > 0 {
				for _, line := range d.expectedOutput {
					must.Contains(output, line, "missing output")
				}
			}

			if d.afterBuild != nil {
				d.afterBuild(sub, cmd)
			}
		})
	}
}
