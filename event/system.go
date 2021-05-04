package event

type SystemEvent struct {
	eventType Type
}

func NewSystemEvent(eventType Type) *SystemEvent {
	return &SystemEvent{
		eventType: eventType,
	}
}

func (e SystemEvent) Topic() string {
	return ""
}

func (e SystemEvent) Type() Type {
	return e.eventType
}

var _ Event = &SystemEvent{}
