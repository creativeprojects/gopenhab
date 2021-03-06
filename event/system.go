package event

// SystemEvent is used for system events not generated by openHAB
type SystemEvent struct {
	eventType Type
	message   string
}

// NewSystemEvent creates ClientStart, ClientConnected and ClientDisconnected event types
func NewSystemEvent(eventType Type) SystemEvent {
	return SystemEvent{
		eventType: eventType,
	}
}

// Topic is always empty on SystemEvent
func (e SystemEvent) Topic() string {
	return e.message
}

// Type is either ClientStart, ClientConnected, ClientDisconnected or ClientStop
func (e SystemEvent) Type() Type {
	return e.eventType
}

var _ Event = &SystemEvent{}
