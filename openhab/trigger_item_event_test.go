package openhab

import (
	"testing"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestMatchingItemEvent(t *testing.T) {
	testEvents := []struct {
		e       event.Event
		trigger Trigger
		match   bool
	}{
		// received command
		{
			event.NewItemReceivedCommand("TestItem", "OnOff", "ON"),
			OnItemReceivedCommand("TestItem", SwitchOFF),
			false,
		},
		{
			event.NewItemReceivedCommand("TestItem", "OnOff", "ON"),
			OnItemReceivedCommand("TestItem", SwitchON),
			true,
		},
		{
			event.NewItemReceivedCommand("TestItem", "OnOff", "ON"),
			OnItemReceivedCommand("TestItem", nil),
			true,
		},
		{
			event.NewItemReceivedCommand("TestItem", "OnOff", "OFF"),
			OnItemReceivedCommand("TestItem", nil),
			true,
		},
		// received state
		{
			event.NewItemReceivedState("TestItem", "OnOff", "ON"),
			OnItemReceivedState("TestItem", SwitchOFF),
			false,
		},
		{
			event.NewItemReceivedState("TestItem", "OnOff", "ON"),
			OnItemReceivedState("TestItem", SwitchON),
			true,
		},
		{
			event.NewItemReceivedState("TestItem", "OnOff", "ON"),
			OnItemReceivedState("TestItem", nil),
			true,
		},
		{
			event.NewItemReceivedState("TestItem", "OnOff", "OFF"),
			OnItemReceivedState("TestItem", nil),
			true,
		},
		// item state changed
		{
			event.NewItemStateChanged("TestItem", "OnOff", "OFF", "OnOff", "ON"),
			OnItemStateChanged("TestItem"),
			true,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "OFF", "OnOff", "ON"),
			OnItemStateChangedFrom("TestItem", SwitchOFF),
			true,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "ON", "OnOff", "OFF"),
			OnItemStateChangedFrom("TestItem", SwitchOFF),
			false,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "OFF", "OnOff", "ON"),
			OnItemStateChangedTo("TestItem", SwitchON),
			true,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "ON", "OnOff", "OFF"),
			OnItemStateChangedTo("TestItem", SwitchON),
			false,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "OFF", "OnOff", "ON"),
			OnItemStateChangedFromTo("TestItem", SwitchOFF, SwitchON),
			true,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "ON", "OnOff", "OFF"),
			OnItemStateChangedFromTo("TestItem", SwitchOFF, SwitchON),
			false,
		},
		// group item state changed
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "OFF", "OnOff", "ON"),
			OnItemStateChanged("TestItem"),
			true,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "OFF", "OnOff", "ON"),
			OnItemStateChangedFrom("TestItem", SwitchOFF),
			true,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "ON", "OnOff", "OFF"),
			OnItemStateChangedFrom("TestItem", SwitchOFF),
			false,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "OFF", "OnOff", "ON"),
			OnItemStateChangedTo("TestItem", SwitchON),
			true,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "ON", "OnOff", "OFF"),
			OnItemStateChangedTo("TestItem", SwitchON),
			false,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "OFF", "OnOff", "ON"),
			OnItemStateChangedFromTo("TestItem", SwitchOFF, SwitchON),
			true,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "ON", "OnOff", "OFF"),
			OnItemStateChangedFromTo("TestItem", SwitchOFF, SwitchON),
			false,
		},
	}

	for _, testEvent := range testEvents {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, testEvent.match, testEvent.trigger.match(testEvent.e))
		})
	}
}

func TestItemEventSubscription(t *testing.T) {
	testFixtures := []struct {
		trigger    Trigger
		event      event.Event
		eventTypes []event.Type
		calls      int
	}{
		{
			OnItemReceivedCommand("item", SwitchON),
			event.NewItemReceivedCommand("item", "OnOff", string(SwitchOFF)),
			[]event.Type{event.TypeItemCommand},
			0,
		},
		{
			OnItemReceivedCommand("item", SwitchON),
			event.NewItemReceivedCommand("item", "OnOff", string(SwitchON)),
			[]event.Type{event.TypeItemCommand},
			1,
		},
		{
			OnItemReceivedState("item", SwitchON),
			event.NewItemReceivedState("item", "OnOff", string(SwitchOFF)),
			[]event.Type{event.TypeItemState},
			0,
		},
		{
			OnItemReceivedState("item", SwitchON),
			event.NewItemReceivedState("item", "OnOff", string(SwitchON)),
			[]event.Type{event.TypeItemState},
			1,
		},
		{
			OnItemStateChanged("item"),
			event.NewItemStateChanged("item", "OnOff", string(SwitchOFF), "OnOff", string(SwitchON)),
			[]event.Type{event.TypeItemStateChanged, event.TypeGroupItemStateChanged},
			2,
		},
		{
			OnItemStateChangedFrom("item", SwitchOFF),
			event.NewItemStateChanged("item", "OnOff", string(SwitchOFF), "OnOff", string(SwitchON)),
			[]event.Type{event.TypeItemStateChanged, event.TypeGroupItemStateChanged},
			2,
		},
		{
			OnItemStateChangedFrom("item", SwitchON),
			event.NewItemStateChanged("item", "OnOff", string(SwitchOFF), "OnOff", string(SwitchON)),
			[]event.Type{event.TypeItemStateChanged, event.TypeGroupItemStateChanged},
			0,
		},
		{
			OnItemStateChangedTo("item", SwitchON),
			event.NewItemStateChanged("item", "OnOff", string(SwitchOFF), "OnOff", string(SwitchON)),
			[]event.Type{event.TypeItemStateChanged, event.TypeGroupItemStateChanged},
			2,
		},
		{
			OnItemStateChangedTo("item", SwitchOFF),
			event.NewItemStateChanged("item", "OnOff", string(SwitchOFF), "OnOff", string(SwitchON)),
			[]event.Type{event.TypeItemStateChanged, event.TypeGroupItemStateChanged},
			0,
		},
		{
			OnItemStateChangedFromTo("item", SwitchOFF, SwitchON),
			event.NewItemStateChanged("item", "OnOff", string(SwitchOFF), "OnOff", string(SwitchON)),
			[]event.Type{event.TypeItemStateChanged, event.TypeGroupItemStateChanged},
			2,
		},
		{
			OnItemStateChangedFromTo("item", SwitchON, SwitchOFF),
			event.NewItemStateChanged("item", "OnOff", string(SwitchOFF), "OnOff", string(SwitchON)),
			[]event.Type{event.TypeItemStateChanged, event.TypeGroupItemStateChanged},
			0,
		},
	}

	for _, testFixture := range testFixtures {
		const subID = 11

		t.Run("", func(t *testing.T) {
			calls := 0

			run := func(ev event.Event) {
				calls++
			}

			subscribedCallbacks := make([]func(e event.Event), 0, len(testFixture.eventTypes))

			client := newMockSubscriber(t)
			for _, eventType := range testFixture.eventTypes {
				client.On("subscribe", "item", eventType, mock.Anything).
					Return(func(name string, eventType event.Type, callback func(e event.Event)) int {
						subscribedCallbacks = append(subscribedCallbacks, callback)
						return subID
					})
			}

			trigger := testFixture.trigger
			err := trigger.activate(client, run, RuleData{})
			assert.NoError(t, err)
			assert.Equal(t, 0, calls)
			assert.Equal(t, len(testFixture.eventTypes), len(subscribedCallbacks))

			for _, subscribedCallback := range subscribedCallbacks {
				subscribedCallback(testFixture.event)
			}
			assert.Equal(t, testFixture.calls, calls)
		})
	}
}
