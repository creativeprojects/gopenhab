package openhab

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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

type DecimalState struct {
	value float64
	unit  string
}

// NewDecimalState creates a DecimalState with a unit
func NewDecimalState(value float64, unit string) DecimalState {
	return DecimalState{
		value: math.Round(value*100) / 100,
		unit:  unit,
	}
}

func (s DecimalState) String() string {
	value := strconv.FormatFloat(float64(s.value), 'f', -1, 64)
	if s.unit == "" {
		return value
	}
	return fmt.Sprintf("%s %s", value, s.unit)
}

func (s DecimalState) Float64() float64 {
	return s.value
}

func (s DecimalState) Unit() string {
	return s.unit
}

// ParseDecimalState converts a string to a DecimalState
func ParseDecimalState(value string) (DecimalState, error) {
	unit := ""
	// check if there's a unit first
	parts := strings.Split(value, " ")
	if len(parts) == 2 {
		value = parts[0]
		unit = parts[1]
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return DecimalState{}, err
	}
	return NewDecimalState(number, unit), nil
}

// MustParseDecimalState does not panic if the string is not a number, it returns 0 instead
func MustParseDecimalState(value string) DecimalState {
	number, _ := ParseDecimalState(value)
	return number
}
