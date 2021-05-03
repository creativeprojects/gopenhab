package openhab

type Rule struct {
	config   RuleConfig
	run      func()
	triggers []Trigger
}

func NewRule(config RuleConfig, run func(), triggers []Trigger) *Rule {
	return &Rule{
		config:   config,
		run:      run,
		triggers: triggers,
	}
}

func (r Rule) String() string {
	if r.config.Name != "" {
		return r.config.Name
	}
	if r.config.ID != "" {
		return "ID " + r.config.ID
	}
	return ""
}

func (r Rule) activate(client *Client) error {
	for _, trigger := range r.triggers {
		err := trigger.activate(client, r.run)
		if err != nil {
			return err
		}
	}
	return nil
}
