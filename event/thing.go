package event

import (
	"github.com/creativeprojects/gopenhab/api"
)

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

// Verify interface
var _ Event = ThingStatusInfoChangedEvent{}
