package openhab

import "errors"

var (
	ErrorNotFound           = errors.New("not found")
	ErrRuleAlreadyActivated = errors.New("rule already activated")
)
