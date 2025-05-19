package command

import (
	"flag"
)

type VersionSubcommand struct {
	cmd *Command
	fs  *flag.FlagSet
}

func (c *VersionSubcommand) Name() string {
	return SubCommandVersion
}

func (c *VersionSubcommand) usage() {
	//if c.fs == nil {
	//	panic(ErrNoInit)
	//}

	c.cmd.Fmt.Printf("<debug>%-15s<reset> show application information on its build version and directories\n\n", c.Name())
	c.cmd.Fmt.Printf("<comment>Examples:<reset>\n  $ %s %s\n\n", c.cmd.Name, c.Name())
	c.cmd.Fmt.Println("  templar:\n    Version: v0.0.1\n    GIT commit: c7a8949\n    Working directory:   /usr/local/bin/templar\n")
	//c.fs.PrintDefaults()
}

func (c *VersionSubcommand) Usage() error {
	if c.fs == nil {
		return ErrNoInit
	}

	c.usage()

	return nil
}

func (c *VersionSubcommand) Init(cmd *Command, _ []string) error {
	if cmd == nil {
		return ErrNoCommand
	}

	c.cmd = cmd
	c.fs = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.fs.SetOutput(c.cmd.Output)
	c.fs.Usage = c.usage

	return nil
}

func (c *VersionSubcommand) IsNil() bool {
	return c == nil
}

func (c *VersionSubcommand) Run() error {
	c.cmd.Fmt.Printf("<info><bold>%s:<reset>\n", c.cmd.Name)
	c.cmd.Fmt.Printf("  <debug>Version:<reset> %s\n", c.cmd.App.Version)
	c.cmd.Fmt.Printf("  <debug>GIT commit:<reset> %s\n", c.cmd.App.GitCommitHash)
	c.cmd.Fmt.Printf("  <debug>Working directory:<reset> %s\n", c.cmd.App.WorkDir)

	return nil
}
