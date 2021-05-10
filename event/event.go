package event

import (
	"encoding/json"
	"fmt"

	"github.com/creativeprojects/gopenhab/api"
)

type Event interface {
	Topic() string
	Type() Type
}

func New(data string) (Event, error) {
	message := api.EventMessage{}
	err := json.Unmarshal([]byte(data), &message)
	if err != nil {
		return nil, fmt.Errorf("invalid event data: %s", err)
	}
	switch message.Type {
	case api.EventItemCommand:
		data := api.EventCommand{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, fmt.Errorf("error decoding message: %w", err)
		}
		return ItemReceivedCommand{
			topic:       message.Topic,
			CommandType: data.Type,
			Command:     data.Value,
		}, nil

	case api.EventItemState:
		data := api.EventState{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, fmt.Errorf("error decoding message: %w", err)
		}
		return ItemReceivedState{
			topic:     message.Topic,
			StateType: data.Type,
			State:     data.Value,
		}, nil

	case api.EventItemStateChanged:
		data := api.EventStateChanged{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, fmt.Errorf("error decoding message: %w", err)
		}
		return ItemChanged{
			topic:        message.Topic,
			StateType:    data.Type,
			State:        data.Value,
			OldStateType: data.OldType,
			OldState:     data.OldValue,
		}, nil

	default:
		return GenericEvent{
			typeName: message.Type,
			topic:    message.Topic,
			payload:  message.Payload,
		}, nil
	}
}
