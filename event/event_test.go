package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEventFromJSON(t *testing.T) {
	testData := []struct {
		source string
		event  Event
	}{
		{`{"topic":"smarthome/items/TestSwitch/other","payload":"[{}]","type":"OtherEvent"}`,
			GenericEvent{topic: "smarthome/items/TestSwitch/other", typeName: "OtherEvent", payload: "[{}]"}},

		{`{"topic":"smarthome/items/TestSwitch/command","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemCommandEvent"}`,
			ItemReceivedCommand{topic: "smarthome/items/TestSwitch/command", CommandType: "OnOff", Command: "OFF"}},

		{`{"topic":"smarthome/items/TestSwitch/state","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemStateEvent"}`,
			ItemReceivedState{topic: "smarthome/items/TestSwitch/state", StateType: "OnOff", State: "OFF"}},

		{`{"topic":"smarthome/items/TestSwitch/command","payload":"{\"type\":\"OnOff\",\"value\":\"ON\"}","type":"ItemCommandEvent"}`,
			ItemReceivedCommand{topic: "smarthome/items/TestSwitch/command", CommandType: "OnOff", Command: "ON"}},

		{`{"topic":"smarthome/items/TestSwitch/state","payload":"{\"type\":\"OnOff\",\"value\":\"ON\"}","type":"ItemStateEvent"}`,
			ItemReceivedState{topic: "smarthome/items/TestSwitch/state", StateType: "OnOff", State: "ON"}},

		{`{"topic":"smarthome/items/TestSwitch/statechanged","payload":"{\"type\":\"OnOff\",\"value\":\"ON\",\"oldType\":\"OnOff\",\"oldValue\":\"OFF\"}","type":"ItemStateChangedEvent"}`,
			ItemChanged{topic: "smarthome/items/TestSwitch/statechanged", StateType: "OnOff", State: "ON", OldStateType: "OnOff", OldState: "OFF"}},
	}

	for _, testItem := range testData {
		t.Run("", func(t *testing.T) {
			e, err := New(testItem.source)
			require.NoError(t, err)
			assert.Equal(t, testItem.event, e)
		})
	}
}

func TestErrorEventFromJSON(t *testing.T) {
	testData := []struct {
		source string
	}{
		{`{["topic":"smarthome/items/TestSwitch/other","payload":"","type":"OtherEvent"]}`},
		{`{"topic":"smarthome/items/TestSwitch/other","payload":"","type":"OtherEvent"}]`},
		{`{"topic":"smarthome/items/TestSwitch/command","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}]","type":"ItemCommandEvent"}`},
		{`{"topic":"smarthome/items/TestSwitch/state","payload":"{\"type\":\"OnOff\",\"value\"]:\"OFF\"}","type":"ItemStateEvent"}`},
		{`{"topic":"smarthome/items/TestSwitch/statechanged","payload":"{\"type\":\"OnOff\",[\"value\":\"ON\"}","type":"ItemStateChangedEvent"}`},
	}

	for _, testItem := range testData {
		t.Run("", func(t *testing.T) {
			_, err := New(testItem.source)
			t.Log(err)
			require.Error(t, err)
		})
	}
}
