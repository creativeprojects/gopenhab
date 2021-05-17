package openhab

import "github.com/creativeprojects/gopenhab/event"

// Trigger is a generic interface for catching incoming messages on the event bus
type Trigger interface {
	// activate the trigger for func() in the context of a *Client
	activate(client *Client, run func(ev event.Event), ruleData RuleData) error
	// deactivate the trigger in the context of a *Client
	deactivate(client *Client)
	// match the event should activate the trigger.
	// This method should NOT be used outside of unit tests
	match(e event.Event) bool
}
