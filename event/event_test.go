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
			ItemStateChanged{topic: "smarthome/items/TestSwitch/statechanged", StateType: "OnOff", State: "ON", OldStateType: "OnOff", OldState: "OFF"}},

		{`{"topic":"smarthome/items/TestSwitch/added","payload":"{\"type\":\"Switch\",\"name\":\"TestSwitch\",\"label\":\"TestSwitch\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]}","type":"ItemAddedEvent"}`,
			ItemAdded{topic: "smarthome/items/TestSwitch/added", Item: Item{Name: "TestSwitch", Label: "TestSwitch", Link: "", Type: "Switch", State: "", TransformedState: "", Editable: false, Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""}}},

		{`{"topic":"smarthome/items/TestGroup/added","payload":"{\"groupType\":\"Switch\",\"function\":{\"name\":\"AND\",\"params\":[\"ON\",\"OFF\"]},\"type\":\"Group\",\"name\":\"TestGroup\",\"label\":\"TestGroup\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]}","type":"ItemAddedEvent"}`,
			ItemAdded{topic: "smarthome/items/TestGroup/added", Item: Item{Name: "TestGroup", Label: "TestGroup", Link: "", Type: "Group", State: "", TransformedState: "", Editable: false, Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: "Switch"}}},

		{`{"topic":"smarthome/items/Dummy/removed","payload":"{\"type\":\"Number:Length\",\"name\":\"Dummy\",\"label\":\"hello\",\"category\":\"Light\",\"tags\":[],\"groupNames\":[]}","type":"ItemRemovedEvent"}`,
			ItemRemoved{topic: "smarthome/items/Dummy/removed", Item: Item{Name: "Dummy", Label: "hello", Link: "", Type: "Number:Length", State: "", TransformedState: "", Editable: false, Category: "Light", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""}}},

		{`{"topic":"smarthome/items/TestSwitch/updated","payload":"[{\"type\":\"Switch\",\"name\":\"TestSwitch\",\"label\":\"Test Switch\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]},{\"type\":\"Switch\",\"name\":\"TestSwitch\",\"label\":\"TestSwitch\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]}]","type":"ItemUpdatedEvent"}`,
			ItemUpdated{topic: "smarthome/items/TestSwitch/updated", Item: Item{Name: "TestSwitch", Label: "Test Switch", Link: "", Type: "Switch", State: "", TransformedState: "", Editable: false, Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""}, OldItem: Item{Name: "TestSwitch", Label: "TestSwitch", Link: "", Type: "Switch", State: "", TransformedState: "", Editable: false, Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""}}},
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
