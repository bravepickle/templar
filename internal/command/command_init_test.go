package command

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitCommand_Basic(t *testing.T) {
	must := require.New(t)
	cmd := &InitCommand{}

	must.Equal("init", cmd.Name())
	must.PanicsWithError(ErrNoInit.Error(), func() {
		cmd.usage()
	})
	must.Contains(cmd.Summary(), "init default files structure for building templates")
	must.Error(ErrNoCommand, cmd.Init(nil, nil))
	must.False(cmd.IsNil())
	must.Error(ErrNoInit, cmd.Usage())
}

func TestInitCommand_Usage(t *testing.T) {
	must := require.New(t)
	targetCmd := SubCommandInit
	buf := bytes.NewBuffer([]byte{})

	sub, cmd := initTestSubcommand(must, targetCmd, buf)

	must.NoError(sub.Init(cmd, []string{}))
	must.NotNil(sub, "subcommand not found")
	must.NoError(sub.Usage(), "usage failed")

	output := buf.String()
	t.Log("output:", output)

	must.Contains(output, "init default files structure", "text on usage output missing")
	must.Contains(output, "test-app [OPTIONS] init [COMMAND_OPTIONS]", "text on usage output missing")
}

func TestInitCommand_Run(t *testing.T) {
	must := require.New(t)
	targetCmd := SubCommandInit
	buf := bytes.NewBuffer([]byte{})

	sub, cmd := initTestSubcommand(must, targetCmd, buf)
	must.NotNil(sub, "subcommand not found")
	must.NotNil(cmd, "command not found")

	cmd.WorkDir = t.TempDir() // Automatically cleaned up after test

	t.Logf("workdir: %s", cmd.WorkDir)

	must.Error(ErrNoCommand, sub.Init(nil, []string{}))
	must.NoError(sub.Init(cmd, []string{}))

	must.NoError(sub.Run(), "run failed")

	output := buf.String()
	t.Log("output:", output)

	must.Contains(output, "Templates created: "+cmd.WorkDir)

	var expectedFilePath string
	var actual []byte
	var err error
	files := map[string]string{
		"variables.json":        ExampleJson,
		"variables.env":         ExampleEnv,
		"batch.json":            ExampleBatchJson,
		"templates/example.tpl": ExampleTemplate,
		"templates/custom.tpl":  ExampleCustomTemplate,
	}

	for filename, expectedContents := range files {
		expectedFilePath = filepath.Join(cmd.WorkDir, filename)
		must.FileExists(expectedFilePath, "%s not present", filename)
		actual, err = os.ReadFile(expectedFilePath)
		must.NoError(err)
		must.Equal(expectedContents, string(actual), "%s file contents mismatch", filename)
	}
}

func TestInitCommand_Run_NoBatch(t *testing.T) {
	must := require.New(t)
	targetCmd := SubCommandInit
	buf := bytes.NewBuffer([]byte{})

	sub, cmd := initTestSubcommand(must, targetCmd, buf)
	must.NotNil(sub, "subcommand not found")
	must.NotNil(cmd, "command not found")

	cmd.WorkDir = t.TempDir() // Automatically cleaned up after test

	t.Logf("workdir: %s", cmd.WorkDir)

	must.Error(ErrNoCommand, sub.Init(nil, []string{}))
	must.NoError(sub.Init(cmd, []string{"--no-batch"}))

	must.NoError(sub.Run(), "run failed")

	output := buf.String()
	t.Log("output:", output)

	must.Contains(output, "Templates created: "+cmd.WorkDir)

	var expectedFilePath string
	var actual []byte
	var err error

	filename := "batch.json"
	expectedFilePath = filepath.Join(cmd.WorkDir, filename)
	must.NoFileExists(expectedFilePath, "%s should not be present", filename)

	filename = "templates/custom.tpl"
	expectedFilePath = filepath.Join(cmd.WorkDir, filename)
	must.NoFileExists(expectedFilePath, "%s should not be present", filename)

	files := map[string]string{
		"variables.json":        ExampleJson,
		"variables.env":         ExampleEnv,
		"templates/example.tpl": ExampleTemplate,
	}

	for filename, expectedContents := range files {
		expectedFilePath = filepath.Join(cmd.WorkDir, filename)
		must.FileExists(expectedFilePath, "%s not present", filename)
		actual, err = os.ReadFile(expectedFilePath)
		must.NoError(err)
		must.Equal(expectedContents, string(actual), "%s file contents mismatch", filename)
	}
}
