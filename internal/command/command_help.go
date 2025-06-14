package command

import (
	"flag"
)

type HelpCommand struct {
	cmd *Command
	fs  *flag.FlagSet
}

func (c *HelpCommand) Name() string {
	return SubCommandHelp
}

func (c *HelpCommand) usage() {
	if c.fs == nil {
		panic(ErrNoInit)
	}

	subName := c.Name()
	c.cmd.Fmt.Printf("Usage: <debug>%s [OPTIONS] %s [COMMAND]<reset>\n\n", c.cmd.Name, subName)
	c.cmd.Fmt.Printf("<debug>%-10s<reset> %s\n\n", subName, c.Summary())

	c.cmd.Fmt.Printf("<info>Examples:<reset>\n  <debug>$ %s %s %s<reset>\n\n", c.cmd.Name, subName, SubCommandVersion)
	c.cmd.Fmt.Printf("  %s:\n    Version: v0.0.1\n    GIT commit: c7a8949\n    Working directory: /home/user/templar\n", c.cmd.Name)
}

func (c *HelpCommand) Summary() string {
	return c.cmd.Fmt.Sprintf("show help information on command or subcommand usage. "+
		"Type \"<debug>%s %s %s<reset>\" to see help command usage information", c.cmd.Name, c.Name(), c.Name())
}

func (c *HelpCommand) Usage() error {
	if c.fs == nil {
		return ErrNoInit
	}

	c.usage()

	return nil
}

func (c *HelpCommand) Init(cmd *Command, args []string) error {
	if cmd == nil {
		return ErrNoCommand
	}

	c.cmd = cmd
	c.fs = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.fs.SetOutput(c.cmd.Output)
	c.fs.Usage = c.usage

	return c.fs.Parse(args)
}

func (c *HelpCommand) IsNil() bool {
	return c == nil
}

func (c *HelpCommand) Run() error {
	if c.fs == nil || !c.fs.Parsed() {
		return ErrNoInit
	}

	targetCmd := c.fs.Arg(0)
	if targetCmd == "" {
		return c.cmd.Usage()
	}

	if targetCmd == SubCommandHelp {
		return c.Usage()
	}

	for _, sub := range c.cmd.commands {
		if sub.Name() == targetCmd {
			return sub.Usage()
		}
	}

	return c.cmd.Usage()
}
