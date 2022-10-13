package event

type fakeEvent struct {
	topic     string
	eventType Type
}

// newFakeEvent creates a fake event: please start the topic at "items/"
func newFakeEvent(topic string, eventType Type) fakeEvent {
	return fakeEvent{
		topic:     topic,
		eventType: eventType,
	}
}

func (e fakeEvent) Topic() string {
	return e.topic
}

func (e fakeEvent) Type() Type {
	return e.eventType
}

var _ Event = &fakeEvent{}
