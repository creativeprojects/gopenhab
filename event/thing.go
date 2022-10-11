package event

import (
	"github.com/creativeprojects/gopenhab/api"
)

type ThingStatusInfoEvent struct {
	topic        string
	ThingName    string
	Status       string
	StatusDetail string
}

// NewThingStatusInfoEvent create a ThingStatusInfoEvent.
// Please note the name of the thing is not present
// in the payload so we pass the topic as first argument.
func NewThingStatusInfoEvent(thingName, status, statusDetail string) ThingStatusInfoEvent {
	topic := thingTopicPrefix + thingName + "/" + api.TopicEventStatus

	return ThingStatusInfoEvent{
		topic:        topic,
		ThingName:    thingName,
		Status:       status,
		StatusDetail: statusDetail,
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
