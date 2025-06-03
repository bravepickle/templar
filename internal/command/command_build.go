package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bravepickle/templar/internal/core"
	"github.com/bravepickle/templar/internal/parser"
)

type BuildCommand struct {
	cmd          *Command
	fs           *flag.FlagSet
	InputFile    string
	OutputFile   string
	InputFormat  string
	TemplateFile string
	BatchFile    string
	SkipExisting bool
	ClearEnv     bool
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
	c.cmd.Fmt.Printf("Usage: <debug>%s [OPTIONS] %s [COMMAND_OPTIONS]<reset>\n", c.cmd.Name, subName)
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

	c.fs.StringVar(&c.InputFile, "input", "", "input file path. Format should match \"-format\" value")
	c.fs.StringVar(&c.InputFormat, "format", "env", "input file format. Allowed: env, json")
	c.fs.StringVar(&c.OutputFile, "output", "", "output file path, If empty, outputs to stdout. If \"-batch\" option is used, specifies output directory")
	c.fs.StringVar(&c.TemplateFile, "template", "", "template file path, If empty and \"-batch\" not defined, reads from stdin")
	c.fs.StringVar(&c.BatchFile, "batch", "", "batch file path. Overrides some other fields, such as --variables")
	c.fs.BoolVar(&c.SkipExisting, "skip", false, "skip generation if target files already exist")
	c.fs.BoolVar(&c.ClearEnv, "clear", false, "clear ENV variables before building variables to avoid collisions")

	return c.fs.Parse(args)
}

func (c *BuildCommand) IsNil() bool {
	return c == nil
}

func (c *BuildCommand) Run() error {
	if c.fs == nil {
		return ErrNoInit
	}

	if c.BatchFile != "" {
		return c.runBatch()
	}

	//c.cmd.Fmt.Printf("<alert><bold>TBD. Needs implementation<reset>\n")

	return c.runOnce()
}

func (c *BuildCommand) readTemplate() (string, error) {
	var tplContents []byte
	var err error

	if c.TemplateFile == "" {
		tplContents, err = io.ReadAll(os.Stdin) // read STDIN
	} else {
		tplFile := c.TemplateFile
		if !filepath.IsAbs(c.TemplateFile) {
			tplFile = filepath.Join(c.cmd.WorkDir, c.TemplateFile)
		}

		tplContents, err = os.ReadFile(tplFile)
	}

	if err != nil {
		return "", err
	}

	return string(tplContents), nil
}

func (c *BuildCommand) readVars() (parser.Params, error) {
	var contents []byte
	var err error

	//var varParser parser.Parser

	//if c.InputFile == "" {
	//	contents, err = io.ReadAll(os.Stdin) // read STDIN
	//} else {
	//	tplFile := c.TemplateFile
	//	if !filepath.IsAbs(c.TemplateFile) {
	//		tplFile = filepath.Join(c.cmd.WorkDir, c.TemplateFile)
	//	}
	//
	//	contents, err = os.ReadFile(tplFile)
	//}

	if c.InputFile != "" {
		if filepath.IsAbs(c.InputFile) {
			contents, err = os.ReadFile(c.InputFile)
		} else {
			contents, err = os.ReadFile(filepath.Join(c.cmd.WorkDir, c.InputFile))
		}

		if err != nil {
			return nil, fmt.Errorf(`variables file: %w`, err)
		}
	}

	var varParser parser.Parser

	switch c.InputFormat {
	case "env":
		varParser = c.getEnvParser(contents)
	case "json":
		varParser = c.getJSONParser(contents)
	default:
		return nil, fmt.Errorf("invalid input format: %s", c.InputFormat)
	}

	if varParser != nil && !varParser.IsNil() {
		return varParser.Parse(string(contents))
	}

	return nil, nil // everything is fine but ono vars input found
}

func (c *BuildCommand) getEnvParser(contents []byte) parser.Parser {
	if c.ClearEnv {
		if len(contents) > 0 {
			return parser.NewEnvParser()
		}
	} else {
		if len(contents) > 0 {
			return parser.NewChainParser(
				parser.NewEnvParser(),
				parser.NewEnvOsParser(),
			)
		} else {
			return parser.NewEnvOsParser()
		}
	}

	return nil
}

func (c *BuildCommand) getJSONParser(contents []byte) parser.Parser {
	if c.ClearEnv {
		if len(contents) > 0 {
			return parser.NewJSONParser()
		}
	} else {
		if len(contents) > 0 {
			return parser.NewChainParser(
				parser.NewJSONParser(),
				parser.NewEnvOsParser(),
			)
		} else {
			return parser.NewEnvOsParser()
		}
	}

	return nil
}

func (c *BuildCommand) selectWriter() (io.Writer, error) {
	if c.OutputFile == "" {
		return os.Stdout, nil
	}

	outputFile := c.OutputFile
	if !filepath.IsAbs(outputFile) {
		outputFile = filepath.Join(c.cmd.WorkDir, c.OutputFile)
	}

	return os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, MkFilePerm)
}

func (c *BuildCommand) runOnce() error {
	tplContents, err := c.readTemplate()
	if err != nil {
		return fmt.Errorf("template read: %w", err)
	}

	//var varParser parser.Parser
	//
	//switch c.InputFormat {
	//case "env":
	//	if c.ClearEnv {
	//		varParser = parser.NewEnvParser()
	//	} else {
	//		varParser = parser.NewChainParser(
	//			parser.NewEnvParser(),
	//			parser.NewEnvOsParser(),
	//		)
	//	}
	//case "json":
	//	if c.ClearEnv {
	//		varParser = parser.NewJSONParser()
	//	} else {
	//		varParser = parser.NewChainParser(
	//			parser.NewJSONParser(),
	//			parser.NewEnvOsParser(),
	//		)
	//	}
	//default:
	//	return fmt.Errorf("invalid input format: %s", c.InputFormat)
	//}

	var params parser.Params

	params, err = c.readVars()
	if err != nil {
		return fmt.Errorf("variables read: %w", err)
	}

	writer, err := c.selectWriter()
	if err != nil {
		return fmt.Errorf("select writer: %w", err)
	}

	if oc, ok := writer.(io.Closer); ok {
		defer oc.Close()
	}

	//fmt.Println("template contents:", tplContents)
	//fmt.Printf("variables: %+v\n", params)

	builder := parser.NewTemplate(c.TemplateFile, tplContents, params)

	return builder.Build(writer)
}

func (c *BuildCommand) runBatch() error {
	contents, err := os.ReadFile(c.BatchFile)
	if err != nil {
		return err
	}

	var batch *core.Batch
	if err := json.Unmarshal(contents, batch); err != nil {
		return err
	}

	c.cmd.Fmt.Printf("<alert><bold>TBD: contents - %+v<reset>\n", batch)

	return nil
}
