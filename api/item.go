package api

type Item struct {
	Name               string              `json:"name"`
	Label              string              `json:"label"`
	Link               string              `json:"link"`
	Type               string              `json:"type"`
	State              string              `json:"state"`
	TransformedState   string              `json:"transformedState,omitempty"`
	Editable           bool                `json:"editable"`
	Category           string              `json:"category"`
	Tags               []string            `json:"tags"`
	GroupNames         []string            `json:"groupNames"`
	Members            []string            `json:"members,omitempty"`
	GroupType          string              `json:"groupType,omitempty"`
	Function           *Function           `json:"function,omitempty"`
	StateDescription   *StateDescription   `json:"stateDescription,omitempty"`
	CommandDescription *CommandDescription `json:"commandDescription,omitempty"`
}

type Function struct {
	Name   string   `json:"name"`
	Params []string `json:"params"`
}

type StateDescription struct {
	Minimum  int            `json:"minimum"`
	Maximum  int            `json:"maximum"`
	Step     int            `json:"step"`
	Pattern  string         `json:"pattern"`
	ReadOnly bool           `json:"readOnly"`
	Options  []StateOptions `json:"options"`
}

type StateOptions struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type CommandDescription struct {
	Options []CommandOptions `json:"commandOptions"`
}

type CommandOptions struct {
	Command string `json:"command"`
	Label   string `json:"label"`
}
