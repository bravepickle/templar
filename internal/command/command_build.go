package command

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bravepickle/templar/internal/core"
	"github.com/bravepickle/templar/internal/parser"
)

type BuildCommand struct {
	cmd *Command
	fs  *flag.FlagSet
	// In is the default stream to read input from for templates
	In io.Reader

	InputFile     string
	OutputFile    string
	InputFormat   string
	TemplateFile  string
	SkipExisting  bool
	ClearEnv      bool
	Dump          string
	NoCloseWriter bool
}

func (c *BuildCommand) Name() string {
	return SubCommandBuild
}

func (c *BuildCommand) usage() {
	if c.fs == nil {
		panic(ErrNoInit)
	}

	subName := c.Name()
	c.cmd.Fmt.Printf("Usage: <debug>%s [OPTIONS] %s [COMMAND_OPTIONS]<reset>\n\n", c.cmd.Name, subName)
	c.cmd.Fmt.Printf("<debug>%-10s<reset> %s\n\n", subName, c.Summary())

	c.cmd.Fmt.Println("<info>Options:<reset>")
	c.fs.PrintDefaults()
	c.cmd.Fmt.Println(``)

	c.cmd.Fmt.Println("<info>Examples:<reset>")
	c.cmd.Fmt.Printf(`  <debug>$ %[1]s build --input .env --format env --template template.tpl --output output.txt<reset> 
      # generates output.txt file from the provided template.tpl and .env variables in env format (is the default one, can be ommitted)

  <debug>$ NAME=John %[1]s build --template template.tpl --output output.txt<reset>
      # generates output.txt file from the provided template.tpl and provided env variable

  <debug>$ echo "My name is {{ .NAME }}" | NAME=John %[1]s build<reset>
      # generates output.txt file from the provided template.tpl and provided env variable
`, c.cmd.Name)

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

	c.fs.StringVar(&c.InputFile, "input", "", "file path which contains variables for template to use or batch file. Format should match \"-format\" value")
	c.fs.StringVar(&c.InputFormat, "format", "env", "input file format for variables' file. Allowed: "+strings.Join(AllowedFormats, ", "))
	c.fs.StringVar(&c.OutputFile, "output", "", "output file path, If empty, outputs to stdout. If \"-batch\" option is used, specifies output directory")
	c.fs.StringVar(&c.TemplateFile, "template", "", "template file path, If empty and \"-batch\" not defined, reads from stdin")
	c.fs.StringVar(&c.Dump, "dump", "env", "show all available variables for the template to use and stop processing. Pass optionally --verbose or --debug flags for more information. Allowed output formats: json, env")
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

	if c.InputFormat == FormatBatch {
		return c.runBatch()
	}

	return c.runOnce()
}

func (c *BuildCommand) readInput(path string) ([]byte, error) {
	var contents []byte
	var err error

	if path == "" {
		if c.In == nil {
			contents, err = io.ReadAll(os.Stdin) // default input stream
		} else {
			contents, err = io.ReadAll(c.In) // read from custom input io.Reader
		}
	} else {
		if !filepath.IsAbs(path) {
			path = filepath.Join(c.cmd.WorkDir, path)
		}

		contents, err = os.ReadFile(path)
	}

	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (c *BuildCommand) readVars() (core.Params, error) {
	var contents []byte
	var err error

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
	case FormatEnv:
		varParser = c.getEnvParser(len(contents) > 0)
	case FormatJson:
		varParser = c.getJSONParser(len(contents) > 0)
	default:
		return nil, fmt.Errorf("invalid input format: %s", c.InputFormat)
	}

	if varParser != nil && !varParser.IsNil() {
		return varParser.Parse(string(contents))
	}

	return nil, nil // everything is fine but ono vars input found
}

