package openhab

import (
	"testing"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
)

func TestMatchingThingEvent(t *testing.T) {
	testEvents := []struct {
		e       event.Event
		trigger Trigger
		match   bool
	}{
		// received status
		{
			event.NewThingStatusInfoEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfo("TestItem", ThingStatusOffline),
			false,
		},
		{
			event.NewThingStatusInfoEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfo("TestItem", ThingStatusOnline),
			true,
		},
		{
			event.NewThingStatusInfoEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfo("TestItem", ""),
			true,
		},
		{
			event.NewThingStatusInfoEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOffline)}),
			OnThingReceivedStatusInfo("TestItem", ThingStatusAny),
			true,
		},
		// thing status changed
	}

	for _, testEvent := range testEvents {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, testEvent.match, testEvent.trigger.match(testEvent.e))
		})
	}
}
