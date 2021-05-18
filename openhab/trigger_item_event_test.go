package openhab

import (
	"testing"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
)

func TestMatchingEvent(t *testing.T) {
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
			event.NewItemStateChanged("TestItem", "OnOff", "OFF", "ON"),
			OnItemStateChanged("TestItem"),
			true,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "OFF", "ON"),
			OnItemStateChangedFrom("TestItem", SwitchOFF),
			true,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "ON", "OFF"),
			OnItemStateChangedFrom("TestItem", SwitchOFF),
			false,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "OFF", "ON"),
			OnItemStateChangedTo("TestItem", SwitchON),
			true,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "ON", "OFF"),
			OnItemStateChangedTo("TestItem", SwitchON),
			false,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "OFF", "ON"),
			OnItemStateChangedFromTo("TestItem", SwitchOFF, SwitchON),
			true,
		},
		{
			event.NewItemStateChanged("TestItem", "OnOff", "ON", "OFF"),
			OnItemStateChangedFromTo("TestItem", SwitchOFF, SwitchON),
			false,
		},
		// group item state changed
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "OFF", "ON"),
			OnItemStateChanged("TestItem"),
			true,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "OFF", "ON"),
			OnItemStateChangedFrom("TestItem", SwitchOFF),
			true,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "ON", "OFF"),
			OnItemStateChangedFrom("TestItem", SwitchOFF),
			false,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "OFF", "ON"),
			OnItemStateChangedTo("TestItem", SwitchON),
			true,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "ON", "OFF"),
			OnItemStateChangedTo("TestItem", SwitchON),
			false,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "OFF", "ON"),
			OnItemStateChangedFromTo("TestItem", SwitchOFF, SwitchON),
			true,
		},
		{
			event.NewGroupItemStateChanged("TestItem", "TriggeringItem", "OnOff", "ON", "OFF"),
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
