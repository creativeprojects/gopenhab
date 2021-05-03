package openhab

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

type StringState string

func (s StringState) String() string {
	return string(s)
}
