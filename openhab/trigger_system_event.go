package openhab

import "github.com/creativeprojects/gopenhab/event"

type SystemEventTrigger struct {
	eventType event.Type
	subId     int
}

func OnConnect() *SystemEventTrigger {
	return &SystemEventTrigger{
		eventType: event.ClientConnected,
	}
}

func OnDisconnect() *SystemEventTrigger {
	return &SystemEventTrigger{
		eventType: event.ClientDisconnected,
	}
}

// activate subscribes to the corresponding event
func (c *SystemEventTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	c.subId = client.eventBus.subscribe("", c.eventType, func(e event.Event) {
		if run == nil {
			return
		}
		run(e)
	})
	return nil
}

func (c *SystemEventTrigger) deactivate(client *Client) {
	if c.subId > 0 {
		client.eventBus.unsubscribe(c.subId)
		c.subId = 0
	}
}

// Interface
var _ Trigger = &SystemEventTrigger{}
