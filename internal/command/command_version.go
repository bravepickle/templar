package command

import (
	"flag"
)

type VersionCommand struct {
	cmd *Command
	fs  *flag.FlagSet
}

func (c *VersionCommand) Name() string {
	return SubCommandVersion
}

func (c *VersionCommand) usage() {
	if c.fs == nil {
		panic(ErrNoInit)
	}

	subName := c.Name()
	c.cmd.Fmt.Printf("Usage: <debug>%s [OPTIONS] %s<reset>\n\n", c.cmd.Name, subName)
	c.cmd.Fmt.Printf("<debug>%-10s<reset> %s\n\n", subName, c.Summary())

	c.cmd.Fmt.Printf("<debug>Examples:<reset>\n  <comment>$ %s %s<reset>\n\n", c.cmd.Name, subName)
	c.cmd.Fmt.Printf("  %s:\n    Version: v0.0.1\n    GIT commit: c7a8949\n    Working directory: /usr/local/bin/templar\n", c.cmd.Name)
	//c.fs.PrintDefaults()
}

func (c *VersionCommand) Summary() string {
	return "show application information on its build version and directories"
}

func (c *VersionCommand) Usage() error {
	if c.fs == nil {
		return ErrNoInit
	}

	c.usage()

	return nil
}

func (c *VersionCommand) Init(cmd *Command, args []string) error {
	if cmd == nil {
		return ErrNoCommand
	}

	c.cmd = cmd
	c.fs = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.fs.SetOutput(c.cmd.Output)
	c.fs.Usage = c.usage

	return c.fs.Parse(args)
}

func (c *VersionCommand) IsNil() bool {
	return c == nil
}

func (c *VersionCommand) Run() error {
	c.cmd.Fmt.Printf("<info><bold>%s:<reset>\n", c.cmd.Name)
	c.cmd.Fmt.Printf("  <debug>Version:<reset> %s\n", c.cmd.App.Version)
	c.cmd.Fmt.Printf("  <debug>GIT commit:<reset> %s\n", c.cmd.App.GitCommitHash)
	c.cmd.Fmt.Printf("  <debug>Working directory:<reset> %s\n", c.cmd.WorkDir)

	return nil
}
