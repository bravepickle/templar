# templar
Simple CLI templates generator written in Go without any external dependencies.

Press *Ctrl+D* to stop writing to STDIN interactively

See also GOLANG (template engine docs)[https://golang.org/pkg/html/template/]

## Install
### Option 1
1. Install GO
2. Run `go install github.com/bravepickle/templar@latest`

### Option 2
1. Go to https://github.com/bravepickle/templar/releases
2. Download binary file according to your OS and architecture

### Option 3
1. Create your own Docker container (for example, using image https://hub.docker.com/_/golang)
2. Clone repository
3. Run `make release` or `make build` to build binary file

## Additional template functions
- for now uses https://masterminds.github.io/sprig/ for additional functions supported
- regular templating functions for Golang packages also supported. See official docs for `text/template`, `html/template` packages.

See also (default functions and usages)[https://golang.org/pkg/text/template/#hdr-Functions]

## Variables overrides order
### Batch
OS ENV -> item variables from batch file -> defaults in batch file

## TODO
- [x] Add binary executables for some of the architectures
- [ ] Read docs with installation and usage instructions
- [ ] Provide examples
- [ ] Support types: text, html
- [ ] Optionally use facter or similar to pass extra params from environment and similar sources
- [ ] As input use ENV values, JSON, key-values from file or directly set as params
- [x] Read data for template by piping
- [ ] Support all formats specified
- [ ] Specify in docs all available template functions
- [ ] Support configs that contain multiple templates to generate - some kind of templates aggregator to easily template files in batches
- [ ] Support both html and text formatters, e.g. by adding flag `templar build --html ...` to a command. By default, use text
- [ ] Add raw|unescape function to FuncMap for html text formatting. Should skip HTML escapes for `html/template` package. E.g. `return template.HTML(text)`.
- [ ] On general help list commands with basic description. See `docker compose help` as an example for formatting and texts
- [ ] On command help view full description. Make new functions for detailed and short description
- [ ] Add debug command to see all resulting variables available for the template
- [ ] If debug is disabled then do not show error messages. Instead, show help blocks with non-zero exit code
- [ ] Add SchemaJSON file for `batch.json` files.
- [ ] Add support of JSONL format. For STDIN and other to stream data. Instead of ReadAll, we should process each template per line on-fly and skip blank lines.
- [ ] Batch support input formats - env, json
- [ ] `--verbose`, `--debug` flags reconsider their application
- [ ] Read from stdin only if some flag passed. E.g. `--input -`

## Usage
### Main command
```
$ templar
Usage: templar [OPTIONS] COMMAND [COMMAND_ARGS]

templar    generate template contents with provided variables

Commands:
  init       init default files structure for building templates
  build      render template contents with provided variables
  help       show help information on command or subcommand usage. Type "templar help help" to see help command usage information
  version    show application information on its build version and directories

Options:
  -debug
        debug mode
  -no-color
        disable color and styles output
  -verbose
        verbose output
  -workdir string
        working directory path (default "/home/user/templar")
```        

### Build template command
```
$ templar --no-color help build
Usage: templar [OPTIONS] build [COMMAND_OPTIONS]

build      render template contents with provided variables

Options:
  -clear
        clear ENV variables before building variables to avoid collisions
  -dump string
        show all available variables for the template to use and stop processing. Pass optionally --verbose or --debug flags for more information. Allowed dump formats: env, json, json_compact
  -format string
        input file format for variables' file. Allowed: env, json, batch (default "env")
  -input string
        file path which contains variables for template to use or batch file. Format should match "-format" value
  -output string
        output file path, If empty, outputs to stdout. If "-batch" option is used, specifies output directory
  -skip
        skip generation if target files already exist
  -template string
        template file path, If empty and "-batch" not defined, reads from stdin

Examples:
  $ templar build --input .env --format env --template template.tpl --output output.txt 
      # generates output.txt file from the provided template.tpl and .env variables in env format (is the default one, can be ommitted)

  $ NAME=John templar build --template template.tpl --output output.txt
      # generates output.txt file from the provided template.tpl and provided env variable

  $ echo "My name is {{ .NAME }}" | NAME=John templar build
      # generates output.txt file from the provided template.tpl and provided env variable

  $ templar build --format json --input vars.json --dump env
      # dumps to stdout combined OS ENV and JSON variables. Used to check what variables are available

  $ templar --debug build --input vars.env --dump json --clear
      # dump variables in JSON format and display their values (--debug flag was added). OS ENV variables will be omitted
```