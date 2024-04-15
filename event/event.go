package event

import (
	"encoding/json"
	"fmt"

	"github.com/creativeprojects/gopenhab/api"
)

type Event interface {
	Topic() string
	Type() Type
	String() string
}

func New(data string) (Event, error) {
	message := api.EventMessage{}
	err := json.Unmarshal([]byte(data), &message)
	if err != nil {
		return nil, fmt.Errorf("invalid event data %q: %w", data, err)
	}
	switch message.Type {
	case api.EventItemCommand:
		return newEventItemCommand(message)

	case api.EventItemState:
		return newEventItemState(message)

	case api.EventItemStateUpdated:
		return newEventItemStateUpdated(message)

	case api.EventItemStateChanged:
		return newEventItemStateChanged(message)

	case api.EventGroupItemStateUpdated:
		return newEventGroupItemStateUpdated(message)

	case api.EventGroupItemStateChanged:
		return newEventGroupItemStateChanged(message)

	case api.EventItemStatePredicted:
		return newEventItemStatePredicted(message)

	case api.EventItemAdded:
		return newEventItemAdded(message)

	case api.EventItemRemoved:
		return newEventItemRemoved(message)

	case api.EventItemUpdated:
		return newEventItemUpdated(message)

	case api.EventThingUpdated:
		return newEventThingUpdated(message)

	case api.EventThingStatusInfo:
		return newEventThingStatusInfo(message)

	case api.EventThingStatusInfoChanged:
		return newEventThingStatusInfoChanged(message)

	case api.EventTypeAlive:
		return NewAliveEvent(), nil

	case api.EventTypeStartlevel:
		return newEventTypeStartlevel(message)

	case api.EventChannelTriggered:
		return newEventChannelTriggered(message)

	default:
		return NewGenericEvent(message.Type, message.Topic, message.Payload), nil
	}
}

func newEventItemCommand(message api.EventMessage) (Event, error) {
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
}

func newEventItemState(message api.EventMessage) (Event, error) {
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
}

func newEventItemStateUpdated(message api.EventMessage) (Event, error) {
	data := api.EventState{}
	err := json.Unmarshal([]byte(message.Payload), &data)
	if err != nil {
		return nil, errDecodingMessage(err)
	}
	itemName, _, _ := splitItemTopic(message.Topic)
	if itemName == "" {
		return nil, errInvalidTopic(message.Topic)
	}
	return NewItemStateUpdated(itemName, data.Type, data.Value), nil
}

func newEventItemStateChanged(message api.EventMessage) (Event, error) {
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
}

func newEventGroupItemStateUpdated(message api.EventMessage) (Event, error) {
	data := api.EventState{}
	err := json.Unmarshal([]byte(message.Payload), &data)
	if err != nil {
		return nil, errDecodingMessage(err)
	}
	itemName, triggeringItem, _ := splitItemTopic(message.Topic)
	if itemName == "" {
		return nil, errInvalidTopic(message.Topic)
	}
	return NewGroupItemStateUpdated(itemName, triggeringItem, data.Type, data.Value), nil
}

func newEventGroupItemStateChanged(message api.EventMessage) (Event, error) {
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
}

func newEventItemStatePredicted(message api.EventMessage) (Event, error) {
	data := api.EventStatePredicted{}
	err := json.Unmarshal([]byte(message.Payload), &data)
	if err != nil {
		return nil, errDecodingMessage(err)
	}
	itemName, _, _ := splitItemTopic(message.Topic)
	if itemName == "" {
		return nil, errInvalidTopic(message.Topic)
	}
	return NewItemStatePredicted(itemName, data.PredictedType, data.PredictedValue), nil
}

func newEventItemAdded(message api.EventMessage) (Event, error) {
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
}

func newEventItemRemoved(message api.EventMessage) (Event, error) {
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
}

func newEventItemUpdated(message api.EventMessage) (Event, error) {
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
}

func newEventThingUpdated(message api.EventMessage) (Event, error) {
	data := make([]api.Thing, 2)
	err := json.Unmarshal([]byte(message.Payload), &data)
	if err != nil {
		return nil, errDecodingMessage(err)
	}
	if len(data) != 2 {
		return nil, fmt.Errorf("error decoding message: expected array with 2 elements, but found %d", len(data))
	}
	return NewThingUpdated(
		Thing{
			UID:           data[1].UID,
			Label:         data[1].Label,
			BridgeUID:     data[1].BridgeUID,
			Configuration: data[1].Configuration,
			Properties:    data[1].Properties,
			ThingTypeUID:  data[1].ThingTypeUID,
		},
		Thing{
			UID:           data[0].UID,
			Label:         data[0].Label,
			BridgeUID:     data[0].BridgeUID,
			Configuration: data[0].Configuration,
			Properties:    data[0].Properties,
			ThingTypeUID:  data[0].ThingTypeUID,
		},
	), nil
}

func newEventThingStatusInfo(message api.EventMessage) (Event, error) {
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
}

func newEventThingStatusInfoChanged(message api.EventMessage) (Event, error) {
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
}

func newEventChannelTriggered(message api.EventMessage) (Event, error) {
	data := api.EventTriggered{}
	err := json.Unmarshal([]byte(message.Payload), &data)
	if err != nil {
		return nil, errDecodingMessage(err)
	}
	channelName, _ := splitChannelTopic(message.Topic)
	return NewChannelTriggered(channelName, data.Event), nil
}

func newEventTypeStartlevel(message api.EventMessage) (Event, error) {
	data := api.Startlevel{}
	err := json.Unmarshal([]byte(message.Payload), &data)
	if err != nil {
		return nil, errDecodingMessage(err)
	}
	return NewStartlevelEvent(message.Topic, data.Startlevel), nil
}

func errInvalidTopic(topic string) error {
	return fmt.Errorf("invalid topic: %q", topic)
}

func errDecodingMessage(err error) error {
	return fmt.Errorf("error decoding message: %w", err)
}
