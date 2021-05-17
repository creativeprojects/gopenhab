package event

import (
	"strings"

	"github.com/creativeprojects/gopenhab/api"
)

type Item struct {
	Name       string
	Label      string
	Type       string
	Category   string
	Tags       []string
	GroupNames []string
	Members    []string
	GroupType  string
}

type ItemReceivedCommand struct {
	topic       string
	ItemName    string
	CommandType string
	Command     string
}

func NewItemReceivedCommand(itemName, commandType, command string) ItemReceivedCommand {
	topic := itemTopicPrefix + itemName + "/" + api.TopicEventCommand
	return ItemReceivedCommand{
		topic:       topic,
		ItemName:    itemName,
		CommandType: commandType,
		Command:     command,
	}
}

func (i ItemReceivedCommand) Topic() string {
	return i.topic
}

func (i ItemReceivedCommand) Type() Type {
	return TypeItemCommand
}

type ItemReceivedState struct {
	topic     string
	ItemName  string
	StateType string
	State     string
}

func NewItemReceivedState(itemName, stateType, state string) ItemReceivedState {
	topic := itemTopicPrefix + itemName + "/" + api.TopicEventState
	return ItemReceivedState{
		topic:     topic,
		ItemName:  itemName,
		StateType: stateType,
		State:     state,
	}
}

func (i ItemReceivedState) Topic() string {
	return i.topic
}

func (i ItemReceivedState) Type() Type {
	return TypeItemState
}

type ItemStateChanged struct {
	topic             string
	ItemName          string
	NewStateType      string
	NewState          string
	PreviousStateType string
	PreviousState     string
}

func NewItemStateChanged(itemName, stateType, previousState, newState string) ItemStateChanged {
	topic := itemTopicPrefix + itemName + "/" + api.TopicEventStateChanged
	return ItemStateChanged{
		topic:             topic,
		ItemName:          itemName,
		PreviousStateType: stateType,
		PreviousState:     previousState,
		NewStateType:      stateType,
		NewState:          newState,
	}
}

func (i ItemStateChanged) Topic() string {
	return i.topic
}

func (i ItemStateChanged) Type() Type {
	return TypeItemStateChanged
}

type GroupItemStateChanged struct {
	topic             string
	ItemName          string
	TriggeringItem    string
	NewStateType      string
	NewState          string
	PreviousStateType string
	PreviousState     string
}

func (i GroupItemStateChanged) Topic() string {
	return i.topic
}

func (i GroupItemStateChanged) Type() Type {
	return TypeGroupItemStateChanged
}

type ItemAdded struct {
	topic string
	Item  Item
}

func (i ItemAdded) Topic() string {
	return i.topic
}

func (i ItemAdded) Type() Type {
	return TypeItemAdded
}

type ItemRemoved struct {
	topic string
	Item  Item
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
	Item    Item
}

func (i ItemUpdated) Topic() string {
	return i.topic
}

func (i ItemUpdated) Type() Type {
	return TypeItemUpdated
}

// splitItemTopic returns the item name, triggering item (if any) and the event type
func splitItemTopic(topic string) (string, string, string) {
	parts := strings.Split(topic, "/")
	if len(parts) < 4 || len(parts) > 5 || parts[0] != "smarthome" || parts[1] != "items" {
		return "", "", ""
	}
	if len(parts) == 5 {
		return parts[2], parts[3], parts[4]
	}
	return parts[2], "", parts[3]
}
