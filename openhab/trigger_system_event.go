package openhab

import "github.com/creativeprojects/gopenhab/event"

// SystemEventTrigger for connection or disconnection
type SystemEventTrigger struct {
	eventType event.Type
	subId     int
}

// OnConnect is a trigger activated when successfully connected (or reconnected) to openHAB.
//
// Please note when openHAB is starting, the event may be triggered many times:
// you might want to look at Debounce to avoid too many triggers.
//
// On my setup, a Debounce of 1 minute works well and the event gets triggered only once
// during a restart
func OnConnect() *SystemEventTrigger {
	return &SystemEventTrigger{
		eventType: event.ClientConnected,
	}
}

// OnDisconnect is a trigger activated when the connection to openHAB is lost
func OnDisconnect() *SystemEventTrigger {
	return &SystemEventTrigger{
		eventType: event.ClientDisconnected,
	}
}

// activate subscribes to the corresponding event
func (c *SystemEventTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	c.subId = client.subscribe("", c.eventType, func(e event.Event) {
		if run == nil {
			return
		}
		run(e)
	})
	return nil
}

func (c *SystemEventTrigger) deactivate(client *Client) {
	if c.subId > 0 {
		client.unsubscribe(c.subId)
		c.subId = 0
	}
}

// Interface
var _ Trigger = &SystemEventTrigger{}
