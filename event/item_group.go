package event

import "github.com/creativeprojects/gopenhab/api"

// GroupItemStateUpdated is sent when the state of a group of items has been updated.
type GroupItemStateUpdated struct {
	topic          string
	ItemName       string
	TriggeringItem string
	StateType      string
	State          string
}

func NewGroupItemStateUpdated(itemName, triggeringItem, stateType, state string) GroupItemStateUpdated {
	topic := itemTopicPrefix + itemName + "/" + triggeringItem + "/" + api.TopicEventStateUpdated
	return GroupItemStateUpdated{
		topic:          topic,
		ItemName:       itemName,
		TriggeringItem: triggeringItem,
		StateType:      stateType,
		State:          state,
	}
}

func (i GroupItemStateUpdated) Topic() string {
	return i.topic
}

func (i GroupItemStateUpdated) Type() Type {
	return TypeItemState // should we create a new type for this?
}

// Verify interface
var _ Event = GroupItemStateUpdated{}

type GroupItemStateChanged struct {
	topic             string
	ItemName          string
	TriggeringItem    string
	NewStateType      string
	NewState          string
	PreviousStateType string
	PreviousState     string
}

func NewGroupItemStateChanged(itemName, triggeringItem, previousStateType, previousState, newStateType, newState string) GroupItemStateChanged {
	topic := itemTopicPrefix + itemName + "/" + triggeringItem + "/" + api.TopicEventStateChanged
	return GroupItemStateChanged{
		topic:             topic,
		ItemName:          itemName,
		TriggeringItem:    triggeringItem,
		PreviousStateType: previousStateType,
		PreviousState:     previousState,
		NewStateType:      newStateType,
		NewState:          newState,
	}
}

func (i GroupItemStateChanged) Topic() string {
	return i.topic
}

func (i GroupItemStateChanged) Type() Type {
	return TypeGroupItemStateChanged
}

// Verify interface
var _ Event = GroupItemStateChanged{}
