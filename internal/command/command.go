package command

// Command processing common logic

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

const (
	SubCommandVersion = "version"
	SubCommandInit    = "init"
	SubCommandHelp    = "help"
	SubCommandBuild   = "build"
)

var ErrNoCommand = errors.New("command not defined")
var ErrNoInit = errors.New("command not initialized")

type Nillable interface {
	// IsNil check if interface has nil value
	IsNil() bool
}

// Subcommand is a subcommand common interface
type Subcommand interface {
	Nillable

	// Init boots command
	//
	// Arguments:
	//   - cmd is a parent command of the sub-command
	//   - args sub-command input arguments
	Init(cmd *Command, args []string) error

	// Name reads subcommand name
	Name() string

	// Run processes subcommand after Init was run
	Run() error

	// Usage show command usage. See flag.Usage
	Usage() error
}

// Command CLI command input options
type Command struct {
	App Application

	// Debug mode
	Debug bool

	// Disable CLI output colors
	NoColor bool

	//// Application environment
	//Environment string

	// Quiet suppresses STDOUT info messages
	Quiet bool

	// Suppress STDOUT info messages
	Verbose bool

	//// Clear ENV params before loading .env files
	//ClearEnv bool
	//
	//// List of DotEnv files to use instead of default ones
	//EnvFiles CmdInputFiles

	// Configuration for processing user resolution and other
	// ConfigFile string

	//// Configuration for processing user resolution and other
	//WorkDir string

	// Name is a name for the command. E.g., os.Args[0]
	Name string

	// Args contain all command arguments
	Args []string

	// fs results of parsing CLI command arguments excluding Subcommand
	fs *flag.FlagSet

	// Subcommand is a subcommand to run
	Subcommand Subcommand

	//// Writer CLI commands. If not defined then use os.Stdout
	//Writer io.Writer

	// Output is the stream to write output to
	Output io.Writer

	// Fmt styler of output
	Fmt *PrinterFormatter

	commands []Subcommand
}

func (c *Command) Init() error {
	c.fs = flag.NewFlagSet(c.Name, flag.ContinueOnError)
	c.fs.SetOutput(c.Output)
	c.fs.BoolVar(&c.NoColor, "nocolor", false, "disable color and styles output")
	c.fs.BoolVar(&c.Debug, "debug", false, "debug mode")
	c.fs.BoolVar(&c.Verbose, "verbose", false, "verbose output")

	c.commands = append(
		c.commands,
		&VersionSubcommand{},
	)

	var err error
	for _, sub := range c.commands {
		if err = sub.Init(c, nil); err != nil {
			return fmt.Errorf("%s %s init: %w", c.Name, sub.Name(), err)
		}
	}

	return nil
}

func (c *Command) Usage() error {
	if c.fs == nil {
		return errors.New(`argument flags is not defined`)
	}

	c.Fmt.Printf("Usage: <debug>%s [OPTIONS] COMMAND [COMMAND_ARGS]<reset>\n\n", c.Name)

	c.Fmt.Println(`<info>Arguments:<reset>`)
	c.Fmt.Printf("  <debug>%-10s<reset>\tshow this help\n", SubCommandHelp)
	c.Fmt.Printf("  <debug>%-10s<reset>\tshow application information\n", SubCommandVersion)
	c.Fmt.Printf("  <debug>%-10s<reset>\tinit default files structure for building templates\n", SubCommandInit)
	c.Fmt.Printf("  <debug>%-10s<reset>\trenders files from templates and configs\n", SubCommandBuild)
	c.Fmt.Println(``)

	c.Fmt.Println(`<info>Options:<reset>`)
	c.fs.PrintDefaults()
	c.Fmt.Println(``)

	c.Fmt.Println(`<info>Commands:<reset>`)

	var err error
	for _, sub := range c.commands {
		if err = sub.Usage(); err != nil {
			return fmt.Errorf("%s: %w", sub.Name(), err)
		}
	}

	return nil
}

func (c *Command) Run() error {
	var sub Subcommand
	var cmdArgs []string
	var subArgs []string

	for k, arg := range c.Args {
		for _, sc := range c.commands {
			if sc.Name() == arg {
				sub = sc
				subArgs = c.Args[k+1:]
				cmdArgs = c.Args[0:k]
			}
		}
	}

	var err error
	if sub == nil {
		//c.Usage()
		// TODO: help sub cmd show

		if err = c.fs.Parse(c.Args); err != nil {
			return fmt.Errorf("%s parse flags: %w", c.Name, err)
		}

		// update formatter coloring scheme
		c.Fmt.NoColor = c.NoColor
		c.Fmt.Init()

		if err = c.Usage(); err != nil {
			return fmt.Errorf("%s usage: %w", c.Name, err)
		}
	} else {
		//if err = fs.Parse(c.Args[0:subIndex]); err != nil {
		if err = c.fs.Parse(cmdArgs); err != nil {
			return fmt.Errorf("%s %s parse flags: %w", c.Name, sub.Name(), err)
		}

		// update formatter coloring scheme
		c.Fmt.NoColor = c.NoColor
		c.Fmt.Init()

		if err = sub.Init(c, subArgs); err != nil {
			return fmt.Errorf("%s %s init: %w", c.Name, sub.Name(), err)
		}

		if err = sub.Run(); err != nil {
			return fmt.Errorf("%s %s run: %w", c.Name, sub.Name(), err)
		}
	}

	// TODO: show help for parent and all sub-commands
	return nil
}

// NewCommandOpts defines options for NewCommand function
type NewCommandOpts struct {
	// Name is the name of binary command
	Name string

	// Args lists the command arguments
	Args []string

	// Output is the stream to write output to
	Output io.Writer

	// NoColor disables coloring and styling
	NoColor bool

	// App is an application for running command
	App Application
}

// NewCommand creates new command
func NewCommand(opts NewCommandOpts) *Command {
	return &Command{
		Name:   opts.Name,
		Args:   opts.Args,
		Output: opts.Output,
		Fmt:    NewPrinterFormatter(opts.NoColor, opts.Output),
		App:    opts.App,
	}
}
