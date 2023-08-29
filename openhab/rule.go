package openhab

import (
	"context"
	"sync"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/segmentio/ksuid"
)

type Runner func(ctx context.Context, client *Client, ruleData RuleData, e event.Event)

type rule struct {
	ruleData     RuleData
	client       *Client
	runner       Runner
	triggers     []Trigger
	runLocker    sync.Mutex
	cancelFunc   context.CancelFunc
	cancelLocker sync.Mutex
}

func newRule(client *Client, ruleData RuleData, runner Runner, triggers []Trigger) *rule {
	if ruleData.ID == "" {
		gen := ksuid.New()
		ruleData.ID = gen.String()
	}
	return &rule{
		ruleData:     ruleData,
		client:       client,
		runner:       runner,
		triggers:     triggers,
		runLocker:    sync.Mutex{},
		cancelFunc:   nil,
		cancelLocker: sync.Mutex{},
	}
}

func (r *rule) String() string {
	if r.ruleData.Name != "" {
		return r.ruleData.Name
	}
	if r.ruleData.ID != "" {
		return "ID " + r.ruleData.ID
	}
	return ""
}

func (r *rule) activate(client subscriber) error {
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

func (r *rule) deactivate(client subscriber) {
	for _, trigger := range r.triggers {
		if trigger == nil {
			continue
		}
		trigger.deactivate(client)
	}
}

func (r *rule) run(e event.Event) {
	// run only one instance of that rule at any time
	r.runLocker.Lock()
	defer r.runLocker.Unlock()

	// make the rule cancellable from the outside
	var cancelFunc context.CancelFunc
	ctx := context.Background()
	if r.ruleData.Timeout > 0 {
		ctx, cancelFunc = context.WithTimeout(ctx, r.ruleData.Timeout)
	} else {
		ctx, cancelFunc = context.WithCancel(context.Background())
	}
	r.setCancelFunc(cancelFunc)
	defer r.cancel()

	r.runner(ctx, r.client, r.ruleData, e)
}

func (r *rule) cancel() {
	r.cancelLocker.Lock()
	defer r.cancelLocker.Unlock()

	if r.cancelFunc != nil {
		r.cancelFunc()
		r.cancelFunc = nil
	}
}

func (t *rule) setCancelFunc(cancelFunc context.CancelFunc) {
	t.cancelLocker.Lock()
	defer t.cancelLocker.Unlock()

	t.cancelFunc = cancelFunc
}
