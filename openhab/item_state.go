package openhab

import (
	"math"
	"strconv"
)

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

type DecimalState float64

// NewDecimalState creates a DecimalState rounded to 2 decimal places
func NewDecimalState(value float64) DecimalState {
	return DecimalState(math.Round(value*100) / 100)
}

func (s DecimalState) String() string {
	return strconv.FormatFloat(float64(s), 'f', -1, 64)
}

// ParseDecimalState converts a string to a DecimalState
func ParseDecimalState(value string) (DecimalState, error) {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, err
	}
	return DecimalState(number), nil
}

// MustParseDecimalState does not panic if the string is not a number, it returns 0 instead
func MustParseDecimalState(value string) DecimalState {
	number, err := ParseDecimalState(value)
	if err != nil {
		return 0.0
	}
	return number
}
