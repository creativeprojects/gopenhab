package openhab

import "time"

type RuleData struct {
	// ID doesn't need to be unique. It will be autogenerated if empty.
	ID string
	// Friendly name for the rule (optional)
	Name string
	// Description of the rule (optional)
	Description string
	// Context can contain anything you need to access inside the rule (database connection, etc)
	Context interface{}
	// Timeout after which the rule will be sent a cancellation through the context (optional)
	Timeout time.Duration
}
