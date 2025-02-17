package openhab

import (
	"github.com/creativeprojects/gopenhab/event"
)

type thingReceivedStatusInfoTrigger struct {
	baseTrigger
	thing  string
	status ThingStatus
	subID  int
}

// OnThingReceivedStatusInfo triggers the rule when the thing received a status update equal to status.
// pass ThingStatusAny or "" to status to receive ANY update of the status
// This is an equivalent of the DSL rule:
//
// Thing <thingUID> received update [<status>]
// Thing <thingUID> changed [from <status>] [to <status>]
func OnThingReceivedStatusInfo(thing string, status ThingStatus) *thingReceivedStatusInfoTrigger {
	return &thingReceivedStatusInfoTrigger{
		thing:  thing,
		status: status,
	}
}

func (c *thingReceivedStatusInfoTrigger) activate(client subscriber, run func(ev event.Event), ruleData RuleData) error {
	if c.subID > 0 {
		return ErrRuleAlreadyActivated
	}
	c.subID = c.subscribe(client, c.thing, event.TypeThingStatusInfo, run, c.match)
	return nil
}

func (c *thingReceivedStatusInfoTrigger) deactivate(client subscriber) {
	if c.subID > 0 {
		client.unsubscribe(c.subID)
		c.subID = 0
	}
}

func (c *thingReceivedStatusInfoTrigger) match(e event.Event) bool {
	if c.status != "" {
		// check for the desired state
		if ev, ok := e.(event.ThingStatusInfoEvent); ok {
			if ev.Status != string(c.status) {
				// not the value we wanted
				return false
			}
		} else {
			panic("expected event of type event.ThingStatusInfoEvent")
			// return false
		}
	}
	return true
}

// Interface
var _ Trigger = &thingReceivedStatusInfoChangedTrigger{}

type thingReceivedStatusInfoChangedTrigger struct {
	baseTrigger
	thing string
	from  ThingStatus
	to    ThingStatus
	subID int
}

// OnThingReceivedStatusInfoChanged triggers the rule when the thing status changed.
// This is an equivalent of the DSL rule:
//
// Thing <thingUID> changed
func OnThingReceivedStatusInfoChanged(thing string) *thingReceivedStatusInfoChangedTrigger {
	return &thingReceivedStatusInfoChangedTrigger{
		thing: thing,
	}
}

// OnThingReceivedStatusInfoChanged triggers the rule when the thing status changed from a value.
// pass ThingStatusAny or "" to status to receive ANY update of the status
// This is an equivalent of the DSL rule:
//
// Thing <thingUID> changed from <status>
func OnThingReceivedStatusInfoChangedFrom(thing string, from ThingStatus) *thingReceivedStatusInfoChangedTrigger {
	return &thingReceivedStatusInfoChangedTrigger{
		thing: thing,
		from:  from,
	}
}

// OnThingReceivedStatusInfoChanged triggers the rule when the thing status changed to a value.
// pass ThingStatusAny or "" to status to receive ANY update of the status
// This is an equivalent of the DSL rule:
//
// Thing <thingUID> changed to <status>
func OnThingReceivedStatusInfoChangedTo(thing string, to ThingStatus) *thingReceivedStatusInfoChangedTrigger {
	return &thingReceivedStatusInfoChangedTrigger{
		thing: thing,
		to:    to,
	}
}

// OnThingReceivedStatusInfoChanged triggers the rule when the thing status changed from a value to another value.
// This is an equivalent of the DSL rule:
//
// Thing <thingUID> changed from <status> to <status>
func OnThingReceivedStatusInfoChangedFromTo(thing string, from, to ThingStatus) *thingReceivedStatusInfoChangedTrigger {
	return &thingReceivedStatusInfoChangedTrigger{
		thing: thing,
		from:  from,
		to:    to,
	}
}

func (c *thingReceivedStatusInfoChangedTrigger) activate(client subscriber, run func(ev event.Event), ruleData RuleData) error {
	if c.subID > 0 {
		return ErrRuleAlreadyActivated
	}
	c.subID = c.subscribe(client, c.thing, event.TypeThingStatusInfoChanged, run, c.match)
	return nil
}

func (c *thingReceivedStatusInfoChangedTrigger) deactivate(client subscriber) {
	if c.subID > 0 {
		client.unsubscribe(c.subID)
		c.subID = 0
	}
}

func (c *thingReceivedStatusInfoChangedTrigger) match(e event.Event) bool {
	if c.from != "" {
		// check for the desired state
		if ev, ok := e.(event.ThingStatusInfoChangedEvent); ok {
			if ev.PreviousStatus != string(c.from) {
				// not the value we wanted
				return false
			}
		} else {
			panic("expected event of type event.ThingStatusInfoChangedEvent")
			// return false
		}
	}
	if c.to != "" {
		// check for the desired state
		if ev, ok := e.(event.ThingStatusInfoChangedEvent); ok {
			if ev.NewStatus != string(c.to) {
				// not the value we wanted
				return false
			}
		} else {
			panic("expected event of type event.ThingStatusInfoChangedEvent")
			// return false
		}
	}
	return true
}

// Interface
var _ Trigger = &thingReceivedStatusInfoChangedTrigger{}
