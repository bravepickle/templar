package command

// MkDirPerm defines default permissions for created directories
const MkDirPerm = 0755

const ExampleEnv = `
# This is an example environment variable configuration file.
# Change it to fit your needs. Comments and quotes are supported.
FOO=bar
NUMBER=55
GREET="Hello, world!"
BOOL=true
`

const ExampleJson = `
{
	"debug": true,
	"foo": "bar",
	"NUMBER": 42,
	"nested": {"baz": "faz"}
}
`

const ExampleTemplate = `
This is an example of template for templar to use
{{ .GREET }} and {{ .BOOL }}.
{{ if .debug -}}
Debug is enabled
{{ end }}
My lucky number is {{ .NUMBER }}.

JSON has nested parameter {{ default "N/A" .nested }}.
`
