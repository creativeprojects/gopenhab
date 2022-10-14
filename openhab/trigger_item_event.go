package openhab

import (
	"errors"

	"github.com/creativeprojects/gopenhab/event"
)

type itemReceivedCommandTrigger struct {
	baseTrigger
	item  string
	state State
	subId int
}

// OnItemReceivedCommand triggers the rule when the item received a command equal to state.
// Use a nil state to receive ANY command sent to the item
// This is an equivalent of the DSL rule:
//
// Item <item> received command [<command>]
func OnItemReceivedCommand(item string, state State) *itemReceivedCommandTrigger {
	return &itemReceivedCommandTrigger{
		item:  item,
		state: state,
	}
}

func (c *itemReceivedCommandTrigger) activate(client subscriber, run func(ev event.Event), ruleData RuleData) error {
	if c.subId > 0 {
		return errors.New("rule already activated")
	}
	c.subId = c.subscribe(client, c.item, event.TypeItemCommand, run, c.match)
	return nil
}

func (c *itemReceivedCommandTrigger) deactivate(client subscriber) {
	if c.subId > 0 {
		client.unsubscribe(c.subId)
		c.subId = 0
	}
}

func (c *itemReceivedCommandTrigger) match(e event.Event) bool {
	if c.state != nil && c.state.String() != "" {
		// check for the desired state
		if ev, ok := e.(event.ItemReceivedCommand); ok {
			if ev.Command != c.state.String() {
				// not the value we wanted
				return false
			}
		} else {
			panic("expected event of type event.ItemReceivedCommand")
			// return false
		}
	}
	return true
}

// Interface
var _ Trigger = &itemReceivedCommandTrigger{}

type itemReceivedStateTrigger struct {
	baseTrigger
	item  string
	state State
	subId int
}

// OnItemReceivedState triggers the rule when the item received an update equal to state.
// pass nil to state to receive ANY update of the state
// This is an equivalent of the DSL rule:
//
// Item <item> received update [<state>]
func OnItemReceivedState(item string, state State) *itemReceivedStateTrigger {
	return &itemReceivedStateTrigger{
		item:  item,
		state: state,
	}
}

func (c *itemReceivedStateTrigger) activate(client subscriber, run func(ev event.Event), ruleData RuleData) error {
	if c.subId > 0 {
		return errors.New("rule already activated")
	}
	c.subId = c.subscribe(client, c.item, event.TypeItemState, run, c.match)
	return nil
}

func (c *itemReceivedStateTrigger) deactivate(client subscriber) {
	if c.subId > 0 {
		client.unsubscribe(c.subId)
		c.subId = 0
	}
}

func (c *itemReceivedStateTrigger) match(e event.Event) bool {
	if c.state != nil && c.state.String() != "" {
		// check for the desired state
		if ev, ok := e.(event.ItemReceivedState); ok {
			if ev.State != c.state.String() {
				// not the value we wanted
				return false
			}
		} else {
			panic("expected event of type event.ItemReceivedState")
			// return false
		}
	}
	return true
}

// Interface
var _ Trigger = &itemReceivedStateTrigger{}

type itemStateChangedTrigger struct {
	baseTrigger
	item   string
	from   State
	to     State
	subId1 int
	subId2 int
}

// OnItemStateChanged triggers the rule when the item received an update with a different state
// This is an equivalent of the DSL rule:
//
// Item <item> changed
func OnItemStateChanged(item string) *itemStateChangedTrigger {
	return &itemStateChangedTrigger{
		item: item,
	}
}

// OnItemStateChangedFrom triggers the rule when the item received an update with a different state
// This is an equivalent of the DSL rule:
//
// Item <item> changed from <state>
func OnItemStateChangedFrom(item string, from State) *itemStateChangedTrigger {
	return &itemStateChangedTrigger{
		item: item,
		from: from,
	}
}

// OnItemStateChangedTo triggers the rule when the item received an update with a different state
// This is an equivalent of the DSL rule:
//
// Item <item> changed to <state>
func OnItemStateChangedTo(item string, to State) *itemStateChangedTrigger {
	return &itemStateChangedTrigger{
		item: item,
		to:   to,
	}
}

// OnItemStateChangedFromTo triggers the rule when the item received an update with a different state
// This is an equivalent of the DSL rule:
//
// Item <item> changed from <state> to <state>
func OnItemStateChangedFromTo(item string, from, to State) *itemStateChangedTrigger {
	return &itemStateChangedTrigger{
		item: item,
		from: from,
		to:   to,
	}
}

func (c *itemStateChangedTrigger) activate(client subscriber, run func(ev event.Event), ruleData RuleData) error {
	if run == nil {
		return errors.New("event callback is nil")
	}
	if c.subId1 > 0 || c.subId2 > 0 {
		return errors.New("rule already activated")
	}
	c.subId1 = c.subscribe(client, c.item, event.TypeItemStateChanged, run, c.match)
	c.subId1 = c.subscribe(client, c.item, event.TypeGroupItemStateChanged, run, c.match)
	return nil
}

func (c *itemStateChangedTrigger) deactivate(client subscriber) {
	if c.subId1 > 0 {
		client.unsubscribe(c.subId1)
		c.subId1 = 0
	}
	if c.subId2 > 0 {
		client.unsubscribe(c.subId2)
		c.subId2 = 0
	}
}

func (c *itemStateChangedTrigger) match(e event.Event) bool {
	if c.from != nil && c.from.String() != "" {
		// check for the desired state
		if ev, ok := e.(event.ItemStateChanged); ok {
			if ev.PreviousState != c.from.String() {
				// not the value we wanted
				return false
			}
		} else if ev, ok := e.(event.GroupItemStateChanged); ok {
			if ev.PreviousState != c.from.String() {
				// not the value we wanted
				return false
			}
		} else {
			panic("expected event of type event.ItemStateChanged or event.GroupItemStateChanged")
			// return false
		}
	}
	if c.to != nil && c.to.String() != "" {
		// check for the desired state
		if ev, ok := e.(event.ItemStateChanged); ok {
			if ev.NewState != c.to.String() {
				// not the value we wanted
				return false
			}
		} else if ev, ok := e.(event.GroupItemStateChanged); ok {
			if ev.NewState != c.to.String() {
				// not the value we wanted
				return false
			}
		} else {
			panic("expected event of type event.ItemStateChanged or event.GroupItemStateChanged")
			// return false
		}
	}
	return true
}

// Interface
var _ Trigger = &itemStateChangedTrigger{}