func (c *BuildCommand) getEnvParser(hasVars bool) parser.Parser {
	if c.ClearEnv {
		if hasVars {
			return parser.NewEnvParser()
		}
	} else {
		if hasVars {
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

func (c *BuildCommand) getJSONParser(hasVars bool) parser.Parser {
	if c.ClearEnv {
		if hasVars {
			return parser.NewJSONParser()
		}
	} else {
		if hasVars {
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

func (c *BuildCommand) selectWriter(outputFile string) (io.Writer, error) {
	if outputFile == "" {
		return c.cmd.Output, nil
	}

	if !filepath.IsAbs(outputFile) {
		outputFile = filepath.Join(c.cmd.WorkDir, outputFile)
	}

	return os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, MkFilePerm)
}

func (c *BuildCommand) runOnce() error {
	tplContents, err := c.readInput(c.TemplateFile)
	if err != nil {
		return fmt.Errorf("template read: %w", err)
	}

	var params core.Params

	params, err = c.readVars()
	if err != nil {
		return fmt.Errorf("variables read: %w", err)
	}

	if c.Dump != "" {
		return c.dumpParams(params)
	}

	writer, err := c.selectWriter(c.OutputFile)
	if err != nil {
		return fmt.Errorf("select writer: %w", err)
	}

	if !c.NoCloseWriter {
		if oc, ok := writer.(io.Closer); ok {
			defer oc.Close()
		}
	}

	builder := parser.NewTemplate(c.TemplateFile, string(tplContents), params)

	return builder.Build(writer)
}

func (c *BuildCommand) prepareVarsForDump(params core.Params) ([]string, map[string]string, map[string]any) {
	if len(params) == 0 {
		return nil, nil, nil
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	if c.cmd.Debug {
		data := map[string]any{}

		for _, k := range keys {
			data[k] = params[k]
		}

		return keys, nil, data
	}

	if c.cmd.Verbose {
		data := map[string]string{}

		for _, k := range keys {
			data[k] = fmt.Sprintf("%T", params[k])
		}

		return keys, data, nil
	}

	return keys, nil, nil
}

func (c *BuildCommand) dumpParams(params core.Params) error {
	keys, strMap, anyMap := c.prepareVarsForDump(params)

	if c.Dump == FormatJson {
		var data any
		if len(strMap) > 0 {
			data = strMap
		} else if len(anyMap) > 0 {
			data = anyMap
		} else if len(keys) > 0 {
			data = keys
		} else {
			return nil
		}

		if output, err := json.MarshalIndent(data, "", "  "); err != nil {
			return err
		} else {
			c.cmd.Fmt.PrintRaw(string(output) + "\n")

			return nil
		}
	}

	if len(keys) == 0 {
		if c.cmd.Debug || c.cmd.Verbose {
			c.cmd.Fmt.PrintfRaw("No variables found\n")
		}

		return nil
	}

	if len(strMap) > 0 {
		for _, k := range keys {
			c.cmd.Fmt.PrintfRaw("%s=%s\n", k, strMap[k])
		}

		return nil
	}

	if len(anyMap) > 0 {
		for _, k := range keys {
			if vm, ok := anyMap[k].(map[string]any); ok {
				if v, err := json.Marshal(vm); err == nil {
					c.cmd.Fmt.PrintfRaw("%s=%s\n", k, v)

					continue
				}
			}

			c.cmd.Fmt.PrintfRaw("%s=%#v\n", k, anyMap[k])
		}

		return nil
	}

	for _, k := range keys {
		c.cmd.Fmt.PrintRaw(k + "\n")
	}

	return nil
}

func (c *BuildCommand) runBatch() error {
	contents, err := c.readInput(c.InputFile)
	if err != nil {
		return err
	}

	var batch core.Batch
	if err := json.Unmarshal(contents, &batch); err != nil {
		return err
	}

	if len(batch.Items) == 0 {
		return errors.New("no items defined")
	}

	for _, item := range batch.Items {
		if err = c.runBatchItem(item, batch.Defaults); err != nil {
			return err
		}
	}

	return nil
}

func (c *BuildCommand) runBatchItem(item core.BatchItem, defaults core.BatchDefault) error {
	cfg := c.combineBatchItem(item, defaults)
	contents, err := c.readInput(cfg.Template)
	if err != nil {
		return err
	}

	writer, err := c.selectWriter(cfg.Target)
	if err != nil {
		return fmt.Errorf("select writer: %w", err)
	}

	if !c.NoCloseWriter {
		if oc, ok := writer.(io.Closer); ok {
			defer oc.Close()
		}
	}

	builder := parser.NewTemplate(cfg.Template, string(contents), cfg.Variables)
	if err = builder.Build(writer); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	return nil
}

func (c *BuildCommand) combineBatchItem(item core.BatchItem, defaults core.BatchDefault) core.BatchItem {
	if len(item.Info) == 0 {
		item.Info = defaults.Info
	}

	if len(item.Template) == 0 {
		item.Template = defaults.Template
	}

	if len(item.Variables) == 0 {
		item.Variables = defaults.Variables
	} else {
		for k, v := range defaults.Variables {
			if _, ok := item.Variables[k]; !ok {
				item.Variables[k] = v
			}
		}
	}

	return item
}
