package openhab

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const DateTimeFormat = "2006-01-02T15:04:05.999-0700"

type State interface {
	String() string
	Raw() interface{}
	Equal(other string) bool
}

// Verify interfaces
var (
	_ State = SwitchState("")
	_ State = StringState("")
	_ State = DecimalState{}
	_ State = DateTimeState{}
)

type SwitchState string

const (
	SwitchON  SwitchState = "ON"
	SwitchOFF SwitchState = "OFF"
)

func (s SwitchState) String() string {
	return strings.ToUpper(string(s))
}

func (s SwitchState) Raw() interface{} {
	return string(s)
}

func (s SwitchState) Equal(other string) bool {
	return s.String() == other
}

const (
	StateNULL = "NULL"
	StateOFF  = "OFF"
	StateON   = "ON"
)

type StringState string

func (s StringState) String() string {
	return string(s)
}

func (s StringState) Raw() interface{} {
	return string(s)
}

func (s StringState) Equal(other string) bool {
	return string(s) == other
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

func (s DecimalState) Raw() interface{} {
	return s.value
}

func (s DecimalState) Float64() float64 {
	return s.value
}

func (s DecimalState) Unit() string {
	return s.unit
}

func (s DecimalState) Equal(other string) bool {
	compare, err := ParseDecimalState(other)
	if err != nil {
		return false
	}
	return s.value == compare.value && s.unit == compare.unit
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

type DateTimeState time.Time

// NewDateTimeState creates a DateState
func NewDateTimeState(value time.Time) DateTimeState {
	return DateTimeState(value)
}

func (s DateTimeState) String() string {
	return time.Time(s).Format(DateTimeFormat)
}

func (s DateTimeState) Raw() interface{} {
	return time.Time(s)
}

func (s DateTimeState) Time() time.Time {
	return time.Time(s)
}

func (s DateTimeState) Equal(other string) bool {
	compare, err := ParseDateTimeState(other)
	if err != nil {
		return false
	}
	return s.Time().Equal(compare.Time())
}

// ParseDateTimeState converts a string to a DateState
func ParseDateTimeState(value string) (DateTimeState, error) {
	date, err := time.Parse(DateTimeFormat, value)
	if err != nil {
		return DateTimeState{}, err
	}
	return NewDateTimeState(date), nil
}

// MustParseDateTimeState does not panic if the string is not a number, it returns 0 instead
func MustParseDateTimeState(value string) DateTimeState {
	number, _ := ParseDateTimeState(value)
	return number
}
