package openhab

import (
	"github.com/creativeprojects/gopenhab/event"
	"github.com/robfig/cron/v3"
)

// timeCronTrigger triggers a rule at a time described by a quartz style cron entry
type timeCronTrigger struct {
	spec    string
	entryID cron.EntryID
}

// OnTimeCron creates a trigger from a cron entry.
// Please note the YEAR field is NOT supported.
// The 6 fields are: "second minute hour dayOfMonth month dayOfWeek"
// For more information, see the quartz format:
// http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/crontrigger.html
func OnTimeCron(spec string) *timeCronTrigger {
	return &timeCronTrigger{
		spec: spec,
	}
}

// activate schedules the run function in the context of a *Client
func (c *timeCronTrigger) activate(client subscriber, run func(ev event.Event), ruleData RuleData) error {
	entryID, err := client.getCron().AddFunc(c.spec, func() {
		defer preventPanic()

		run(event.NewSystemEvent(event.TypeTimeCron))
	})
	if err != nil {
		return err
	}
	c.entryID = entryID
	return nil
}

func (c *timeCronTrigger) deactivate(client subscriber) {
	if c.entryID > 0 {
		client.getCron().Remove(c.entryID)
		c.entryID = 0
	}
}

func (c *timeCronTrigger) match(e event.Event) bool {
	return true
}

// Interface
var _ Trigger = &timeCronTrigger{}
