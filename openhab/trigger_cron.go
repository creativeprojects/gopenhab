package openhab

import (
	"github.com/creativeprojects/gopenhab/event"
	"github.com/robfig/cron/v3"
)

// TimeCronTrigger triggers a rule at a time described by a quartz style cron entry
type TimeCronTrigger struct {
	spec    string
	entryID cron.EntryID
}

// OnTimeCron creates a trigger from a cron entry.
// Please note the YEAR field is NOT supported.
// The 6 fields are: "second minute hour dayOfMonth month dayOfWeek"
// For more information, see the quartz format:
// http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/crontrigger.html
func OnTimeCron(spec string) *TimeCronTrigger {
	return &TimeCronTrigger{
		spec: spec,
	}
}

// activate schedules the run function in the context of a *Client
func (c *TimeCronTrigger) activate(client *Client, run func(ev event.Event), ruleData RuleData) error {
	entryID, err := client.cron.AddFunc(c.spec, func() {
		run(event.NewSystemEvent(event.TimeCron))
	})
	if err != nil {
		return err
	}
	c.entryID = entryID
	return nil
}

func (c *TimeCronTrigger) deactivate(client *Client) {
	if c.entryID > 0 {
		client.cron.Remove(c.entryID)
		c.entryID = 0
	}
}

// Interface
var _ Trigger = &TimeCronTrigger{}
