package openhab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenHABDate(t *testing.T) {
	fixtures := []string{"2022-03-08T07:01:00.000+0000", "2022-03-08T17:38:00.000+0000"}

	for _, date := range fixtures {
		t.Run(date, func(t *testing.T) {
			result, err := ParseDateTimeState(date)
			assert.NoError(t, err)
			assert.NotZero(t, result)
		})
	}
}
