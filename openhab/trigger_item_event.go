package openhab

import (
	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
)

const itemTopicPrefix = "smarthome/items/"

type ItemReceivedCommandTrigger struct {
	topic string
	state StateValue
	subId int
}

// OnItemReceivedCommand triggers the rule when the item received a command equal to state.
// pass nil to state to receive ANY command
func OnItemReceivedCommand(item string, state StateValue) *ItemReceivedCommandTrigger {
	return &ItemReceivedCommandTrigger{
		topic: itemTopicPrefix + item + "/" + api.TopicEventCommand,
		state: state,
	}
}

func (c *ItemReceivedCommandTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	c.subId = client.eventBus.subscribe(c.topic, event.ItemCommand, func(e event.Event) {
		if run == nil {
			return
		}
		if c.state != nil && c.state.String() != "" {
			// check for the desired state
			if ev, ok := e.(event.ItemReceivedCommand); ok {
				if ev.Command != c.state.String() {
					// not the value we wanted
					return
				}
			}
		}
		run(e)
	})
	return nil
}

func (c *ItemReceivedCommandTrigger) deactivate(client *Client) {
	if c.subId > 0 {
		client.eventBus.unsubscribe(c.subId)
		c.subId = 0
	}
}

// Interface
var _ Trigger = &ItemReceivedCommandTrigger{}

type ItemReceivedStateTrigger struct {
	topic string
	state StateValue
	subId int
}

// OnItemReceivedState triggers the rule when the item received an update equal to state.
// pass nil to state to receive ANY update of the state
func OnItemReceivedState(item string, state StateValue) *ItemReceivedStateTrigger {
	return &ItemReceivedStateTrigger{
		topic: itemTopicPrefix + item + "/" + api.TopicEventState,
		state: state,
	}
}

func (c *ItemReceivedStateTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	c.subId = client.eventBus.subscribe(c.topic, event.ItemState, func(e event.Event) {
		if run == nil {
			return
		}
		if c.state != nil && c.state.String() != "" {
			// check for the desired state
			if ev, ok := e.(event.ItemReceivedState); ok {
				if ev.State != c.state.String() {
					// not the value we wanted
					return
				}
			}
		}
		run(e)
	})
	return nil
}

func (c *ItemReceivedStateTrigger) deactivate(client *Client) {
	if c.subId > 0 {
		client.eventBus.unsubscribe(c.subId)
		c.subId = 0
	}
}

// Interface
var _ Trigger = &ItemReceivedStateTrigger{}

type ItemChangedTrigger struct {
	topic string
	from  StateValue
	to    StateValue
	subId int
}

func OnItemChanged(item string) *ItemChangedTrigger {
	return &ItemChangedTrigger{
		topic: itemTopicPrefix + item + "/" + api.TopicEventStateChanged,
	}
}

func OnItemChangedFrom(item string, from StateValue) *ItemChangedTrigger {
	return &ItemChangedTrigger{
		topic: itemTopicPrefix + item + "/" + api.TopicEventStateChanged,
		from:  from,
	}
}

func OnItemChangedTo(item string, to StateValue) *ItemChangedTrigger {
	return &ItemChangedTrigger{
		topic: itemTopicPrefix + item + "/" + api.TopicEventStateChanged,
		to:    to,
	}
}

func OnItemChangedFromTo(item string, from, to StateValue) *ItemChangedTrigger {
	return &ItemChangedTrigger{
		topic: itemTopicPrefix + item + "/" + api.TopicEventStateChanged,
		from:  from,
		to:    to,
	}
}

func (c *ItemChangedTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	c.subId = client.eventBus.subscribe(c.topic, event.ItemStateChanged, func(e event.Event) {
		if run == nil {
			return
		}
		if c.from != nil && c.from.String() != "" {
			// check for the desired state
			if ev, ok := e.(event.ItemChanged); ok {
				if ev.OldState != c.from.String() {
					// not the value we wanted
					return
				}
			}
		}
		if c.to != nil && c.to.String() != "" {
			// check for the desired state
			if ev, ok := e.(event.ItemChanged); ok {
				if ev.State != c.to.String() {
					// not the value we wanted
					return
				}
			}
		}
		run(e)
	})
	return nil
}

func (c *ItemChangedTrigger) deactivate(client *Client) {
	if c.subId > 0 {
		client.eventBus.unsubscribe(c.subId)
		c.subId = 0
	}
}

// Interface
var _ Trigger = &ItemChangedTrigger{}
