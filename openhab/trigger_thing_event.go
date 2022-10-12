package openhab

import (
	"errors"

	"github.com/creativeprojects/gopenhab/event"
)

type thingReceivedStatusInfoTrigger struct {
	thing  string
	status ThingStatus
	subId  int
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

func (c *thingReceivedStatusInfoTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	if c.subId > 0 {
		return errors.New("rule already activated")
	}
	c.subId = client.subscribe(c.thing, event.TypeItemState, func(e event.Event) {
		if run == nil {
			return
		}
		if c.match(e) {
			run(e)
		}
	})
	return nil
}

func (c *thingReceivedStatusInfoTrigger) deactivate(client *Client) {
	if c.subId > 0 {
		client.unsubscribe(c.subId)
		c.subId = 0
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
	thing string
	from  ThingStatus
	to    ThingStatus
	subId int
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

func (c *thingReceivedStatusInfoChangedTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	if c.subId > 0 {
		return errors.New("rule already activated")
	}
	c.subId = client.subscribe(c.thing, event.TypeItemState, func(e event.Event) {
		if run == nil {
			return
		}
		if c.match(e) {
			run(e)
		}
	})
	return nil
}

func (c *thingReceivedStatusInfoChangedTrigger) deactivate(client *Client) {
	if c.subId > 0 {
		client.unsubscribe(c.subId)
		c.subId = 0
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
			panic("expected event of type event.ThingStatusInfoEvent")
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
			panic("expected event of type event.ThingStatusInfoEvent")
			// return false
		}
	}
	return true
}

// Interface
var _ Trigger = &thingReceivedStatusInfoChangedTrigger{}
