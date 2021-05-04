package openhab

import "github.com/creativeprojects/gopenhab/event"

type Trigger interface {
	// activate the trigger for func() in the context of a *Client
	activate(client *Client, run func(ev event.Event), ruleData RuleData) error
	// deactivate the trigger in the context of a *Client
	deactivate(client *Client)
}
