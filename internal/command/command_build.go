package command

import (
	"errors"
	"flag"
)

type BuildCommand struct {
	cmd          *Command
	fs           *flag.FlagSet
	InputFile    string
	OutputFile   string
	InputFormat  string
	TemplateFile string
	BatchFile    string
}

func (c *BuildCommand) Name() string {
	return SubCommandBuild
}

func (c *BuildCommand) usage() {
	if c.fs == nil {
		panic(ErrNoInit)
	}

	subName := c.Name()
	c.cmd.Fmt.Printf("<debug>%-15s<reset> %s\n\n", subName, c.Summary())
	c.cmd.Fmt.Printf("Usage: <debug>%s %s [OPTIONS]<reset>\n", c.cmd.Name, subName)
	c.cmd.Fmt.Println(``)

	c.cmd.Fmt.Println("<info>Options:<reset>")
	c.fs.PrintDefaults()
	c.cmd.Fmt.Println(``)

	c.cmd.Fmt.Printf("<info>Examples:<reset>\n  $ %s %s\n\n", c.cmd.Name, subName)
	c.cmd.Fmt.Println("  TBD\n") // <<<<<<<<<<<<<<<<<<<<<!!!!!

}

func (c *BuildCommand) Summary() string {
	return "render template contents with provided variables"
}

func (c *BuildCommand) Usage() error {
	if c.fs == nil {
		return ErrNoInit
	}

	c.usage()

	return nil
}

func (c *BuildCommand) Init(cmd *Command, args []string) error {
	if cmd == nil {
		return ErrNoCommand
	}

	c.cmd = cmd
	c.fs = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.fs.SetOutput(c.cmd.Output)
	c.fs.Usage = c.usage

	format := c.cmd.Fmt.Sprintf

	c.fs.StringVar(&c.InputFile, "in", "", format("input file path. Format should match \"<debug>-format<reset>\" value"))
	c.fs.StringVar(&c.InputFormat, "format", "env", "input file format. Allowed: env, json")
	c.fs.StringVar(&c.OutputFile, "out", "", format("output file path, If empty, outputs to stdout. If \"<debug>-batch<reset>\" option is used, specifies output directory"))
	c.fs.StringVar(&c.TemplateFile, "template", "", format("template file path, If empty and \"<debug>-batch<reset>\" not defined, reads from stdin"))
	c.fs.StringVar(&c.BatchFile, "batch", "", format("build multiple files from templates. Supersedes \"<debug>-template<reset>\", \"<debug>-in<reset>\" options. See examples for details"))

	return c.fs.Parse(args)
}

func (c *BuildCommand) IsNil() bool {
	return c == nil
}

func (c *BuildCommand) Run() error {
	c.cmd.Fmt.Printf("<alert><bold>TBD. Needs implementation<reset>\n")

	return errors.New("implement \"build\" command")
}
