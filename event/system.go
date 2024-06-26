package event

import "fmt"

// SystemEvent is used for system events not generated by openHAB
type SystemEvent struct {
	eventType Type
}

// NewSystemEvent creates ClientStart, ClientConnected, ClientConnectionStable, ClientDisconnected, ClientStop and TimeCron event types
func NewSystemEvent(eventType Type) SystemEvent {
	return SystemEvent{
		eventType: eventType,
	}
}

// Topic is always empty on SystemEvent
func (e SystemEvent) Topic() string {
	return ""
}

// Type is either ClientStart, ClientConnected, ClientConnectionStable, ClientDisconnected, ClientStop or TimeCron
func (e SystemEvent) Type() Type {
	return e.eventType
}

func (e SystemEvent) String() string {
	return fmt.Sprintf("System event #%d", e.eventType)
}

var _ Event = &SystemEvent{}
