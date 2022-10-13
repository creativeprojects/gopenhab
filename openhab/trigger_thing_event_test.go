package openhab

import (
	"testing"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
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

func TestSubscription(t *testing.T) {
	testFixtures := []struct {
		trigger   Trigger
		eventType event.Type
		event     event.Event
		calls     int
	}{
		{
			OnThingReceivedStatusInfo("thing", ThingStatusOnline),
			event.TypeThingStatusInfo,
			event.NewThingStatusInfoEvent("thing", event.ThingStatus{Status: string(ThingStatusOffline)}),
			0,
		},
		{
			OnThingReceivedStatusInfo("thing", ThingStatusOnline),
			event.TypeThingStatusInfo,
			event.NewThingStatusInfoEvent("thing", event.ThingStatus{Status: string(ThingStatusOnline)}),
			1,
		},
		{
			OnThingReceivedStatusInfoChangedFrom("thing", ThingStatusOnline),
			event.TypeThingStatusInfoChanged,
			event.NewThingStatusInfoChangedEvent("thing",
				event.ThingStatus{Status: string(ThingStatusOffline)},
				event.ThingStatus{Status: string(ThingStatusOnline)},
			),
			0,
		},
		{
			OnThingReceivedStatusInfoChangedFrom("thing", ThingStatusOnline),
			event.TypeThingStatusInfoChanged,
			event.NewThingStatusInfoChangedEvent("thing",
				event.ThingStatus{Status: string(ThingStatusOnline)},
				event.ThingStatus{Status: string(ThingStatusOffline)},
			),
			1,
		},
		{
			OnThingReceivedStatusInfoChangedTo("thing", ThingStatusOnline),
			event.TypeThingStatusInfoChanged,
			event.NewThingStatusInfoChangedEvent("thing",
				event.ThingStatus{Status: string(ThingStatusOnline)},
				event.ThingStatus{Status: string(ThingStatusOffline)},
			),
			0,
		},
		{
			OnThingReceivedStatusInfoChangedTo("thing", ThingStatusOnline),
			event.TypeThingStatusInfoChanged,
			event.NewThingStatusInfoChangedEvent("thing",
				event.ThingStatus{Status: string(ThingStatusOffline)},
				event.ThingStatus{Status: string(ThingStatusOnline)},
			),
			1,
		},
	}

	for _, testFixture := range testFixtures {
		const subID = 11

		t.Run("", func(t *testing.T) {
			calls := 0

			run := func(ev event.Event) {
				calls++
			}

			var subscribedCallback func(e event.Event)

			client := newMockSubscriber(t)
			client.On("subscribe", "thing", mock.Anything, mock.Anything).
				Return(func(name string, eventType event.Type, callback func(e event.Event)) int {
					assert.Equal(t, testFixture.eventType, eventType, "incorrect event type")
					subscribedCallback = callback
					return subID
				})

			trigger := testFixture.trigger
			err := trigger.activate(client, run, RuleData{})
			assert.NoError(t, err)
			assert.Equal(t, 0, calls)
			assert.NotNil(t, subscribedCallback)

			// non matching event
			subscribedCallback(testFixture.event)
			assert.Equal(t, testFixture.calls, calls)
		})
	}
}
