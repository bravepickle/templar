# templar
Small files generator from templates written in Golang

Press *Ctrl+D* to stop writing to STDIN interactively

See also GOLANG (template engine docs)[https://golang.org/pkg/text/template/]

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
