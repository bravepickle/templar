package command

import (
	"flag"
	"os"
	"path/filepath"
)

type InitCommand struct {
	cmd *Command
	fs  *flag.FlagSet

	// NoBatch disables generation of batch file examples
	NoBatch bool
}

func (c *InitCommand) Name() string {
	return SubCommandInit
}

func (c *InitCommand) usage() {
	if c.fs == nil {
		panic(ErrNoInit)
	}

	subName := c.Name()
	c.cmd.Fmt.Printf("Usage: <debug>%s [OPTIONS] %s [COMMAND_OPTIONS]<reset>\n\n", c.cmd.Name, subName)
	c.cmd.Fmt.Printf("<debug>%-10s<reset> %s\n\n", subName, c.Summary())

	c.cmd.Fmt.Println(`<info>Options:<reset>`)
	c.fs.PrintDefaults()
	c.cmd.Fmt.Println(``)

	c.cmd.Fmt.Println("<info>Examples:<reset>")
	c.cmd.Fmt.Printf("  <debug>$ %-40s<reset> # init project under current working directory\n", c.cmd.Name+" init")
	c.cmd.Fmt.Printf("  <debug>$ %-40s<reset> # init project in custom working directory\n", c.cmd.Name+" --workdir ~/.templar init")
}

func (c *InitCommand) Usage() error {
	if c.fs == nil {
		return ErrNoInit
	}

	c.usage()

	return nil
}

func (c *InitCommand) Summary() string {
	return "init default files structure for building templates"
}

func (c *InitCommand) Init(cmd *Command, args []string) error {
	if cmd == nil {
		return ErrNoCommand
	}

	c.cmd = cmd
	c.fs = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.fs.SetOutput(c.cmd.Output)
	c.fs.BoolVar(&c.NoBatch, "no-batch", false, "skip batch examples generation")
	c.fs.Usage = c.usage

	return c.fs.Parse(args)
}

func (c *InitCommand) IsNil() bool {
	return c == nil
}

func (c *InitCommand) Run() error {
	var err error
	if c.cmd.WorkDir == "" {
		if c.cmd.WorkDir, err = os.Getwd(); err != nil {
			return err
		}
	}

	if _, err = os.Stat(c.cmd.WorkDir); os.IsNotExist(err) {
		if err = os.Mkdir(c.cmd.WorkDir, MkDirPerm); err != nil {
			return err
		}
	}

	path, err := filepath.Abs(c.cmd.WorkDir + `/templates`)
	if err != nil {
		return err
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, MkDirPerm)
		if err != nil {
			return err
		}

		if c.cmd.Verbose {
			c.cmd.Fmt.Println(`Directory created:`, path)
		}
	} else {
		if c.cmd.Verbose {
			c.cmd.Fmt.Println(`Directory already created:`, path)
		}
	}

	if err = os.WriteFile(c.cmd.WorkDir+`/variables.env`, []byte(ExampleEnv), MkFilePerm); err != nil {
		return err
	}

	if err = os.WriteFile(c.cmd.WorkDir+`/variables.json`, []byte(ExampleJson), MkFilePerm); err != nil {
		return err
	}

	if err = os.WriteFile(c.cmd.WorkDir+`/templates/example.tpl`, []byte(ExampleTemplate), MkFilePerm); err != nil {
		return err
	}

	if !c.NoBatch {
		if err = os.WriteFile(c.cmd.WorkDir+`/batch.json`, []byte(ExampleBatchJson), MkFilePerm); err != nil {
			return err
		}

		if err = os.WriteFile(c.cmd.WorkDir+`/templates/custom.tpl`, []byte(ExampleCustomTemplate), MkFilePerm); err != nil {
			return err
		}
	}

	if !c.cmd.Quiet {
		c.cmd.Fmt.Printf("Templates created: <debug>%s<reset>\n", path)
	}

	return nil
}
