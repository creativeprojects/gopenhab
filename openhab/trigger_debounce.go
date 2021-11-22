package openhab

import (
	"sync"
	"time"

	"github.com/creativeprojects/gopenhab/event"
)

// triggerDebounce is a special trigger used to debounce a trigger
type triggerDebounce struct {
	lock    sync.Locker
	after   time.Duration
	timer   *time.Timer
	trigger Trigger
}

// Debounce will trigger the event after some time, in case the subscription is triggered multiple times in a row.
// Typically this is the case of Connection and Disconnection system events when openHAB is starting
func Debounce(trigger Trigger, after time.Duration) *triggerDebounce {
	return &triggerDebounce{
		lock:    &sync.Mutex{},
		after:   after,
		trigger: trigger,
	}
}

func (c *triggerDebounce) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	debounced := func(ev event.Event) {
		c.lock.Lock()
		defer c.lock.Unlock()

		if c.timer != nil {
			c.timer.Stop()
		}

		c.timer = time.AfterFunc(c.after, func() {
			run(ev)
		})
	}
	return c.trigger.activate(client, debounced, ruleData)
}

func (c *triggerDebounce) deactivate(client *Client) {
	c.trigger.deactivate(client)

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.timer != nil {
		c.timer.Stop()
	}
}

func (c *triggerDebounce) match(e event.Event) bool {
	return c.trigger.match(e)
}

// Interface
var _ Trigger = &triggerDebounce{}
