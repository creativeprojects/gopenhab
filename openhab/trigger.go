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

//go:generate mockery --name subscriber --inpackage
type subscriber interface {
	subscribe(name string, eventType event.Type, callback func(e event.Event)) int
}

type baseTrigger struct{}

func (t baseTrigger) subscribe(
	client subscriber,
	thing string,
	eventType event.Type,
	run func(ev event.Event),
	match func(e event.Event) bool,
) int {
	return client.subscribe(thing, eventType, func(e event.Event) {
		if run == nil {
			return
		}
		if match(e) {
			run(e)
		}
	})
}
