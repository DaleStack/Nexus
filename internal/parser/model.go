package parser

// Module represents a parsed Nexus module
type Module struct {
	Name     string     `json:"name"`
	State    []Property `json:"state"`
	Views    []View     `json:"views"`
	Actions  []Action   `json:"actions"`
	Template []string   `json:"template"`
}

// Property represents a state variable
type Property struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// View represents a view declaration
type View struct {
	Name    string   `json:"name"`
	Content []string `json:"content"`
}

// Action represents an action that can be performed
type Action struct {
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters"`
}

// Parameter represents a parameter for an action
type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// StateProperty is an alias for Property for backward compatibility
type StateProperty = Property
