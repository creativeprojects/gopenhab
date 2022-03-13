package openhab

import (
	"sync"
	"time"

	"github.com/creativeprojects/gopenhab/event"
)

// triggerDebounce is a special trigger used to debounce a trigger
type triggerDebounce struct {
	lock     sync.Locker
	after    time.Duration
	timer    *time.Timer
	triggers []Trigger
}

// Debounce will trigger the event after some time, in case the subscription is triggered multiple times in a row.
// Typically this is the case of Connection and Disconnection system events when openHAB is starting
func Debounce(after time.Duration, triggers ...Trigger) *triggerDebounce {
	return &triggerDebounce{
		lock:     &sync.Mutex{},
		after:    after,
		triggers: triggers,
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
	for _, trigger := range c.triggers {
		err := trigger.activate(client, debounced, ruleData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *triggerDebounce) deactivate(client *Client) {
	for _, trigger := range c.triggers {
		trigger.deactivate(client)
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.timer != nil {
		c.timer.Stop()
	}
}

func (c *triggerDebounce) match(e event.Event) bool {
	for _, trigger := range c.triggers {
		if trigger.match(e) {
			return true
		}
	}
	return false
}

// Interface
var _ Trigger = &triggerDebounce{}
