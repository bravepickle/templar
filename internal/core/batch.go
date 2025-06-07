package core

// Params parser params list
type Params map[string]any

type BatchVariables map[string]any

type BatchItem struct {
	// Info description of the item
	Info string `json:"info,omitempty"`

	// InputFormat is a input format. Allowed: env, json
	InputFormat string `json:"format,omitempty"`

	// Input is a source file for input variables
	Input string `json:"input,omitempty"`

	// Variables is a list of variables to apply. Exclusive to Input
	Variables Params `json:"variables,omitempty"`

	// Template is a template file
	Template string `json:"template,omitempty"`

	// Output is a target file to write results to. Will overwrite contents
	Output string `json:"output,omitempty"`
}

type BatchDefault struct {
	// Info description of the item
	Info string `json:"info,omitempty"`

	// InputFormat is a input format
	InputFormat string `json:"format,omitempty"`

	// Input is a source file for input variables
	Input string `json:"input,omitempty"`

	// Variables is a list of variables to apply. Exclusive to Input
	Variables Params `json:"variables,omitempty"`

	// Template is a template file
	Template string `json:"template,omitempty"`
}

type Batch struct {
	// Items is a list of items to process
	Items []BatchItem `json:"items"`

	// Defaults defines some default values or expands Batch.Items if they are undefined.
	// Will be added to each item in the file.
	Defaults BatchDefault `json:"defaults"`
}
