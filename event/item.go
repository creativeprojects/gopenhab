package event

import (
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

// Verify interface
var _ Event = ItemReceivedCommand{}

// ItemReceivedState is sent when the state of an item is about to get updated.
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

// Verify interface
var _ Event = ItemReceivedState{}

// ItemStateUpdated is sent when the state of an item has been updated.
type ItemStateUpdated struct {
	topic     string
	ItemName  string
	StateType string
	State     string
}

func NewItemStateUpdated(itemName, stateType, state string) ItemStateUpdated {
	topic := itemTopicPrefix + itemName + "/" + api.TopicEventStateUpdated
	return ItemStateUpdated{
		topic:     topic,
		ItemName:  itemName,
		StateType: stateType,
		State:     state,
	}
}

func (i ItemStateUpdated) Topic() string {
	return i.topic
}

func (i ItemStateUpdated) Type() Type {
	return TypeItemState // should we create a new type for this?
}

// Verify interface
var _ Event = ItemStateUpdated{}

// ItemStateChanged is sent when the state of an item has changed.
type ItemStateChanged struct {
	topic             string
	ItemName          string
	NewStateType      string
	NewState          string
	PreviousStateType string
	PreviousState     string
}

func NewItemStateChanged(itemName, previousStateType, previousState, newStateType, newState string) ItemStateChanged {
	topic := itemTopicPrefix + itemName + "/" + api.TopicEventStateChanged
	return ItemStateChanged{
		topic:             topic,
		ItemName:          itemName,
		PreviousStateType: previousStateType,
		PreviousState:     previousState,
		NewStateType:      newStateType,
		NewState:          newState,
	}
}

func (i ItemStateChanged) Topic() string {
	return i.topic
}

func (i ItemStateChanged) Type() Type {
	return TypeItemStateChanged
}

// Verify interface
var _ Event = ItemStateChanged{}

type ItemStatePredicted struct {
	topic          string
	ItemName       string
	PredictedType  string
	PredictedState string
}

func NewItemStatePredicted(itemName, predictedStateType, predictedState string) ItemStatePredicted {
	topic := itemTopicPrefix + itemName + "/" + api.TopicEventStatePredicted
	return ItemStatePredicted{
		topic:          topic,
		ItemName:       itemName,
		PredictedType:  predictedStateType,
		PredictedState: predictedState,
	}
}

func (i ItemStatePredicted) Topic() string {
	return i.topic
}

func (i ItemStatePredicted) Type() Type {
	return TypeItemStatePredicted
}

// Verify interface
var _ Event = ItemStatePredicted{}

type ItemAdded struct {
	topic string
	Item  Item
}

func NewItemAdded(item Item) ItemAdded {
	topic := itemTopicPrefix + item.Name + "/" + api.TopicEventAdded
	return ItemAdded{
		topic: topic,
		Item:  item,
	}
}

func (i ItemAdded) Topic() string {
	return i.topic
}

func (i ItemAdded) Type() Type {
	return TypeItemAdded
}

// Verify interface
var _ Event = ItemAdded{}

type ItemRemoved struct {
	topic string
	Item  Item
}

func NewItemRemoved(item Item) ItemRemoved {
	topic := itemTopicPrefix + item.Name + "/" + api.TopicEventRemoved
	return ItemRemoved{
		topic: topic,
		Item:  item,
	}
}

func (i ItemRemoved) Topic() string {
	return i.topic
}

func (i ItemRemoved) Type() Type {
	return TypeItemRemoved
}

// Verify interface
var _ Event = ItemRemoved{}

type ItemUpdated struct {
	topic   string
	OldItem Item
	Item    Item
}

func NewItemUpdated(oldItem, newItem Item) ItemUpdated {
	topic := itemTopicPrefix + newItem.Name + "/" + api.TopicEventUpdated
	return ItemUpdated{
		topic:   topic,
		OldItem: oldItem,
		Item:    newItem,
	}
}

func (i ItemUpdated) Topic() string {
	return i.topic
}

func (i ItemUpdated) Type() Type {
	return TypeItemUpdated
}

// Verify interface
var _ Event = ItemUpdated{}
