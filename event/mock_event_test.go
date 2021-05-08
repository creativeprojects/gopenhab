package event

type mockEvent struct {
	topic     string
	eventType Type
}

func newMockEvent(topic string, eventType Type) mockEvent {
	return mockEvent{
		topic:     topic,
		eventType: eventType,
	}
}

func (e mockEvent) Topic() string {
	return e.topic
}

func (e mockEvent) Type() Type {
	return e.eventType
}

var _ Event = &mockEvent{}
