package command

import (
	"flag"
	"os"
	"path/filepath"
)

type InitSubcommand struct {
	cmd  *Command
	fs   *flag.FlagSet
	Help bool
}

func (c *InitSubcommand) Name() string {
	return SubCommandInit
}

func (c *InitSubcommand) usage() {
	if c.fs == nil {
		panic(ErrNoInit)
	}

	subName := c.Name()
	c.cmd.Fmt.Printf("<debug>%-15s<reset> init default files structure for building templates\n\n", subName)
	c.cmd.Fmt.Printf("Usage: <debug>%s [OPTIONS] %s [COMMAND_OPTIONS]<reset>\n\n", c.cmd.Name, subName)

	c.cmd.Fmt.Println(`<info>Options:<reset>`)
	c.fs.PrintDefaults()
	c.cmd.Fmt.Println(``)

	c.cmd.Fmt.Println("<info>Examples:<reset>")
	c.cmd.Fmt.Printf("  $ %-40s # init project under current working directory\n", c.cmd.Name+" init")
	c.cmd.Fmt.Printf("  $ %-40s # show command usage help\n", c.cmd.Name+" init -h")
	c.cmd.Fmt.Printf("  $ %-40s # init project in custom working directory\n", c.cmd.Name+" --workdir ~/.templar init")
}

func (c *InitSubcommand) Usage() error {
	if c.fs == nil {
		return ErrNoInit
	}

	c.usage()

	return nil
}

func (c *InitSubcommand) Init(cmd *Command, args []string) error {
	if cmd == nil {
		return ErrNoCommand
	}

	c.cmd = cmd
	c.fs = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.fs.BoolVar(&c.Help, `h`, false, `print command usage suboptions`)
	c.fs.SetOutput(c.cmd.Output)
	c.fs.Usage = c.usage

	return c.fs.Parse(args)
}

func (c *InitSubcommand) IsNil() bool {
	return c == nil
}

func (c *InitSubcommand) Run() error {
	if c.Help {
		return c.Usage()
	}

	//c.cmd.Fmt.Printf("%s %+v\n", c.cmd.Name, c)

	// TODO: create mkdir -p if not exists. Check custom workdir c.cmd.WorkDir

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

	if err = os.WriteFile(c.cmd.WorkDir+`/variables.env`, []byte(ExampleEnv), 0644); err != nil {
		return err
	}

	if err = os.WriteFile(c.cmd.WorkDir+`/variables.json`, []byte(ExampleJson), 0644); err != nil {
		return err
	}

	if err = os.WriteFile(c.cmd.WorkDir+`/templates/example.tmpl`, []byte(ExampleTemplate), 0644); err != nil {
		return err
	}

	if !c.cmd.Quiet {
		c.cmd.Fmt.Printf("Templates created: <debug>%s<reset>\n", path)
	}

	return nil
}
