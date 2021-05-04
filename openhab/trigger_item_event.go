package openhab

import (
	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
)

type ItemReceivedCommandTrigger struct {
	item  string
	state StateValue
	subId int
}

// OnItemReceivedCommand triggers the rule when the item received a command equal to state.
// pass nil to state to receive ANY command
func OnItemReceivedCommand(item string, state StateValue) *ItemReceivedCommandTrigger {
	return &ItemReceivedCommandTrigger{
		item:  item,
		state: state,
	}
}

func (c *ItemReceivedCommandTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	c.subId = client.eventBus.subscribe("smarthome/items/"+c.item+"/"+api.TopicEventCommand, event.ItemCommand, func(e event.Event) {
		if run == nil {
			return
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

type ItemReceivedUpdateTrigger struct {
	item  string
	state StateValue
}

// OnItemReceivedUpdate triggers the rule when the item received an update equal to state.
// pass nil to state to receive ANY update of the state
func OnItemReceivedUpdate(item string, state StateValue) *ItemReceivedUpdateTrigger {
	return &ItemReceivedUpdateTrigger{
		item:  item,
		state: state,
	}
}

type ItemChangedTrigger struct {
	item string
	from StateValue
	to   StateValue
}

func OnItemChanged(item string) *ItemChangedTrigger {
	return &ItemChangedTrigger{
		item: item,
	}
}

func OnItemChangedFrom(item string, from StateValue) *ItemChangedTrigger {
	return &ItemChangedTrigger{
		item: item,
		from: from,
	}
}

func OnItemChangedTo(item string, to StateValue) *ItemChangedTrigger {
	return &ItemChangedTrigger{
		item: item,
		to:   to,
	}
}

func OnItemChangedFromTo(item string, from, to StateValue) *ItemChangedTrigger {
	return &ItemChangedTrigger{
		item: item,
		from: from,
		to:   to,
	}
}
