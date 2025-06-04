package command

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
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
		expectedStdout []string
		beforeBuild    func(sub Subcommand, cmd *Command)
		afterBuild     func(sub Subcommand, cmd *Command)
	}{
		//{
		//	name:           "no input",
		//	args:           nil,
		//	expectedErr:    "no template file specified",
		//	expectedOutput: nil,
		//	beforeBuild:    nil,
		//	afterBuild:     nil,
		//},
		{
			name:           "invalid batch path",
			args:           []string{"--batch", "unknown.txt"},
			expectedErr:    "unknown.txt: no such file or directory",
			expectedOutput: nil,
			beforeBuild:    nil,
			afterBuild:     nil,
		},
		{
			name:           "env",
			args:           []string{"--input", ".env", "--format", "env", "--template", "template.tpl", "--output", "result.txt"},
			expectedErr:    "",
			expectedOutput: nil,
			beforeBuild: func(sub Subcommand, cmd *Command) {
				cmd.WorkDir = t.TempDir()

				// Init vars file
				envFilename := filepath.Join(cmd.WorkDir, ".env")
				envFile, err := os.OpenFile(envFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
				must.NoError(err)

				defer func() {
					must.NoError(envFile.Close())
				}()

				must.NoError(os.WriteFile(envFilename, []byte("# comment\nTEST_TEMPLAR=success"), 0666))

				// Init template file
				tplFilename := filepath.Join(cmd.WorkDir, "template.tpl")
				tplFile, err := os.OpenFile(tplFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
				must.NoError(err)

				defer func() {
					must.NoError(tplFile.Close())
				}()

				must.NoError(os.WriteFile(
					tplFilename,
					[]byte("Test value: {{ .TEST_TEMPLAR }}\nDefault: {{ default \"zero\" .UNDEF_VAR }}"),
					0666,
				))
			},
			afterBuild: func(sub Subcommand, cmd *Command) {
				out, err := os.ReadFile(filepath.Join(cmd.WorkDir, "result.txt"))
				must.NoError(err)

				output := string(out)
				t.Log("rendered:", output)
				must.Contains(output, "Test value: success", "placeholder failed")
				must.Contains(output, "Default: zero", "defaults failed")
			},
		},
		{
			name:           "stdout with clear env",
			args:           []string{"--input", ".env", "--template", "template.tpl"},
			expectedErr:    "",
			expectedOutput: nil,
			beforeBuild: func(sub Subcommand, cmd *Command) {
				cmd.WorkDir = t.TempDir()
				if s, ok := sub.(*BuildCommand); ok {
					s.NoCloseWriter = true
				}

				// Init vars file
				envFilename := filepath.Join(cmd.WorkDir, ".env")
				envFile, err := os.OpenFile(envFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
				must.NoError(err)

				defer func() {
					must.NoError(envFile.Close())
				}()

				must.NoError(os.WriteFile(envFilename, []byte("TEST_TEMPLAR=foo"), 0666))

				// Init template file
				tplFilename := filepath.Join(cmd.WorkDir, "template.tpl")
				tplFile, err := os.OpenFile(tplFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
				must.NoError(err)

				defer func() {
					must.NoError(tplFile.Close())
				}()

				must.NoError(os.WriteFile(
					tplFilename,
					[]byte("Test value: {{ .TEST_TEMPLAR }}"),
					0666,
				))
			},
			//afterBuild: func(sub Subcommand, cmd *Command) {
			//	out, err := os.ReadFile(filepath.Join(cmd.WorkDir, "result.txt"))
			//	must.NoError(err)
			//
			//	output := string(out)
			//	t.Log("rendered:", output)
			//	must.Contains(output, "Test value: success", "placeholder failed")
			//	must.Contains(output, "Default: zero", "defaults failed")
			//},
			expectedStdout: []string{"Test value: foo"},
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

			var err error

			actualStdout := captureOutput(t, func() {
				err = sub.Run()
			})

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

			t.Log("stdout:", actualStdout)
			if len(d.expectedStdout) > 0 {
				for _, line := range d.expectedOutput {
					must.Contains(actualStdout, line, "missing stdout output")
				}
			}

			if d.afterBuild != nil {
				d.afterBuild(sub, cmd)
			}
		})
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	// Backup original os.Stdout
	origStdout := os.Stdout

	// Create pipe to capture output
	r, w, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = w

	// Run the function while stdout is redirected
	fn()

	// Close writer and restore original stdout
	err = w.Close()
	require.NoError(t, err)
	os.Stdout = origStdout

	// Read captured output
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	_ = r.Close()

	return buf.String()
}
