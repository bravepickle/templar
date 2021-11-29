# templar
Small files generator from templates written in Golang

Press *Ctrl+D* to stop writing to STDIN interactively

See also GOLANG (template engine docs)[https://golang.org/pkg/html/template/]

## Additional template functions
- sub - subtract one value from another (int): `{{ sub .num 4 }}`
- sum - sum of two values (int): `{{ sum .num 1 }}`
- repeat - repeat n times: `{{ repeat 3 "Hurray! " }}`

See also (default functions and usages)[https://golang.org/pkg/text/template/#hdr-Functions]

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

## Usage
```
Usage: templar [OPTIONS] [COMMAND] [COMMAND_OPTIONS]
  -h	Print command usage options [Optional].
  -help
    	Print command usage options [Optional].
  -v	Run in verbose mode [Optional].
  -verbose
    	Run in verbose mode [Optional].

Commands:
    list    List all available commands.
    init    Initialize project for templated within current dir.
    build   Build file from template.

Examples:
    templar -h      See help for using this command
    templar init    Initialize current working directory as new project
    templar --verbose
	init Initialize new project in verbose mode
    templar build --format=env -d /tmp --format=env --input=./data.env --batch ./batch.json Build templates batch from file
    templar build --format=env --input=./data.env --template=./templates/test.tpl --output=./out.txt
	 Create out.txt file from test.tpl and environment parameters found in data.env file

```