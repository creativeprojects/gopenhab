package openhab

import "github.com/robfig/cron/v3"

type Cron struct {
	spec    string
	entryID cron.EntryID
}

// NewCron creates a trigger from a cron entry.
// Please note the YEAR field is NOT supported.
// The 6 fields are: "second minute hour dayOfMonth month dayOfWeek"
// For more information, see the quartz format:
// http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/crontrigger.html
func NewCron(spec string) *Cron {
	return &Cron{
		spec: spec,
	}
}

// activate schedules the run function in the context of a *Client
func (c Cron) activate(client *Client, run func()) error {
	entryID, err := client.cron.AddFunc(c.spec, run)
	if err != nil {
		return err
	}
	c.entryID = entryID
	return nil
}

// Interface
var _ Trigger = &Cron{}
