package openhab

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPastSchedule(t *testing.T) {
	t.Parallel()
	schedule := dateTimeSchedule{time.Now().Add(-time.Minute)}
	assert.Zero(t, schedule.Next(time.Now()))
}

func TestFutureSchedule(t *testing.T) {
	t.Parallel()
	schedule := dateTimeSchedule{time.Now().Add(time.Minute)}
	assert.NotZero(t, schedule.Next(time.Now()))
}

func TestDateTimeTriggerPast(t *testing.T) {
	t.Parallel()
	trigger := OnDateTime(time.Now().Add(-time.Minute))
	assert.NotNil(t, trigger)
}
