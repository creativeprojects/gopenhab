package openhab

import (
	"sync"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/segmentio/ksuid"
)

type Runner func(client RuleClient, ruleData RuleData, e event.Event)

type Rule struct {
	ruleData  RuleData
	client    *Client
	runner    Runner
	triggers  []Trigger
	runLocker sync.Mutex
}

func newRule(client *Client, ruleData RuleData, runner Runner, triggers []Trigger) *Rule {
	if ruleData.ID == "" {
		gen := ksuid.New()
		ruleData.ID = gen.String()
	}
	return &Rule{
		ruleData: ruleData,
		client:   client,
		runner:   runner,
		triggers: triggers,
	}
}

func (r *Rule) String() string {
	if r.ruleData.Name != "" {
		return r.ruleData.Name
	}
	if r.ruleData.ID != "" {
		return "ID " + r.ruleData.ID
	}
	return ""
}

func (r *Rule) activate(client *Client) error {
	for _, trigger := range r.triggers {
		if trigger == nil {
			errorlog.Printf("nil trigger encountered")
			continue
		}
		err := trigger.activate(client, r.run, r.ruleData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Rule) deactivate(client *Client) {
	for _, trigger := range r.triggers {
		if trigger == nil {
			continue
		}
		trigger.deactivate(client)
	}
}

func (r *Rule) run(e event.Event) {
	// only run one instance of that rule at any time
	r.runLocker.Lock()
	defer r.runLocker.Unlock()

	r.runner(r.client, r.ruleData, e)
}
