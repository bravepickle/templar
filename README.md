# templar
Simple CLI templates generator written in Go without any external dependencies.

Press *Ctrl+D* to stop writing to STDIN interactively

See also GOLANG (template engine docs)[https://golang.org/pkg/html/template/]

## Additional template functions
- for now uses https://masterminds.github.io/sprig/ for additional functions supported
- regular templating functions for Golang packages also supported. See official docs for `text/template`, `html/template` packages.

See also (default functions and usages)[https://golang.org/pkg/text/template/#hdr-Functions]

## Variables overrides order
### Batch
OS ENV -> item variables from batch file -> defaults in batch file

## TODO
- [ ] Add binary executables for some of the architectures
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

## Usage
```
Usage: templar [OPTIONS] COMMAND [COMMAND_ARGS]

Commands:
  help       show help information on command or subcommand usage. Type "templar help help" to see help command usage information
  version    show application information on its build version and directories
  init       init default files structure for building templates
  build      render template contents with provided variables

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