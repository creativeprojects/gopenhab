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
		{
			event.NewThingStatusInfoChangedEvent("TestItem", event.ThingStatus{}, event.ThingStatus{}),
			OnThingReceivedStatusInfoChanged("TestItem"),
			true,
		},
		{
			event.NewThingStatusInfoChangedEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOffline)}, event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfoChanged("TestItem"),
			true,
		},
		{
			event.NewThingStatusInfoChangedEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOffline)}, event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfoChangedFrom("TestItem", ThingStatusOffline),
			true,
		},
		{
			event.NewThingStatusInfoChangedEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOffline)}, event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfoChangedFrom("TestItem", ThingStatusOnline),
			false,
		},
		{
			event.NewThingStatusInfoChangedEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOffline)}, event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfoChangedTo("TestItem", ThingStatusOnline),
			true,
		},
		{
			event.NewThingStatusInfoChangedEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOffline)}, event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfoChangedTo("TestItem", ThingStatusOffline),
			false,
		},
		{
			event.NewThingStatusInfoChangedEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOffline)}, event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfoChangedFromTo("TestItem", ThingStatusOnline, ThingStatusOffline),
			false,
		},
		{
			event.NewThingStatusInfoChangedEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOffline)}, event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfoChangedFromTo("TestItem", ThingStatusOnline, ThingStatusOnline),
			false,
		},
		{
			event.NewThingStatusInfoChangedEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOffline)}, event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfoChangedFromTo("TestItem", ThingStatusOffline, ThingStatusOffline),
			false,
		},
		{
			event.NewThingStatusInfoChangedEvent("TestItem", event.ThingStatus{Status: string(ThingStatusOffline)}, event.ThingStatus{Status: string(ThingStatusOnline)}),
			OnThingReceivedStatusInfoChangedFromTo("TestItem", ThingStatusOffline, ThingStatusOnline),
			true,
		},
	}

	for _, testEvent := range testEvents {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, testEvent.match, testEvent.trigger.match(testEvent.e))
		})
	}
}
