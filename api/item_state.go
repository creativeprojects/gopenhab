package api

type StateValue interface {
	String() string
}

type SwitchState string

const (
	SwitchON  SwitchState = "ON"
	SwitchOFF SwitchState = "OFF"
)

func (s SwitchState) String() string {
	return string(s)
}

// Verify interface
var _ StateValue = SwitchState("")

const (
	StateNULL = "NULL"
	StateOFF  = "OFF"
	StateON   = "ON"
)
