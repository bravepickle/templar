package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCommand(t *testing.T) {
	must := require.New(t)

	datasets := []struct {
		name            string
		command         string
		args            []string
		version         string
		commit          string
		workdir         string
		expectedError   string
		expectedVersion string
		expectedCommit  string
		expectedWorkdir string
	}{
		{
			name:            "show version",
			command:         "test-app",
			args:            []string{"--nocolor", "version"},
			version:         "v1.0.0",
			commit:          "777",
			workdir:         "/tmp",
			expectedError:   "",
			expectedVersion: "v1.0.0",
			expectedCommit:  "777",
			expectedWorkdir: "/tmp",
		},
	}

	for _, d := range datasets {
		t.Run(d.name, func(t *testing.T) {
			buf := bytes.NewBufferString("")
			err := RunCommand(d.command, d.args, buf, d.version, d.commit, d.workdir)

			if d.expectedError != "" {
				must.ErrorContains(err, d.expectedError)
			} else {
				must.NoError(err)

				output := buf.String()
				t.Logf("Output: %s", output)

				must.Contains(output, d.command+":", "name mismatch")
				must.Contains(output, "Version: "+d.expectedVersion, "version mismatch")
				must.Contains(output, "GIT commit: "+d.expectedCommit, "commit mismatch")
				must.Contains(output, "Working directory: "+d.workdir, "workdir mismatch")
			}
		})
	}
}
