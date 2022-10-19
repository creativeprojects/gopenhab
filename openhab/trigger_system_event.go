package openhab

import "github.com/creativeprojects/gopenhab/event"

// systemEventTrigger for connection or disconnection
type systemEventTrigger struct {
	baseTrigger
	eventType event.Type
	subId     int
}

// OnConnect is a trigger activated when successfully connected (or reconnected) to openHAB.
//
// # This event is only activated the first time openHAB sends any data to the event bus API
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

// OnStableConnection is a trigger activated when the connection has been stable for some time.
// When openHAB restarts, a connection to the API could be disconnected many times until openHAB finished
// all its initialization process.
func OnStableConnection() *systemEventTrigger {
	return &systemEventTrigger{
		eventType: event.TypeClientConnectionStable,
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

// OnAlive is a trigger activated when the client receives an ALIVE event from the server
func OnAlive() *systemEventTrigger {
	return &systemEventTrigger{
		eventType: event.TypeServerAlive,
	}
}

// OnStartlevel is a trigger activated when the server sends Startlevel events during startup.
func OnStartlevel() *systemEventTrigger {
	return &systemEventTrigger{
		eventType: event.TypeServerStartlevel,
	}
}

// activate subscribes to the corresponding event
func (c *systemEventTrigger) activate(client subscriber, run func(ev event.Event), ruleData RuleData) error {
	c.subId = c.subscribe(client, "", c.eventType, run, c.match)
	return nil
}

func (c *systemEventTrigger) deactivate(client subscriber) {
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
