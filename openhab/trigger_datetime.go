package openhab

import (
	"time"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/robfig/cron/v3"
)

type dateTimeTrigger struct {
	schedule dateTimeSchedule
	entryID  cron.EntryID
}

// OnDateTime will trigger the event only once: at the exact time.Time passed in argument.
// If the time argument is in the past, no trigger will ever run.
func OnDateTime(dateTime time.Time) *dateTimeTrigger {
	return &dateTimeTrigger{
		schedule: dateTimeSchedule{dateTime},
	}
}

// activate schedules the run function in the context of a *Client
func (c *dateTimeTrigger) activate(client subscriber, run func(ev event.Event), ruleData RuleData) error {
	entryID := client.getCron().Schedule(c.schedule, cron.FuncJob(func() {
		run(event.NewSystemEvent(event.TypeTimeCron))
	}))
	c.entryID = entryID
	return nil
}

func (c *dateTimeTrigger) deactivate(client subscriber) {
	if c.entryID > 0 {
		client.getCron().Remove(c.entryID)
		c.entryID = 0
	}
}

func (c *dateTimeTrigger) match(e event.Event) bool {
	return true
}

// Interface
var _ Trigger = &dateTimeTrigger{}

type dateTimeSchedule struct {
	next time.Time
}

func (s dateTimeSchedule) Next(after time.Time) time.Time {
	if time.Now().After(s.next) {
		// passed the activation time
		return time.Time{}
	}
	return s.next
}
