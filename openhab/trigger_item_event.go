package openhab

import (
	"errors"

	"github.com/creativeprojects/gopenhab/event"
)

type ItemReceivedCommandTrigger struct {
	item  string
	state StateValue
	subId int
}

// OnItemReceivedCommand triggers the rule when the item received a command equal to state.
// Use a nil state to receive ANY command sent to the item
// This is an equivalent of the DSL rule:
//
// Item <item> received command [<command>]
func OnItemReceivedCommand(item string, state StateValue) *ItemReceivedCommandTrigger {
	return &ItemReceivedCommandTrigger{
		item:  item,
		state: state,
	}
}

func (c *ItemReceivedCommandTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	if c.subId > 0 {
		return errors.New("rule already activated")
	}
	c.subId = client.subscribe(c.item, event.TypeItemCommand, func(e event.Event) {
		if run == nil {
			return
		}
		if c.match(e) {
			run(e)
		}
	})
	return nil
}

func (c *ItemReceivedCommandTrigger) deactivate(client *Client) {
	if c.subId > 0 {
		client.unsubscribe(c.subId)
		c.subId = 0
	}
}

func (c *ItemReceivedCommandTrigger) match(e event.Event) bool {
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
var _ Trigger = &ItemReceivedCommandTrigger{}

type ItemReceivedStateTrigger struct {
	item  string
	state StateValue
	subId int
}

// OnItemReceivedState triggers the rule when the item received an update equal to state.
// pass nil to state to receive ANY update of the state
// This is an equivalent of the DSL rule:
//
// Item <item> received update [<state>]
func OnItemReceivedState(item string, state StateValue) *ItemReceivedStateTrigger {
	return &ItemReceivedStateTrigger{
		item:  item,
		state: state,
	}
}

func (c *ItemReceivedStateTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	if c.subId > 0 {
		return errors.New("rule already activated")
	}
	c.subId = client.subscribe(c.item, event.TypeItemState, func(e event.Event) {
		if run == nil {
			return
		}
		if c.match(e) {
			run(e)
		}
	})
	return nil
}

func (c *ItemReceivedStateTrigger) deactivate(client *Client) {
	if c.subId > 0 {
		client.unsubscribe(c.subId)
		c.subId = 0
	}
}

func (c *ItemReceivedStateTrigger) match(e event.Event) bool {
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
var _ Trigger = &ItemReceivedStateTrigger{}

type ItemStateChangedTrigger struct {
	item   string
	from   StateValue
	to     StateValue
	subId1 int
	subId2 int
}

// OnItemStateChanged triggers the rule when the item received an update with a different state
// This is an equivalent of the DSL rule:
//
// Item <item> changed
func OnItemStateChanged(item string) *ItemStateChangedTrigger {
	return &ItemStateChangedTrigger{
		item: item,
	}
}

// OnItemStateChangedFrom triggers the rule when the item received an update with a different state
// This is an equivalent of the DSL rule:
//
// Item <item> changed from <state>
func OnItemStateChangedFrom(item string, from StateValue) *ItemStateChangedTrigger {
	return &ItemStateChangedTrigger{
		item: item,
		from: from,
	}
}

// OnItemStateChangedTo triggers the rule when the item received an update with a different state
// This is an equivalent of the DSL rule:
//
// Item <item> changed to <state>
func OnItemStateChangedTo(item string, to StateValue) *ItemStateChangedTrigger {
	return &ItemStateChangedTrigger{
		item: item,
		to:   to,
	}
}

// OnItemStateChangedFromTo triggers the rule when the item received an update with a different state
// This is an equivalent of the DSL rule:
//
// Item <item> changed from <state> to <state>
func OnItemStateChangedFromTo(item string, from, to StateValue) *ItemStateChangedTrigger {
	return &ItemStateChangedTrigger{
		item: item,
		from: from,
		to:   to,
	}
}

func (c *ItemStateChangedTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	if run == nil {
		return errors.New("event callback is nil")
	}
	if c.subId1 > 0 || c.subId2 > 0 {
		return errors.New("rule already activated")
	}
	c.subId1 = client.subscribe(c.item, event.TypeItemStateChanged, func(e event.Event) {
		if c.match(e) {
			run(e)
		}
	})
	c.subId2 = client.subscribe(c.item, event.TypeGroupItemStateChanged, func(e event.Event) {
		if c.match(e) {
			run(e)
		}
	})
	return nil
}

func (c *ItemStateChangedTrigger) deactivate(client *Client) {
	if c.subId1 > 0 {
		client.unsubscribe(c.subId1)
		c.subId1 = 0
	}
	if c.subId2 > 0 {
		client.unsubscribe(c.subId2)
		c.subId2 = 0
	}
}

func (c *ItemStateChangedTrigger) match(e event.Event) bool {
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
var _ Trigger = &ItemStateChangedTrigger{}
