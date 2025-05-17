package models

// FunctionDefinition represents an OpenAI function definition
type FunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// FunctionCall represents the details of a function called by OpenAI
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
