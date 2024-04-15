package event

import "strings"

// GenericEvent is used when its type is unknown
type GenericEvent struct {
	typeName string
	topic    string
	payload  string
}

func NewGenericEvent(eventType, topic, payload string) GenericEvent {
	topic = strings.TrimPrefix(topic, "smarthome/")
	topic = strings.TrimPrefix(topic, "openhab/")
	return GenericEvent{
		typeName: eventType,
		topic:    topic,
		payload:  payload,
	}
}

func (e GenericEvent) Topic() string {
	return e.topic
}

func (e GenericEvent) Type() Type {
	return TypeUnknown
}

func (e GenericEvent) Payload() string {
	return e.payload
}

func (e GenericEvent) TypeName() string {
	return e.typeName
}

func (e GenericEvent) String() string {
	return "Received unknown event " + e.typeName + " on topic " + e.topic + " with payload " + e.payload
}

var _ Event = &GenericEvent{}
