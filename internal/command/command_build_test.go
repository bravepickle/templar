package command

import (
	"bytes"
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

	must.NotNil(sub)
	must.NotNil(cmd)
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
			expectedErr:    "no template contents provided",
			expectedOutput: nil,
			beforeBuild:    nil,
			afterBuild:     nil,
		},
		{
			name:           "invalid batch path",
			args:           []string{"--input", "unknown.txt", "--format", "batch"},
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

				must.NoError(os.WriteFile(envFilename, []byte("# comment\nTEST_TEMPLAR=success"), 0666))

				// Init template file
				tplFilename := filepath.Join(cmd.WorkDir, "template.tpl")

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
			expectedOutput: []string{"Test value: foo"},
			beforeBuild: func(sub Subcommand, cmd *Command) {
				cmd.WorkDir = t.TempDir()
				if s, ok := sub.(*BuildCommand); ok {
					s.NoCloseWriter = true
				}

				// Init vars file
				envFilename := filepath.Join(cmd.WorkDir, ".env")

				must.NoError(os.WriteFile(envFilename, []byte("TEST_TEMPLAR=foo"), 0666))

				// Init template file
				tplFilename := filepath.Join(cmd.WorkDir, "template.tpl")

				must.NoError(os.WriteFile(
					tplFilename,
					[]byte("Test value: {{ .TEST_TEMPLAR }}"),
					0666,
				))
			},
		},
		{
			name:           "json",
			args:           []string{"--input", "vars.json", "--format", "json", "--template", "file.tpl", "--output", "result.txt"},
			expectedErr:    "",
			expectedOutput: nil,
			beforeBuild: func(sub Subcommand, cmd *Command) {
				var ok bool
				var buildCmd *BuildCommand

				if buildCmd, ok = sub.(*BuildCommand); !ok {
					t.Fatal("sub is not a BuildCommand")
				}

				cmd.WorkDir = t.TempDir()

				// Init vars file
				varsFilepath := filepath.Join(cmd.WorkDir, buildCmd.InputFile)

				must.NoError(os.WriteFile(varsFilepath, []byte(`{"magic":{"status": "real"}}`), 0666))

				// Init template file
				tplFilename := filepath.Join(cmd.WorkDir, buildCmd.TemplateFile)

				must.NoError(os.WriteFile(
					tplFilename,
					[]byte("I am sure the magic is {{ default \"boring\" .magic.status }}!\nPS: {{ env \"TEST_QUOTE\" }}"),
					0666,
				))

				t.Setenv("TEST_QUOTE", "you are wonder!")
			},
			afterBuild: func(sub Subcommand, cmd *Command) {
				out, err := os.ReadFile(filepath.Join(cmd.WorkDir, "result.txt"))
				must.NoError(err)

				output := string(out)
				t.Log("rendered:", output)
				must.Contains(output, "I am sure the magic is real!", "placeholder failed")
				must.Contains(output, "PS: you are wonder!", "placeholder failed")
			},
		},
		{
			name:           "batch",
			args:           []string{"--input", "batch.json", "--format", "batch"},
			expectedErr:    "",
			expectedOutput: nil,
			beforeBuild: func(sub Subcommand, cmd *Command) {
				var ok bool
				var buildCmd *BuildCommand

				if buildCmd, ok = sub.(*BuildCommand); !ok {
					t.Fatal("sub is not a BuildCommand")
				}

				cmd.WorkDir = t.TempDir()

				t.Setenv("TEST_QUOTE", "ENV VAR!")

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, buildCmd.InputFile),
					[]byte(`{
  "items": [
    {
      "info": "This template will combine defaults with other values.",
      "target": "rendered.txt",
      "template": "custom.tpl",
      "variables": {
        "foo": "custom",
        "extra": "extra value"
      }
    },
    {
      "target": "with_defaults.txt"
    }
  ],
  "defaults": {
    "info": "This section defines default values to be used in the \"items\" section.",
    "template": "default.tpl",
    "variables": {
      "foo": "bar",
      "size": 42,
      "nested": {"baz": "faz"}
    }
  }
}
`), 0666))

				t.Setenv("TEST_QUOTE", "this is env variable")

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, "default.tpl"),
					[]byte(`Default template with values:
foo = {{ default "UNDEFINED" .foo }}
size = {{ default "UNDEFINED" .size }}
nested = {{ toJson .nested }}
extra = {{ default "UNDEFINED" .extra }}
ENV TEST_QUOTE = {{ env "TEST_QUOTE" }}
`),
					0666,
				))

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, "custom.tpl"),
					[]byte(`Custom template with values:
foo = {{ default "UNDEFINED" .foo }}
size = {{ default "UNDEFINED" .size }}
nested = {{ toJson .nested }}
extra = {{ default "UNDEFINED" .extra }}
ENV TEST_QUOTE = {{ env "TEST_QUOTE" }}
`),
					0666,
				))

			},
			afterBuild: func(sub Subcommand, cmd *Command) {
				out, err := os.ReadFile(filepath.Join(cmd.WorkDir, "rendered.txt"))
				must.NoError(err)

				output := string(out)
				t.Log("rendered.txt:", output)
				must.Equal(`Custom template with values:
foo = custom
size = UNDEFINED
nested = null
extra = extra value
ENV TEST_QUOTE = this is env variable
`, output, "rendered.txt file is invalid")

				out, err = os.ReadFile(filepath.Join(cmd.WorkDir, "with_defaults.txt"))
				must.NoError(err)

				output = string(out)
				t.Log("with_defaults.txt:", output)
				must.Equal(`Default template with values:
foo = bar
size = 42
nested = {"baz":"faz"}
extra = UNDEFINED
ENV TEST_QUOTE = this is env variable
`, output, "with_defaults.txt file is invalid")
			},
		},
		{
			name:           "batch with input files",
			args:           []string{"--input", "batch.json", "--format", "batch"},
			expectedErr:    "",
			expectedOutput: nil,
			beforeBuild: func(sub Subcommand, cmd *Command) {
				var ok bool
				var buildCmd *BuildCommand

				if buildCmd, ok = sub.(*BuildCommand); !ok {
					t.Fatal("sub is not a BuildCommand")
				}

				cmd.WorkDir = t.TempDir()

				t.Setenv("TEST_QUOTE", "SRC_OS_ENV")

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, buildCmd.InputFile),
					[]byte(`{
  "items": [
    {
      "target": "rendered.txt",
      "template": "custom.tpl",
      "format": "env",
      "input": "custom.env"
    },
    {
      "target": "with_defaults.txt"
    }
  ],
  "defaults": {
    "template": "default.tpl",
    "format": "json",
    "input": "default.json"
  }
}
`), 0666))

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, "default.json"),
					[]byte(`{"test_var": "SRC_JSON"}`),
					0666,
				))

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, "default.tpl"),
					[]byte(`Default template:
test_var = {{ .test_var }}
TEST_QUOTE = {{ env "TEST_QUOTE" }}
`),
					0666,
				))

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, "custom.env"),
					[]byte("test_var=SRC_ENV"),
					0666,
				))

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, "custom.tpl"),
					[]byte(`Custom template:
test_var = {{ .test_var }}
TEST_QUOTE = {{ env "TEST_QUOTE" }}
`),
					0666,
				))

			},
			afterBuild: func(sub Subcommand, cmd *Command) {
				out, err := os.ReadFile(filepath.Join(cmd.WorkDir, "rendered.txt"))
				must.NoError(err)

				output := string(out)
				t.Log("rendered.txt:", output)
				must.Equal(`Custom template:
test_var = SRC_ENV
TEST_QUOTE = SRC_OS_ENV
`, output, "rendered.txt file is invalid")

				out, err = os.ReadFile(filepath.Join(cmd.WorkDir, "with_defaults.txt"))
				must.NoError(err)

				output = string(out)
				t.Log("with_defaults.txt:", output)
				must.Equal(`Default template:
test_var = SRC_JSON
TEST_QUOTE = SRC_OS_ENV
`, output, "with_defaults.txt file is invalid")
			},
		},
		{
			name:        "debug dump env",
			args:        []string{"--input", "vars.json", "--format", "json", "--clear", "--dump", "env"},
			expectedErr: "",
			expectedOutput: []string{
				"foo=\"myJSON\"",
				`user={"age":42,"name":"John"}`,
			},
			beforeBuild: func(sub Subcommand, cmd *Command) {
				cmd.Debug = true

				var ok bool
				var buildCmd *BuildCommand

				if buildCmd, ok = sub.(*BuildCommand); !ok {
					t.Fatal("sub is not a BuildCommand")
				}

				cmd.WorkDir = t.TempDir()

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, buildCmd.InputFile),
					[]byte(`{"foo": "myJSON", "user":{"name":"John", "age": 42}}`), 0666))

				t.Setenv("TEST_QUOTE", "this is env variable")
			},
		},
		{
			name:        "debug dump json",
			args:        []string{"--input", "vars.json", "--format", "json", "--clear", "--dump", "json"},
			expectedErr: "",
			expectedOutput: []string{
				`"foo": "myJSON"`,
				`"age": 42`,
			},
			beforeBuild: func(sub Subcommand, cmd *Command) {
				cmd.Debug = true

				var ok bool
				var buildCmd *BuildCommand

				if buildCmd, ok = sub.(*BuildCommand); !ok {
					t.Fatal("sub is not a BuildCommand")
				}

				cmd.WorkDir = t.TempDir()

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, buildCmd.InputFile),
					[]byte(`{"foo": "myJSON", "user":{"name":"John", "age": 42}}`), 0666))
			},
		},
		{
			name:           "debug dump json",
			args:           []string{"--input", "vars.json", "--format", "json", "--clear", "--dump", "json_compact"},
			expectedErr:    "",
			expectedOutput: []string{`{"foo":"myJSON","user":{"age":42,"name":"John"}}`},
			beforeBuild: func(sub Subcommand, cmd *Command) {
				cmd.Debug = true

				var ok bool
				var buildCmd *BuildCommand

				if buildCmd, ok = sub.(*BuildCommand); !ok {
					t.Fatal("sub is not a BuildCommand")
				}

				cmd.WorkDir = t.TempDir()

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, buildCmd.InputFile),
					[]byte(`{"foo": "myJSON", "user":{"name":"John", "age": 42}}`), 0666))
			},
		},
		{
			name:           "verbose dump env",
			args:           []string{"--input", "vars.json", "--format", "json", "--clear", "--dump", "env"},
			expectedErr:    "",
			expectedOutput: []string{`foo=string`, `user=map[string]interface {}`},
			beforeBuild: func(sub Subcommand, cmd *Command) {
				cmd.Verbose = true

				var ok bool
				var buildCmd *BuildCommand

				if buildCmd, ok = sub.(*BuildCommand); !ok {
					t.Fatal("sub is not a BuildCommand")
				}

				cmd.WorkDir = t.TempDir()

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, buildCmd.InputFile),
					[]byte(`{"foo": "myJSON", "user":{"name":"John", "age": 42}}`), 0666))
			},
		},
		{
			name:           "basic dump env",
			args:           []string{"--input", "vars.json", "--format", "json", "--clear", "--dump", "env"},
			expectedErr:    "",
			expectedOutput: []string{"foo\nuser\n"},
			beforeBuild: func(sub Subcommand, cmd *Command) {
				var ok bool
				var buildCmd *BuildCommand

				if buildCmd, ok = sub.(*BuildCommand); !ok {
					t.Fatal("sub is not a BuildCommand")
				}

				cmd.WorkDir = t.TempDir()

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, buildCmd.InputFile),
					[]byte(`{"foo": "myJSON", "user":{"name":"John", "age": 42}}`), 0666))
			},
		},
		{
			name:           "empty verbose dump env",
			args:           []string{"--input", "vars.json", "--format", "env", "--clear", "--dump", "env"},
			expectedErr:    "",
			expectedOutput: []string{"No variables found"},
			beforeBuild: func(sub Subcommand, cmd *Command) {
				cmd.Verbose = true

				var ok bool
				var buildCmd *BuildCommand

				if buildCmd, ok = sub.(*BuildCommand); !ok {
					t.Fatal("sub is not a BuildCommand")
				}

				cmd.WorkDir = t.TempDir()

				must.NoError(os.WriteFile(
					filepath.Join(cmd.WorkDir, buildCmd.InputFile),
					[]byte(``), 0666))
			},
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

func TestBuildCommand_ErrorsHandling(t *testing.T) {
	must := require.New(t)
	//targetCmd := SubCommandBuild
	buf := bytes.NewBuffer([]byte{})

	sc, cmd := initTestSubcommand(must, SubCommandBuild, buf)
	sub, ok := sc.(*BuildCommand)
	must.True(ok)

	must.NoError(sub.Init(cmd, nil))
	must.Error(ErrNoInit, sub.Run())

	// read from stdin empty string
	in, err := sub.readInput("")
	must.NoError(err)
	must.Equal([]byte{}, in, "stdin is not empty")

	r, w, err := os.Pipe()
	must.NoError(err)
	closed := false
	defer func() {
		if !closed {
			must.NoError(w.Close())
		}
	}()
	_, err = w.Write([]byte("test me"))
	must.NoError(err)

	closed = true
	must.NoError(w.Close())

	// read from the input file
	sub.In = r
	in, err = sub.readInput("")
	must.NoError(err)
	must.Equal("test me", string(in), "input reader mismatch")
}
