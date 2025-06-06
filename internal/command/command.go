package command

// Command processing common logic

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bravepickle/templar/internal/core"
)

const (
	SubCommandVersion = "version"
	SubCommandInit    = "init"
	SubCommandHelp    = "help"
	SubCommandBuild   = "build"
	//SubCommandVars    = "vars"
)

var ErrNoCommand = errors.New("command not defined")
var ErrNoInit = errors.New("command not initialized")

// Subcommand is a subcommand common interface
type Subcommand interface {
	core.Nillable

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

	// Usage shows command extended usage description with examples. See flag.Usage
	Usage() error

	// Summary generates short description for the command
	Summary() string
}

// Command CLI command input options
type Command struct {
	App core.Application

	// Debug mode
	Debug bool

	// Disable CLI output colors
	NoColor bool

	// Quiet suppresses STDOUT info messages
	Quiet bool

	// Suppress STDOUT info messages
	Verbose bool

	// Configuration for processing user resolution and other
	// ConfigFile string

	// DefaultWorkDir is a default working directory
	DefaultWorkDir string

	// WorkDir is a selected working directory
	WorkDir string

	// Name is a name for the command. E.g., os.Args[0]
	Name string

	// Args contain all command arguments
	Args []string

	// fs results of parsing CLI command arguments excluding Subcommand
	fs *flag.FlagSet

	// Output is the stream to write output to
	Output io.Writer

	// Input streams input to command for interaction.
	// E.g. STDIN, file, keyboard etc.
	Input *os.File

	// Fmt styler of output
	Fmt *core.PrinterFormatter

	commands map[string]Subcommand
}

func (c *Command) Init() error {
	if c.DefaultWorkDir == "" {
		if pwd, err := os.Getwd(); err != nil {
			return err
		} else {
			c.DefaultWorkDir = pwd
		}
	}

	c.fs = flag.NewFlagSet(c.Name, flag.ContinueOnError)
	c.fs.SetOutput(c.Output)
	c.fs.BoolVar(&c.NoColor, "no-color", false, "disable color and styles output")
	c.fs.BoolVar(&c.Debug, "debug", false, "debug mode")
	c.fs.BoolVar(&c.Verbose, "verbose", false, "verbose output")
	c.fs.StringVar(&c.WorkDir, "workdir", c.DefaultWorkDir, "working directory path")

	if c.Input == nil {
		c.Input = os.Stdin // defaults
	}

	c.commands = map[string]Subcommand{
		SubCommandHelp:    &HelpCommand{},
		SubCommandVersion: &VersionCommand{},
		SubCommandInit:    &InitCommand{},
		SubCommandBuild:   &BuildCommand{In: c.Input},
	}

	var err error
	for _, sub := range c.commands {
		if err = sub.Init(c, nil); err != nil {
			return fmt.Errorf("%s %s init: %w", c.Name, sub.Name(), err)
		}
	}

	return nil
}

func (c *Command) Summary() string {
	return "generate template contents with provided variables"
}

func (c *Command) Usage() error {
	if c.fs == nil {
		return errors.New(`argument flags is not defined`)
	}

	c.Fmt.Printf("Usage: <debug>%s [OPTIONS] COMMAND [COMMAND_ARGS]<reset>\n\n", c.Name)
	c.Fmt.Printf("<debug>%-10s<reset> %s\n\n", c.Name, c.Summary())

	c.Fmt.Println(`<info>Commands:<reset>`)
	for _, sub := range c.commands {
		c.Fmt.Printf("  <debug>%-10s<reset> %s\n", sub.Name(), sub.Summary())
	}
	c.Fmt.Println(``)

	c.Fmt.Println(`<info>Options:<reset>`)
	c.fs.PrintDefaults()
	c.Fmt.Println(``)

	return nil
}

func (c *Command) Run() error {
	var sub Subcommand
	var cmdArgs []string
	var subArgs []string

loop:
	for k, arg := range c.Args {
		if strings.HasPrefix(arg, "-") {
			continue
		}

		for _, sc := range c.commands {
			if sc.Name() == arg {
				sub = sc
				subArgs = c.Args[k+1:]
				cmdArgs = c.Args[0:k]

				break loop
			}
		}
	}

	var err error
	if sub == nil {
		if err = c.fs.Parse(c.Args); err != nil {
			return fmt.Errorf("%s parse flags: %w", c.Name, err)
		}

		// update formatter coloring scheme
		c.Fmt.NoColor = c.NoColor
		c.Fmt.Init()

		if err = c.Init(); err != nil {
			return fmt.Errorf("%s usage: %w", c.Name, err)
		}

		if err = c.Usage(); err != nil {
			return fmt.Errorf("%s usage: %w", c.Name, err)
		}
	} else {
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

	// Input input data stream
	Input *os.File

	// NoColor disables coloring and styling
	NoColor bool

	// WorkDir working directory
	WorkDir string

	// App is an application for running command
	App core.Application
}

// NewCommand creates new command
func NewCommand(opts NewCommandOpts) *Command {
	return &Command{
		Name:           opts.Name,
		Args:           opts.Args,
		Input:          opts.Input,
		Output:         opts.Output,
		DefaultWorkDir: opts.WorkDir,
		Fmt:            core.NewPrinterFormatter(opts.NoColor, opts.Output),
		App:            opts.App,
	}
}
