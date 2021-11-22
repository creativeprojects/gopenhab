package openhab

import "github.com/creativeprojects/gopenhab/event"

// systemEventTrigger for connection or disconnection
type systemEventTrigger struct {
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
func OnConnect() *systemEventTrigger {
	return &systemEventTrigger{
		eventType: event.TypeClientConnected,
	}
}

// OnDisconnect is a trigger activated when the connection to openHAB is lost
func OnDisconnect() *systemEventTrigger {
	return &systemEventTrigger{
		eventType: event.TypeClientDisconnected,
	}
}

// OnStart is a trigger activated when the client has started. At this stage, the client may be connected to the openHAB event bus.
func OnStart() *systemEventTrigger {
	return &systemEventTrigger{
		eventType: event.TypeClientStarted,
	}
}

// OnStop is a trigger activated when the client is about to stop.
// There's no guarantee it is going to be the last running event, some other queued events may still run after.
func OnStop() *systemEventTrigger {
	return &systemEventTrigger{
		eventType: event.TypeClientStopped,
	}
}

// OnError is a trigger activated when the client wasn't able to contact openHAB server
func OnError() *systemEventTrigger {
	return &systemEventTrigger{
		eventType: event.TypeClientError,
	}
}

// activate subscribes to the corresponding event
func (c *systemEventTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	c.subId = client.subscribe("", c.eventType, func(e event.Event) {
		if run == nil {
			return
		}
		run(e)
	})
	return nil
}

func (c *systemEventTrigger) deactivate(client *Client) {
	if c.subId > 0 {
		client.unsubscribe(c.subId)
		c.subId = 0
	}
}

func (c *systemEventTrigger) match(e event.Event) bool {
	return e.Type() == c.eventType
}

// Interface
var _ Trigger = &systemEventTrigger{}
