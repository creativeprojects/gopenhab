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
			return nil, errDecodingMessage(err)
		}
		itemName, _, _ := splitItemTopic(message.Topic)
		if itemName == "" {
			return nil, errInvalidTopic(message.Topic)
		}
		return NewItemReceivedCommand(itemName, data.Type, data.Value), nil

	case api.EventItemState:
		data := api.EventState{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, errDecodingMessage(err)
		}
		itemName, _, _ := splitItemTopic(message.Topic)
		if itemName == "" {
			return nil, errInvalidTopic(message.Topic)
		}
		return NewItemReceivedState(itemName, data.Type, data.Value), nil

	case api.EventItemStateChanged:
		data := api.EventStateChanged{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, errDecodingMessage(err)
		}
		itemName, _, _ := splitItemTopic(message.Topic)
		if itemName == "" {
			return nil, errInvalidTopic(message.Topic)
		}
		return NewItemStateChanged(itemName, data.OldType, data.OldValue, data.Type, data.Value), nil

	case api.EventGroupItemStateChanged:
		data := api.EventStateChanged{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, errDecodingMessage(err)
		}
		itemName, triggeringItem, _ := splitItemTopic(message.Topic)
		if itemName == "" {
			return nil, errInvalidTopic(message.Topic)
		}
		return NewGroupItemStateChanged(itemName, triggeringItem, data.OldType, data.OldValue, data.Type, data.Value), nil

	case api.EventItemAdded:
		data := api.Item{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, errDecodingMessage(err)
		}
		return NewItemAdded(Item{
			Type:       data.Type,
			GroupType:  data.GroupType,
			Name:       data.Name,
			Label:      data.Label,
			Category:   data.Category,
			Tags:       data.Tags,
			GroupNames: data.GroupNames,
		},
		), nil

	case api.EventItemRemoved:
		data := api.Item{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, errDecodingMessage(err)
		}
		return NewItemRemoved(Item{
			Type:       data.Type,
			GroupType:  data.GroupType,
			Name:       data.Name,
			Label:      data.Label,
			Category:   data.Category,
			Tags:       data.Tags,
			GroupNames: data.GroupNames,
		},
		), nil

	case api.EventItemUpdated:
		data := make([]api.Item, 2)
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, errDecodingMessage(err)
		}
		if len(data) != 2 {
			return nil, fmt.Errorf("error decoding message: expected array with 2 elements, but found %d", len(data))
		}
		return NewItemUpdated(
			Item{
				Type:       data[1].Type,
				GroupType:  data[1].GroupType,
				Name:       data[1].Name,
				Label:      data[1].Label,
				Category:   data[1].Category,
				Tags:       data[1].Tags,
				GroupNames: data[1].GroupNames,
			},
			Item{
				Type:       data[0].Type,
				GroupType:  data[0].GroupType,
				Name:       data[0].Name,
				Label:      data[0].Label,
				Category:   data[0].Category,
				Tags:       data[0].Tags,
				GroupNames: data[0].GroupNames,
			},
		), nil

	case api.EventThingStatusInfo:
		data := api.ThingStatusInfo{}
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, errDecodingMessage(err)
		}
		thingName, _ := splitThingTopic(message.Topic)
		return NewThingStatusInfoEvent(thingName, ThingStatus{
			Status:       data.Status,
			StatusDetail: data.StatusDetail,
		}), nil

	case api.EventThingStatusInfoChanged:
		data := make([]api.ThingStatusInfo, 0, 2)
		err := json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			return nil, errDecodingMessage(err)
		}
		if len(data) != 2 {
			return nil, fmt.Errorf("error decoding message: expected array with 2 elements, but found %d", len(data))
		}
		thingName, _ := splitThingTopic(message.Topic)
		return NewThingStatusInfoChangedEvent(thingName,
			ThingStatus{
				Status:       data[1].Status,
				StatusDetail: data[1].StatusDetail,
				Description:  data[1].Description,
			},
			ThingStatus{
				Status:       data[0].Status,
				StatusDetail: data[0].StatusDetail,
				Description:  data[0].Description,
			}), nil

	case api.EventTypeAlive:
		return NewAliveEvent(), nil

	default:
		return NewGenericEvent(message.Type, message.Topic, message.Payload), nil
	}
}

func errInvalidTopic(topic string) error {
	return fmt.Errorf("invalid topic: %q", topic)
}

func errDecodingMessage(err error) error {
	return fmt.Errorf("error decoding message: %w", err)
}
