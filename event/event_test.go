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
			ItemReceivedCommand{topic: "smarthome/items/TestSwitch/command", ItemName: "TestSwitch", CommandType: "OnOff", Command: "OFF"}},

		{`{"topic":"smarthome/items/TestSwitch/state","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemStateEvent"}`,
			ItemReceivedState{topic: "smarthome/items/TestSwitch/state", ItemName: "TestSwitch", StateType: "OnOff", State: "OFF"}},

		{`{"topic":"smarthome/items/TestSwitch/command","payload":"{\"type\":\"OnOff\",\"value\":\"ON\"}","type":"ItemCommandEvent"}`,
			ItemReceivedCommand{topic: "smarthome/items/TestSwitch/command", ItemName: "TestSwitch", CommandType: "OnOff", Command: "ON"}},

		{`{"topic":"smarthome/items/TestSwitch/state","payload":"{\"type\":\"OnOff\",\"value\":\"ON\"}","type":"ItemStateEvent"}`,
			ItemReceivedState{topic: "smarthome/items/TestSwitch/state", ItemName: "TestSwitch", StateType: "OnOff", State: "ON"}},

		{`{"topic":"smarthome/items/TestSwitch/statechanged","payload":"{\"type\":\"OnOff\",\"value\":\"ON\",\"oldType\":\"OnOff\",\"oldValue\":\"OFF\"}","type":"ItemStateChangedEvent"}`,
			ItemStateChanged{topic: "smarthome/items/TestSwitch/statechanged", ItemName: "TestSwitch", NewStateType: "OnOff", NewState: "ON", PreviousStateType: "OnOff", PreviousState: "OFF"}},

		{`{"topic":"smarthome/items/HouseTemperature/UpstairsTemperature/statechanged","payload":"{\"type\":\"Decimal\",\"value\":\"18.43\",\"oldType\":\"Decimal\",\"oldValue\":\"18.32\"}","type":"GroupItemStateChangedEvent"}`,
			GroupItemStateChanged{topic: "smarthome/items/HouseTemperature/UpstairsTemperature/statechanged", ItemName: "HouseTemperature", TriggeringItem: "UpstairsTemperature", NewStateType: "Decimal", NewState: "18.43", PreviousStateType: "Decimal", PreviousState: "18.32"}},

		{`{"topic":"smarthome/items/TestSwitch/added","payload":"{\"type\":\"Switch\",\"name\":\"TestSwitch\",\"label\":\"TestSwitch\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]}","type":"ItemAddedEvent"}`,
			ItemAdded{topic: "smarthome/items/TestSwitch/added", Item: Item{Name: "TestSwitch", Label: "TestSwitch", Type: "Switch", Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""}}},

		{`{"topic":"smarthome/items/TestGroup/added","payload":"{\"groupType\":\"Switch\",\"function\":{\"name\":\"AND\",\"params\":[\"ON\",\"OFF\"]},\"type\":\"Group\",\"name\":\"TestGroup\",\"label\":\"TestGroup\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]}","type":"ItemAddedEvent"}`,
			ItemAdded{topic: "smarthome/items/TestGroup/added", Item: Item{Name: "TestGroup", Label: "TestGroup", Type: "Group", Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: "Switch"}}},

		{`{"topic":"smarthome/items/Dummy/removed","payload":"{\"type\":\"Number:Length\",\"name\":\"Dummy\",\"label\":\"hello\",\"category\":\"Light\",\"tags\":[],\"groupNames\":[]}","type":"ItemRemovedEvent"}`,
			ItemRemoved{topic: "smarthome/items/Dummy/removed", Item: Item{Name: "Dummy", Label: "hello", Type: "Number:Length", Category: "Light", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""}}},

		{`{"topic":"smarthome/items/TestSwitch/updated","payload":"[{\"type\":\"Switch\",\"name\":\"TestSwitch\",\"label\":\"Test Switch\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]},{\"type\":\"Switch\",\"name\":\"TestSwitch\",\"label\":\"TestSwitch\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]}]","type":"ItemUpdatedEvent"}`,
			ItemUpdated{
				topic:   "smarthome/items/TestSwitch/updated",
				Item:    Item{Name: "TestSwitch", Label: "Test Switch", Type: "Switch", Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""},
				OldItem: Item{Name: "TestSwitch", Label: "TestSwitch", Type: "Switch", Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""},
			},
		},
		// {`{"topic":"smarthome/channels/astro:sun:local:set#event/triggered","payload":"{\"event\":\"START\",\"channel\":\"astro:sun:local:set#event\"}","type":"ChannelTriggeredEvent"}`,
		// },
		// {`{"topic":"smarthome/links/Presence_Mobile_Fred-network:pingdevice:3aadd7c9:online/added","payload":"{\"channelUID\":\"network:pingdevice:3aadd7c9:online\",\"configuration\":{\"profile\":\"system:default\"},\"itemName\":\"Presence_Mobile_Fred\"}","type":"ItemChannelLinkAddedEvent"}`,
		// }
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
		{`{"topic":"smarthome/items//command","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemCommandEvent"}`},
		{`{"topic":"smarthome/items//state","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemStateEvent"}`},
		{`{"topic":"smarthome/items//state","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemStateEvent"}`},
		{`{"topic":"smarthome/items/TestSwitch/state","payload":"{\"type\":\"OnOff\",\"value\"]:\"OFF\"}","type":"ItemStateEvent"}`},
		{`{"topic":"smarthome/items//statechanged","payload":"{\"type\":\"OnOff\",\"value\":\"ON\"}","type":"ItemStateChangedEvent"}`},
		{`{"topic":"smarthome/items/TestSwitch/statechanged","payload":"{\"type\":\"OnOff\",[\"value\":\"ON\"}","type":"ItemStateChangedEvent"}`},
		{`{"topic":"smarthome/items/HouseTemperature/UpstairsTemperature/statechanged","payload":"{\"type\":\"Decimal\",\"value\":\"18.43\",\"oldType\":\"Decimal\",\"oldValue\":\"18.32\"}]","type":"GroupItemStateChangedEvent"}`},
		{`{"topic":"smarthome/items///statechanged","payload":"{\"type\":\"Decimal\",\"value\":\"18.43\",\"oldType\":\"Decimal\",\"oldValue\":\"18.32\"}","type":"GroupItemStateChangedEvent"}`},
		{`{"topic":"smarthome/items/TestGroup/added","payload":"{\"groupType\":\"Switch\",\"function\":{\"name\":\"AND\",\"params\":[\"ON\",\"OFF\"]},\"type\":\"Group\",\"name\":}","type":"ItemAddedEvent"}`},
		{`{"topic":"smarthome/items/Dummy/removed","payload":"{\"type\":\"Number:Length\",\"name\":\"Dummy}","type":"ItemRemovedEvent"}`},
		{`{"topic":"smarthome/items/TestSwitch/updated","payload":"[{\"type\":\"Switch\",\"name\":\"TestSwitch\",]","type":"ItemUpdatedEvent"}`},
		{`{"topic":"smarthome/items/TestSwitch/updated","payload":"[{\"type\":\"Switch\",\"name\":\"TestSwitch\"}]","type":"ItemUpdatedEvent"}`},
	}

	for _, testItem := range testData {
		t.Run("", func(t *testing.T) {
			_, err := New(testItem.source)
			t.Log(err)
			require.Error(t, err)
		})
	}
}
