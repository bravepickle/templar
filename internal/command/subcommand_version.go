package command

import (
	"flag"
	"os"
)

// Subcommand is a subcommand common interface
//type Subcommand interface {
//	Nillable
//
//	// Init boots command
//	//
//	// Arguments:
//	//   - name sub-command name
//	//   - args sub-command input arguments
//	//   - cmd is a parent command of the sub-command
//	Init(name string, args []string, cmd Command)
//
//	// Name reads subcommand name
//	Name() string
//
//	// Run processes subcommand after Init was run
//	Run() error
//}

type VersionSubcommand struct {
	cmd *Command
	fs  *flag.FlagSet
}

func (c *VersionSubcommand) Name() string {
	return SubCommandVersion
}

func (c *VersionSubcommand) usage() {
	if c.cmd == nil {
		panic(ErrNoCommand)
	}

}

func (c *VersionSubcommand) Usage() (string, error) {
	if c.fs == nil {
		return "", ErrNoInit
	}

	c.fs.Usage()

	return "", nil
}

func (c *VersionSubcommand) Init(cmd *Command, _ []string) error {
	if cmd == nil {
		return ErrNoCommand
	}

	c.cmd = cmd
	c.fs = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.fs.Usage = c.usage

	return nil
}

func (c *VersionSubcommand) IsNil() bool {
	return c == nil
}

func (c *VersionSubcommand) Run() error {
	c.cmd.Fmt.Printf("<info><bold>%s:<reset>\n", c.cmd.Name)
	c.cmd.Fmt.Printf("  <debug>Version:<reset> %s\n", c.cmd.AppVersion)
	c.cmd.Fmt.Printf("  <debug>GIT commit:<reset> %s\n", c.cmd.GitCommitHash)

	if c.cmd.App.WorkDir == "" {
		c.cmd.App.WorkDir, _ = os.Getwd()
	}

	return nil
}
