package openhab

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOpenHABDate(t *testing.T) {
	fixtures := []string{"2022-03-08T07:01:00.123+0000", "2022-03-08T17:38:00.456+0100"}

	for _, date := range fixtures {
		t.Run(date, func(t *testing.T) {
			result, err := ParseDateTimeState(date)
			assert.NoError(t, err)
			assert.NotZero(t, result)
			assert.Equal(t, date, result.String())
		})
	}
}

func TestItemState(t *testing.T) {
	dateTime := time.Date(2022, 3, 8, 7, 1, 0, 0, time.Local)
	fixtures := []struct {
		state    State
		raw      interface{}
		value    string
		notEqual string
	}{
		{SwitchState("ON"), "ON", "ON", "OFF"},
		{StringState("test"), "test", "test", "other"},
		{NewDecimalState(2.3, "cm"), float64(2.3), "2.3 cm", "2.3"},
		{NewDateTimeState(dateTime), dateTime, "2022-03-08T07:01:00+0000", "2022-03-08T07:01:01+0000"},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.state.String(), func(t *testing.T) {
			assert.Equal(t, fixture.raw, fixture.state.Raw())
			assert.Equal(t, fixture.value, fixture.state.String())
			assert.True(t, fixture.state.Equal(fixture.value))
			assert.False(t, fixture.state.Equal(fixture.notEqual))
		})
	}
}
