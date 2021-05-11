package event

type Item struct {
	Name             string
	Label            string
	Link             string
	Type             string
	State            string
	TransformedState string
	Editable         bool
	Category         string
	Tags             []string
	GroupNames       []string
	Members          []string
	GroupType        string
}

type ItemReceivedCommand struct {
	topic       string
	CommandType string
	Command     string
}

func (i ItemReceivedCommand) Topic() string {
	return i.topic
}

func (i ItemReceivedCommand) Type() Type {
	return TypeItemCommand
}

type ItemReceivedState struct {
	topic     string
	StateType string
	State     string
}

func (i ItemReceivedState) Topic() string {
	return i.topic
}

func (i ItemReceivedState) Type() Type {
	return TypeItemState
}

type ItemStateChanged struct {
	topic        string
	StateType    string
	State        string
	OldStateType string
	OldState     string
}

func (i ItemStateChanged) Topic() string {
	return i.topic
}

func (i ItemStateChanged) Type() Type {
	return TypeItemStateChanged
}

type ItemAdded struct {
	topic string
	Item
}

func (i ItemAdded) Topic() string {
	return i.topic
}

func (i ItemAdded) Type() Type {
	return TypeItemAdded
}

type ItemRemoved struct {
	topic string
	Item
}

func (i ItemRemoved) Topic() string {
	return i.topic
}

func (i ItemRemoved) Type() Type {
	return TypeItemRemoved
}

type ItemUpdated struct {
	topic   string
	OldItem Item
	Item
}

func (i ItemUpdated) Topic() string {
	return i.topic
}

func (i ItemUpdated) Type() Type {
	return TypeItemUpdated
}
