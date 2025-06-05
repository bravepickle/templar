package core

// Params parser params list
type Params map[string]any

type BatchVariables map[string]any

type BatchItem struct {
	Info      string `json:"info,omitempty"`
	Variables Params `json:"variables,omitempty"`
	Template  string `json:"template,omitempty"`
	Target    string `json:"target,omitempty"`
}

type BatchDefault struct {
	Info      string `json:"info,omitempty"`
	Variables Params `json:"variables,omitempty"`
	Template  string `json:"template,omitempty"`
}

type Batch struct {
	// Items is a list of items to process
	Items []BatchItem `json:"items"`

	// Defaults defines some default values or expands Batch.Items if they are undefined.
	// Will be added to each item in the file.
	Defaults BatchDefault `json:"defaults"`
}
