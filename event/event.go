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
		return nil, fmt.Errorf("invalid event data %q: %w", data, err)
	}
	switch message.Type {
	case api.EventItemCommand:
		data := api.EventCommand{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, fmt.Errorf("error decoding message: %w", err)
		}
		itemName, _, _ := splitItemTopic(message.Topic)
		if itemName == "" {
			return nil, fmt.Errorf("invalid topic: %q", message.Topic)
		}
		return ItemReceivedCommand{
			topic:       message.Topic,
			ItemName:    itemName,
			CommandType: data.Type,
			Command:     data.Value,
		}, nil

	case api.EventItemState:
		data := api.EventState{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, fmt.Errorf("error decoding message: %w", err)
		}
		itemName, _, _ := splitItemTopic(message.Topic)
		if itemName == "" {
			return nil, fmt.Errorf("invalid topic: %q", message.Topic)
		}
		return ItemReceivedState{
			topic:     message.Topic,
			ItemName:  itemName,
			StateType: data.Type,
			State:     data.Value,
		}, nil

	case api.EventItemStateChanged:
		data := api.EventStateChanged{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, fmt.Errorf("error decoding message: %w", err)
		}
		itemName, _, _ := splitItemTopic(message.Topic)
		if itemName == "" {
			return nil, fmt.Errorf("invalid topic: %q", message.Topic)
		}
		return ItemStateChanged{
			topic:             message.Topic,
			ItemName:          itemName,
			NewStateType:      data.Type,
			NewState:          data.Value,
			PreviousStateType: data.OldType,
			PreviousState:     data.OldValue,
		}, nil

	case api.EventGroupItemStateChanged:
		data := api.EventStateChanged{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, fmt.Errorf("error decoding message: %w", err)
		}
		itemName, triggeringItem, _ := splitItemTopic(message.Topic)
		if itemName == "" {
			return nil, fmt.Errorf("invalid topic: %q", message.Topic)
		}
		return GroupItemStateChanged{
			topic:             message.Topic,
			ItemName:          itemName,
			TriggeringItem:    triggeringItem,
			NewStateType:      data.Type,
			NewState:          data.Value,
			PreviousStateType: data.OldType,
			PreviousState:     data.OldValue,
		}, nil

	case api.EventItemAdded:
		data := api.Item{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, fmt.Errorf("error decoding message: %w", err)
		}
		return ItemAdded{
			topic: message.Topic,
			Item: Item{
				Type:       data.Type,
				GroupType:  data.GroupType,
				Name:       data.Name,
				Label:      data.Label,
				Category:   data.Category,
				Tags:       data.Tags,
				GroupNames: data.GroupNames,
			},
		}, nil

	case api.EventItemRemoved:
		data := api.Item{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, fmt.Errorf("error decoding message: %w", err)
		}
		return ItemRemoved{
			topic: message.Topic,
			Item: Item{
				Type:       data.Type,
				GroupType:  data.GroupType,
				Name:       data.Name,
				Label:      data.Label,
				Category:   data.Category,
				Tags:       data.Tags,
				GroupNames: data.GroupNames,
			},
		}, nil

	case api.EventItemUpdated:
		data := make([]api.Item, 2)
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, fmt.Errorf("error decoding message: %w", err)
		}
		if len(data) != 2 {
			return nil, fmt.Errorf("error decoding message: expected array with 2 elements, but found %d", len(data))
		}
		return ItemUpdated{
			topic: message.Topic,
			Item: Item{
				Type:       data[0].Type,
				GroupType:  data[0].GroupType,
				Name:       data[0].Name,
				Label:      data[0].Label,
				Category:   data[0].Category,
				Tags:       data[0].Tags,
				GroupNames: data[0].GroupNames,
			},
			OldItem: Item{
				Type:       data[1].Type,
				GroupType:  data[1].GroupType,
				Name:       data[1].Name,
				Label:      data[1].Label,
				Category:   data[1].Category,
				Tags:       data[1].Tags,
				GroupNames: data[1].GroupNames,
			},
		}, nil

	default:
		return GenericEvent{
			typeName: message.Type,
			topic:    message.Topic,
			payload:  message.Payload,
		}, nil
	}
}
