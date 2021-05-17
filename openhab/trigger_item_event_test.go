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
	}

	for _, testEvent := range testEvents {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, testEvent.match, testEvent.trigger.match(testEvent.e))
		})
	}
}
