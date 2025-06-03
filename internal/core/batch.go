package core

type BatchVariables map[string]any

type BatchItem struct {
	Info      string         `json:"info,omitempty"`
	Template  string         `json:"template,omitempty"`
	Variables BatchVariables `json:"variables,omitempty"`
}

type BatchDefault struct {
	Info      string         `json:"info,omitempty"`
	Template  string         `json:"template,omitempty"`
	Variables BatchVariables `json:"variables,omitempty"`
}

type Batch struct {
	// Items is a list of items to process
	Items []BatchItem `json:"items"`

	// Defaults defines some default values or expands Batch.Items if they are undefined.
	// Will be added to each item in the file.
	Defaults BatchDefault `json:"defaults"`
}
