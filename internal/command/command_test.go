package command

import (
	"bytes"

	"github.com/bravepickle/templar/internal/core"
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

	must.NotNil(cmd)
	must.NoError(cmd.Init(), "init failed")

	for _, c := range cmd.commands {
		if c.Name() == targetCmd {
			return c, cmd
		}
	}

	return nil, cmd
}
