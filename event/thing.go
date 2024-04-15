package event

import (
	"github.com/creativeprojects/gopenhab/api"
)

type Thing struct {
	UID           string
	Label         string
	BridgeUID     string
	Configuration map[string]any
	Properties    map[string]string
	ThingTypeUID  string
}

type ThingStatus struct {
	Status       string
	StatusDetail string
	Description  string
}

type ThingStatusInfoEvent struct {
	topic        string
	ThingName    string
	Status       string
	StatusDetail string
}

// NewThingStatusInfoEvent create a ThingStatusInfoEvent.
func NewThingStatusInfoEvent(thingName string, status ThingStatus) ThingStatusInfoEvent {
	topic := thingTopicPrefix + thingName + "/" + api.TopicEventStatus

	return ThingStatusInfoEvent{
		topic:        topic,
		ThingName:    thingName,
		Status:       status.Status,
		StatusDetail: status.StatusDetail,
	}
}

func (i ThingStatusInfoEvent) Topic() string {
	return i.topic
}

func (i ThingStatusInfoEvent) Type() Type {
	return TypeItemCommand
}

func (i ThingStatusInfoEvent) String() string {
	return "Thing " + i.ThingName + " status is " + i.Status
}

// Verify interface
var _ Event = ThingStatusInfoEvent{}

type ThingStatusInfoChangedEvent struct {
	topic                string
	ThingName            string
	PreviousStatus       string
	PreviousStatusDetail string
	PreviousDescription  string
	NewStatus            string
	NewStatusDetail      string
	NewDescription       string
}

// NewThingStatusInfoChangedEvent create a ThingStatusInfoChangedEvent.
func NewThingStatusInfoChangedEvent(thingName string, previousStatus, newStatus ThingStatus) ThingStatusInfoChangedEvent {
	topic := thingTopicPrefix + thingName + "/" + api.TopicEventStatusChanged

	return ThingStatusInfoChangedEvent{
		topic:                topic,
		ThingName:            thingName,
		PreviousStatus:       previousStatus.Status,
		PreviousStatusDetail: previousStatus.StatusDetail,
		PreviousDescription:  previousStatus.Description,
		NewStatus:            newStatus.Status,
		NewStatusDetail:      newStatus.StatusDetail,
		NewDescription:       newStatus.Description,
	}
}

func (i ThingStatusInfoChangedEvent) Topic() string {
	return i.topic
}

func (i ThingStatusInfoChangedEvent) Type() Type {
	return TypeItemCommand
}

func (i ThingStatusInfoChangedEvent) String() string {
	return "Thing " + i.ThingName + " status changed from " + i.PreviousStatus + " to " + i.NewStatus
}

// Verify interface
var _ Event = ThingStatusInfoChangedEvent{}

type ThingUpdated struct {
	topic    string
	OldThing Thing
	Thing    Thing
}

func NewThingUpdated(oldThing, newThing Thing) ThingUpdated {
	topic := thingTopicPrefix + newThing.UID + "/" + api.TopicEventUpdated
	return ThingUpdated{
		topic:    topic,
		OldThing: oldThing,
		Thing:    newThing,
	}
}

func (i ThingUpdated) Topic() string {
	return i.topic
}

func (i ThingUpdated) Type() Type {
	return TypeThingUpdated
}

func (i ThingUpdated) String() string {
	return "Thing " + i.Thing.UID + " updated"
}

// Verify interface
var _ Event = ThingUpdated{}
