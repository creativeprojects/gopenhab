package openhab

import (
	"testing"

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
