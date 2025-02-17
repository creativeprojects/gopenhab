package openhab

import "errors"

var (
	ErrNotFound             = errors.New("not found")
	ErrRuleAlreadyActivated = errors.New("rule already activated")
)
