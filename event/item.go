package event

type ItemReceivedCommand struct {
	topic       string
	CommandType string
	Command     string
}

// func NewItemReceivedCommand(topic, payload string) (ItemReceivedCommand, error) {
// 	data := api.EventCommand{}
// 	err := json.Unmarshal([]byte(payload), &data)
// 	if err != nil {
// 		return ItemReceivedCommand{}, err
// 	}
// 	return ItemReceivedCommand{
// 		topic:       topic,
// 		CommandType: data.Type,
// 		Command:     data.Value,
// 	}, nil
// }

func (i ItemReceivedCommand) Topic() string {
	return i.topic
}

func (i ItemReceivedCommand) Type() Type {
	return ItemCommand
}

type ItemReceivedState struct {
	topic     string
	StateType string
	State     string
}

// func NewItemReceivedState(topic, payload string) (ItemReceivedState, error) {
// 	data := api.EventState{}
// 	err := json.Unmarshal([]byte(payload), &data)
// 	if err != nil {
// 		return ItemReceivedState{}, err
// 	}
// 	return ItemReceivedState{
// 		topic:     topic,
// 		StateType: data.Type,
// 		State:     data.Value,
// 	}, nil
// }

func (i ItemReceivedState) Topic() string {
	return i.topic
}

func (i ItemReceivedState) Type() Type {
	return ItemState
}

type ItemChanged struct {
	topic        string
	StateType    string
	State        string
	OldStateType string
	OldState     string
}

// func NewItemChanged(topic, payload string) (ItemChanged, error) {
// 	data := api.EventStateChanged{}
// 	err := json.Unmarshal([]byte(payload), &data)
// 	if err != nil {
// 		return ItemChanged{}, err
// 	}
// 	return ItemChanged{
// 		topic:        topic,
// 		StateType:    data.Type,
// 		State:        data.Value,
// 		OldStateType: data.OldType,
// 		OldState:     data.OldValue,
// 	}, nil
// }

func (i ItemChanged) Topic() string {
	return i.topic
}

func (i ItemChanged) Type() Type {
	return ItemStateChanged
}
