package event

import "time"

// RulePanicEvent is used for errors generated by the client
type RulePanicEvent struct {
	message     string
	id          string
	name        string
	description string
	event       string
	timestamp   time.Time
}

// NewRulePanicEvent creates RulePanic event type
func NewRulePanicEvent(message, id, name, description, event string, timestamp time.Time) RulePanicEvent {
	return RulePanicEvent{
		message:     message,
		id:          id,
		name:        name,
		description: description,
		event:       event,
		timestamp:   timestamp,
	}
}

// Topic is always empty on RulePanicEvent
func (e RulePanicEvent) Topic() string {
	return ""
}

// Type is ClientError
func (e RulePanicEvent) Type() Type {
	return TypeRulePanic
}

func (e RulePanicEvent) String() string {
	return "Caught a panic from inside rule code: " + e.message
}

func (e RulePanicEvent) Message() string {
	return e.message
}

func (e RulePanicEvent) RuleID() string {
	return e.id
}

func (e RulePanicEvent) RuleName() string {
	return e.name
}

func (e RulePanicEvent) RuleDescription() string {
	return e.description
}

func (e RulePanicEvent) Event() string {
	return e.event
}

func (e RulePanicEvent) Timestamp() time.Time {
	return e.timestamp
}

var _ Event = &RulePanicEvent{}
