package command

// MkDirPerm defines default permissions for created directories
const MkDirPerm = 0755

// MkFilePerm defines default permissions for created files
const MkFilePerm = 0644

const ExampleEnv = `
# This is an example environment variable configuration file.
# Change it to fit your needs. Comments and quotes are supported.
FOO=bar
NUMBER=55
GREET="Hello, world!"
BOOL=true
`

const ExampleJson = `{
  "debug": true,
  "foo": "bar",
  "NUMBER": 42,
  "nested": {"baz": "faz"}
}
`

const ExampleBatchJson = `{
  "items": [
    {
      "info": "This template will combine defaults with other values.",
      "template": "rendered.txt",
      "source": "./templates/custom.tpl",
      "variables": {
        "FOO": "custom",
        "extra": "extra value"
      }
    },
    {
	  "info": "Will use all defaults to render the template.",
      "template": "with_defaults.txt"
    },
    {
	  "info": "Without defaults to render the template.",
      "template": "without_defaults.txt",
      "no_defaults": true
    }
  ],
  "defaults": {
    "info": "This section defines default values to be used in the \"items\" section.",
    "source": "./templates/template.txt.tpl",
    "variables": {
      "FOO": "bar",
      "NUMBER": 42,
      "nested": {"baz": "faz"}
    }
  }
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
