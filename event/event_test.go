package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEventFromJSON(t *testing.T) {
	t.Parallel()
	testData := []struct {
		source string
		event  Event
	}{
		{
			`{"topic":"smarthome/items/TestSwitch/other","payload":"[{}]","type":"OtherEvent"}`,
			GenericEvent{topic: "items/TestSwitch/other", typeName: "OtherEvent", payload: "[{}]"},
		},

		{
			`{"topic":"openhab/items/TestSwitch/other","payload":"[{}]","type":"OtherEvent"}`,
			GenericEvent{topic: "items/TestSwitch/other", typeName: "OtherEvent", payload: "[{}]"},
		},

		{
			`{"topic":"smarthome/items/TestSwitch/command","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemCommandEvent"}`,
			ItemReceivedCommand{topic: "items/TestSwitch/command", ItemName: "TestSwitch", CommandType: "OnOff", Command: "OFF"},
		},

		{
			`{"topic":"smarthome/items/TestSwitch/state","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemStateEvent"}`,
			ItemReceivedState{topic: "items/TestSwitch/state", ItemName: "TestSwitch", StateType: "OnOff", State: "OFF"},
		},

		{
			`{"topic":"smarthome/items/TestSwitch/command","payload":"{\"type\":\"OnOff\",\"value\":\"ON\"}","type":"ItemCommandEvent"}`,
			ItemReceivedCommand{topic: "items/TestSwitch/command", ItemName: "TestSwitch", CommandType: "OnOff", Command: "ON"},
		},

		{
			`{"topic":"smarthome/items/TestSwitch/state","payload":"{\"type\":\"OnOff\",\"value\":\"ON\"}","type":"ItemStateEvent"}`,
			ItemReceivedState{topic: "items/TestSwitch/state", ItemName: "TestSwitch", StateType: "OnOff", State: "ON"},
		},

		{
			`{"topic":"openhab/items/TestTemperature/stateupdated","payload":"{\"type\":\"Quantity\",\"value\":\"20.0 °C\"}","type":"ItemStateUpdatedEvent"}`,
			ItemStateUpdated{topic: "items/TestTemperature/stateupdated", ItemName: "TestTemperature", StateType: "Quantity", State: "20.0 °C"},
		},

		{
			`{"topic":"smarthome/items/TestSwitch/statechanged","payload":"{\"type\":\"OnOff\",\"value\":\"ON\",\"oldType\":\"OnOff\",\"oldValue\":\"OFF\"}","type":"ItemStateChangedEvent"}`,
			ItemStateChanged{topic: "items/TestSwitch/statechanged", ItemName: "TestSwitch", NewStateType: "OnOff", NewState: "ON", PreviousStateType: "OnOff", PreviousState: "OFF"},
		},

		{
			`{"topic":"openhab/items/ChristmasLightsGroup/Back_Garden_Lighting_Switch/stateupdated","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"GroupStateUpdatedEvent"}`,
			GroupItemStateUpdated{topic: "items/ChristmasLightsGroup/Back_Garden_Lighting_Switch/stateupdated", ItemName: "ChristmasLightsGroup", TriggeringItem: "Back_Garden_Lighting_Switch", StateType: "OnOff", State: "OFF"},
		},

		{
			`{"topic":"smarthome/items/HouseTemperature/UpstairsTemperature/statechanged","payload":"{\"type\":\"Decimal\",\"value\":\"18.43\",\"oldType\":\"Decimal\",\"oldValue\":\"18.32\"}","type":"GroupItemStateChangedEvent"}`,
			GroupItemStateChanged{topic: "items/HouseTemperature/UpstairsTemperature/statechanged", ItemName: "HouseTemperature", TriggeringItem: "UpstairsTemperature", NewStateType: "Decimal", NewState: "18.43", PreviousStateType: "Decimal", PreviousState: "18.32"},
		},

		{
			`{"topic":"smarthome/items/TestSwitch/added","payload":"{\"type\":\"Switch\",\"name\":\"TestSwitch\",\"label\":\"TestSwitch\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]}","type":"ItemAddedEvent"}`,
			ItemAdded{topic: "items/TestSwitch/added", Item: Item{Name: "TestSwitch", Label: "TestSwitch", Type: "Switch", Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""}},
		},

		{
			`{"topic":"smarthome/items/TestGroup/added","payload":"{\"groupType\":\"Switch\",\"function\":{\"name\":\"AND\",\"params\":[\"ON\",\"OFF\"]},\"type\":\"Group\",\"name\":\"TestGroup\",\"label\":\"TestGroup\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]}","type":"ItemAddedEvent"}`,
			ItemAdded{topic: "items/TestGroup/added", Item: Item{Name: "TestGroup", Label: "TestGroup", Type: "Group", Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: "Switch"}},
		},

		{
			`{"topic":"smarthome/items/Dummy/removed","payload":"{\"type\":\"Number:Length\",\"name\":\"Dummy\",\"label\":\"hello\",\"category\":\"Light\",\"tags\":[],\"groupNames\":[]}","type":"ItemRemovedEvent"}`,
			ItemRemoved{topic: "items/Dummy/removed", Item: Item{Name: "Dummy", Label: "hello", Type: "Number:Length", Category: "Light", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""}},
		},

		{
			`{"topic":"smarthome/items/TestSwitch/updated","payload":"[{\"type\":\"Switch\",\"name\":\"TestSwitch\",\"label\":\"Test Switch\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]},{\"type\":\"Switch\",\"name\":\"TestSwitch\",\"label\":\"TestSwitch\",\"category\":\"Switch\",\"tags\":[],\"groupNames\":[]}]","type":"ItemUpdatedEvent"}`,
			ItemUpdated{
				topic:   "items/TestSwitch/updated",
				Item:    Item{Name: "TestSwitch", Label: "Test Switch", Type: "Switch", Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""},
				OldItem: Item{Name: "TestSwitch", Label: "TestSwitch", Type: "Switch", Category: "Switch", Tags: []string{}, GroupNames: []string{}, Members: []string(nil), GroupType: ""},
			},
		},
		{
			`{"topic":"openhab/things/mqtt:homie300:6a75cc6119:test/status","payload":"{\"status\":\"ONLINE\",\"statusDetail\":\"DUTY_CYCLE\"}","type":"ThingStatusInfoEvent"}`,
			ThingStatusInfoEvent{topic: "things/mqtt:homie300:6a75cc6119:test/status", ThingName: "mqtt:homie300:6a75cc6119:test", Status: "ONLINE", StatusDetail: "DUTY_CYCLE"},
		},
		{
			`{"topic":"openhab/things/mqtt:homie300:6a75cc6119:test/statuschanged","payload":"[{\"status\":\"ONLINE\",\"statusDetail\":\"CONFIGURATION_PENDING\"},{\"status\":\"OFFLINE\",\"statusDetail\":\"COMMUNICATION_ERROR\",\"description\":\"Did not receive all required topics\"}]","type":"ThingStatusInfoChangedEvent"}`,
			ThingStatusInfoChangedEvent{topic: "things/mqtt:homie300:6a75cc6119:test/statuschanged", ThingName: "mqtt:homie300:6a75cc6119:test", PreviousStatus: "OFFLINE", PreviousStatusDetail: "COMMUNICATION_ERROR", PreviousDescription: "Did not receive all required topics", NewStatus: "ONLINE", NewStatusDetail: "CONFIGURATION_PENDING"},
		},
		{
			`{"topic":"openhab/things/mqtt:homie300:6a75cc6119:test/statuschanged","payload":"[{\"status\":\"ONLINE\",\"statusDetail\":\"DUTY_CYCLE\"},{\"status\":\"ONLINE\",\"statusDetail\":\"CONFIGURATION_PENDING\"}]","type":"ThingStatusInfoChangedEvent"}`,
			ThingStatusInfoChangedEvent{topic: "things/mqtt:homie300:6a75cc6119:test/statuschanged", ThingName: "mqtt:homie300:6a75cc6119:test", PreviousStatus: "ONLINE", PreviousStatusDetail: "CONFIGURATION_PENDING", NewStatus: "ONLINE", NewStatusDetail: "DUTY_CYCLE"},
		},
		{
			`{"type":"ALIVE"}`,
			AliveEvent{},
		},
		{
			`{"topic":"openhab/system/startlevel","payload":"{\"startlevel\":30}","type":"StartlevelEvent"}`,
			StartlevelEvent{topic: "system/startlevel", level: 30},
		},
		{
			`{"topic":"smarthome/channels/astro:sun:local:set#event/triggered","payload":"{\"event\":\"START\",\"channel\":\"astro:sun:local:set#event\"}","type":"ChannelTriggeredEvent"}`,
			ChannelTriggered{topic: "channels/astro:sun:local:set#event/triggered", ChannelName: "astro:sun:local:set#event", Event: "START"},
		},
		{
			`{
				"topic": "smarthome/things/zwave:device:c4dcc784:node8/updated",
				"payload": "[{\"label\":\"RaZberry 2 controller\",\"bridgeUID\":\"zwave:serial_zstick:c4dcc784\",\"configuration\":{\"node_id\":8},\"properties\":{\"zwave_class_basic\":\"BASIC_TYPE_STATIC_CONTROLLER\",\"zwave_class_generic\":\"GENERIC_TYPE_STATIC_CONTROLLER\",\"zwave_frequent\":\"false\",\"zwave_neighbours\":\"1,6\",\"zwave_version\":\"0.0\",\"zwave_listening\":\"true\",\"zwave_plus_devicetype\":\"NODE_TYPE_ZWAVEPLUS_NODE\",\"zwave_nodeid\":\"8\",\"zwave_lastheal\":\"2022-12-10T01:27:54Z\",\"zwave_routing\":\"true\",\"zwave_plus_roletype\":\"ROLE_TYPE_CONTROLLER_CENTRAL_STATIC\",\"zwave_beaming\":\"true\",\"zwave_secure\":\"false\",\"zwave_class_specific\":\"SPECIFIC_TYPE_GATEWAY\"},\"UID\":\"zwave:device:c4dcc784:node8\",\"thingTypeUID\":\"zwave:device\",\"channels\":[]},{\"label\":\"RaZberry 2 controller\",\"bridgeUID\":\"zwave:serial_zstick:c4dcc784\",\"configuration\":{\"node_id\":8},\"properties\":{\"zwave_class_basic\":\"BASIC_TYPE_STATIC_CONTROLLER\",\"zwave_class_generic\":\"GENERIC_TYPE_STATIC_CONTROLLER\",\"zwave_frequent\":\"false\",\"zwave_neighbours\":\"1,6\",\"zwave_version\":\"0.0\",\"zwave_listening\":\"true\",\"zwave_plus_devicetype\":\"NODE_TYPE_ZWAVEPLUS_NODE\",\"zwave_nodeid\":\"8\",\"zwave_lastheal\":\"2022-12-09T01:27:42Z\",\"zwave_routing\":\"true\",\"zwave_plus_roletype\":\"ROLE_TYPE_CONTROLLER_CENTRAL_STATIC\",\"zwave_beaming\":\"true\",\"zwave_secure\":\"false\",\"zwave_class_specific\":\"SPECIFIC_TYPE_GATEWAY\"},\"UID\":\"zwave:device:c4dcc784:node8\",\"thingTypeUID\":\"zwave:device\",\"channels\":[]}]",
				"type": "ThingUpdatedEvent"
			  }`,
			ThingUpdated{
				topic:    "things/zwave:device:c4dcc784:node8/updated",
				OldThing: Thing{UID: "zwave:device:c4dcc784:node8", Label: "RaZberry 2 controller", BridgeUID: "zwave:serial_zstick:c4dcc784", Configuration: map[string]interface{}{"node_id": float64(8)}, Properties: map[string]string{"zwave_beaming": "true", "zwave_class_basic": "BASIC_TYPE_STATIC_CONTROLLER", "zwave_class_generic": "GENERIC_TYPE_STATIC_CONTROLLER", "zwave_class_specific": "SPECIFIC_TYPE_GATEWAY", "zwave_frequent": "false", "zwave_lastheal": "2022-12-09T01:27:42Z", "zwave_listening": "true", "zwave_neighbours": "1,6", "zwave_nodeid": "8", "zwave_plus_devicetype": "NODE_TYPE_ZWAVEPLUS_NODE", "zwave_plus_roletype": "ROLE_TYPE_CONTROLLER_CENTRAL_STATIC", "zwave_routing": "true", "zwave_secure": "false", "zwave_version": "0.0"}, ThingTypeUID: "zwave:device"},
				Thing:    Thing{UID: "zwave:device:c4dcc784:node8", Label: "RaZberry 2 controller", BridgeUID: "zwave:serial_zstick:c4dcc784", Configuration: map[string]interface{}{"node_id": float64(8)}, Properties: map[string]string{"zwave_beaming": "true", "zwave_class_basic": "BASIC_TYPE_STATIC_CONTROLLER", "zwave_class_generic": "GENERIC_TYPE_STATIC_CONTROLLER", "zwave_class_specific": "SPECIFIC_TYPE_GATEWAY", "zwave_frequent": "false", "zwave_lastheal": "2022-12-10T01:27:54Z", "zwave_listening": "true", "zwave_neighbours": "1,6", "zwave_nodeid": "8", "zwave_plus_devicetype": "NODE_TYPE_ZWAVEPLUS_NODE", "zwave_plus_roletype": "ROLE_TYPE_CONTROLLER_CENTRAL_STATIC", "zwave_routing": "true", "zwave_secure": "false", "zwave_version": "0.0"}, ThingTypeUID: "zwave:device"},
			},
		},
		// {`{"topic":"smarthome/links/Presence_Mobile_Fred-network:pingdevice:3aadd7c9:online/added","payload":"{\"channelUID\":\"network:pingdevice:3aadd7c9:online\",\"configuration\":{\"profile\":\"system:default\"},\"itemName\":\"Presence_Mobile_Fred\"}","type":"ItemChannelLinkAddedEvent"}`,
		// }
		// {
		// 	`"[{\"channels\":[{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#humidity\",\"id\":\"sensors#humidity\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fhumidity\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Humidity\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Humidity\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"%\",\"datatype\":\"float_\"}},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#luminance\",\"id\":\"sensors#luminance\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fluminance\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Luminance\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Luminance\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"float_\"}},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#motion\",\"id\":\"sensors#motion\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fmotion\",\"itemType\":\"Switch\",\"kind\":\"STATE\",\"label\":\"Motion\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Motion\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"boolean_\"}},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#temperature\",\"id\":\"sensors#temperature\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Ftemperature\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Temperature\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Temperature\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"°C\",\"datatype\":\"float_\"}}],\"label\":\"multisensor1\",\"bridgeUID\":\"mqtt:broker:6a75cc6119\",\"configuration\":{\"deviceid\":\"multisensor1\",\"removetopics\":false,\"basetopic\":\"homie\"},\"properties\":{\"homieversion\":\"4.0.0\"},\"UID\":\"mqtt:homie300:6a75cc6119:multisensor1\",\"thingTypeUID\":\"mqtt:homie300\"},{\"channels\":[{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#humidity\",\"id\":\"sensors#humidity\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fhumidity\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Humidity\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Humidity\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"%\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#luminance\",\"id\":\"sensors#luminance\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fluminance\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Luminance\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Luminance\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#motion\",\"id\":\"sensors#motion\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fmotion\",\"itemType\":\"Switch\",\"kind\":\"STATE\",\"label\":\"Motion\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Motion\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"boolean_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#temperature\",\"id\":\"sensors#temperature\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Ftemperature\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Temperature\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Temperature\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"°C\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"}],\"label\":\"multisensor1\",\"bridgeUID\":\"mqtt:broker:6a75cc6119\",\"configuration\":{\"deviceid\":\"multisensor1\",\"removetopics\":false,\"basetopic\":\"homie\"},\"properties\":{\"homieversion\":\"4.0.0\"},\"UID\":\"mqtt:homie300:6a75cc6119:multisensor1\",\"thingTypeUID\":\"mqtt:homie300\"}]" ({"topic":"openhab/things/mqtt:homie300:6a75cc6119:multisensor1/updated","payload":"[{\"channels\":[{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#humidity\",\"id\":\"sensors#humidity\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fhumidity\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Humidity\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Humidity\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"%\",\"datatype\":\"float_\"}},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#luminance\",\"id\":\"sensors#luminance\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fluminance\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Luminance\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Luminance\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"float_\"}},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#motion\",\"id\":\"sensors#motion\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fmotion\",\"itemType\":\"Switch\",\"kind\":\"STATE\",\"label\":\"Motion\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Motion\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"boolean_\"}},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#temperature\",\"id\":\"sensors#temperature\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Ftemperature\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Temperature\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Temperature\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"°C\",\"datatype\":\"float_\"}}],\"label\":\"multisensor1\",\"bridgeUID\":\"mqtt:broker:6a75cc6119\",\"configuration\":{\"deviceid\":\"multisensor1\",\"removetopics\":false,\"basetopic\":\"homie\"},\"properties\":{\"homieversion\":\"4.0.0\"},\"UID\":\"mqtt:homie300:6a75cc6119:multisensor1\",\"thingTypeUID\":\"mqtt:homie300\"},{\"channels\":[{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#humidity\",\"id\":\"sensors#humidity\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fhumidity\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Humidity\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Humidity\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"%\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#luminance\",\"id\":\"sensors#luminance\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fluminance\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Luminance\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Luminance\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#motion\",\"id\":\"sensors#motion\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fmotion\",\"itemType\":\"Switch\",\"kind\":\"STATE\",\"label\":\"Motion\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Motion\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"boolean_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#temperature\",\"id\":\"sensors#temperature\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Ftemperature\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Temperature\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Temperature\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"°C\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"}],\"label\":\"multisensor1\",\"bridgeUID\":\"mqtt:broker:6a75cc6119\",\"configuration\":{\"deviceid\":\"multisensor1\",\"removetopics\":false,\"basetopic\":\"homie\"},\"properties\":{\"homieversion\":\"4.0.0\"},\"UID\":\"mqtt:homie300:6a75cc6119:multisensor1\",\"thingTypeUID\":\"mqtt:homie300\"}]","type":"ThingUpdatedEvent"})`,
		// },
		// {
		// 	`"[{\"channels\":[{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#humidity\",\"id\":\"sensors#humidity\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fhumidity\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Humidity\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Humidity\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"%\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#luminance\",\"id\":\"sensors#luminance\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fluminance\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Luminance\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Luminance\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#motion\",\"id\":\"sensors#motion\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fmotion\",\"itemType\":\"Switch\",\"kind\":\"STATE\",\"label\":\"Motion\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Motion\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"boolean_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#temperature\",\"id\":\"sensors#temperature\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Ftemperature\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Temperature\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Temperature\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"°C\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"}],\"label\":\"multisensor1\",\"bridgeUID\":\"mqtt:broker:6a75cc6119\",\"configuration\":{\"deviceid\":\"multisensor1\",\"removetopics\":false,\"basetopic\":\"homie\"},\"properties\":{\"homieversion\":\"4.0.0\"},\"UID\":\"mqtt:homie300:6a75cc6119:multisensor1\",\"thingTypeUID\":\"mqtt:homie300\"},{\"channels\":[{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#humidity\",\"id\":\"sensors#humidity\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fhumidity\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Humidity\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Humidity\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"%\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#luminance\",\"id\":\"sensors#luminance\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fluminance\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Luminance\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Luminance\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#motion\",\"id\":\"sensors#motion\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fmotion\",\"itemType\":\"Switch\",\"kind\":\"STATE\",\"label\":\"Motion\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Motion\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"boolean_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#temperature\",\"id\":\"sensors#temperature\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Ftemperature\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Temperature\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Temperature\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"°C\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"}],\"label\":\"multisensor1\",\"bridgeUID\":\"mqtt:broker:6a75cc6119\",\"configuration\":{\"deviceid\":\"multisensor1\",\"removetopics\":false,\"basetopic\":\"homie\"},\"properties\":{\"homieversion\":\"4.0.0\"},\"UID\":\"mqtt:homie300:6a75cc6119:multisensor1\",\"thingTypeUID\":\"mqtt:homie300\"}]" ({"topic":"openhab/things/mqtt:homie300:6a75cc6119:multisensor1/updated","payload":"[{\"channels\":[{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#humidity\",\"id\":\"sensors#humidity\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fhumidity\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Humidity\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Humidity\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"%\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#luminance\",\"id\":\"sensors#luminance\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fluminance\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Luminance\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Luminance\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#motion\",\"id\":\"sensors#motion\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fmotion\",\"itemType\":\"Switch\",\"kind\":\"STATE\",\"label\":\"Motion\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Motion\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"boolean_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#temperature\",\"id\":\"sensors#temperature\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Ftemperature\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Temperature\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Temperature\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"°C\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"}],\"label\":\"multisensor1\",\"bridgeUID\":\"mqtt:broker:6a75cc6119\",\"configuration\":{\"deviceid\":\"multisensor1\",\"removetopics\":false,\"basetopic\":\"homie\"},\"properties\":{\"homieversion\":\"4.0.0\"},\"UID\":\"mqtt:homie300:6a75cc6119:multisensor1\",\"thingTypeUID\":\"mqtt:homie300\"},{\"channels\":[{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#humidity\",\"id\":\"sensors#humidity\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fhumidity\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Humidity\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Humidity\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"%\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#luminance\",\"id\":\"sensors#luminance\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fluminance\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Luminance\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Luminance\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#motion\",\"id\":\"sensors#motion\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Fmotion\",\"itemType\":\"Switch\",\"kind\":\"STATE\",\"label\":\"Motion\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Motion\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"\",\"datatype\":\"boolean_\"},\"autoUpdatePolicy\":\"DEFAULT\"},{\"uid\":\"mqtt:homie300:6a75cc6119:multisensor1:sensors#temperature\",\"id\":\"sensors#temperature\",\"channelTypeUID\":\"mqtt:homie_2Fmultisensor1_2Fsensors_2Ftemperature\",\"itemType\":\"Number\",\"kind\":\"STATE\",\"label\":\"Temperature\",\"defaultTags\":[],\"properties\":{},\"configuration\":{\"format\":\"\",\"name\":\"Temperature\",\"retained\":\"true\",\"settable\":\"false\",\"unit\":\"°C\",\"datatype\":\"float_\"},\"autoUpdatePolicy\":\"DEFAULT\"}],\"label\":\"multisensor1\",\"bridgeUID\":\"mqtt:broker:6a75cc6119\",\"configuration\":{\"deviceid\":\"multisensor1\",\"removetopics\":false,\"basetopic\":\"homie\"},\"properties\":{\"homieversion\":\"4.0.0\"},\"UID\":\"mqtt:homie300:6a75cc6119:multisensor1\",\"thingTypeUID\":\"mqtt:homie300\"}]","type":"ThingUpdatedEvent"})`,
		// },
		// {
		// 	`{\"channelUID\":\"mqtt:homie300:6a75cc6119:feathers2n2:bme280#humidity\",\"configuration\":{},\"itemName\":\"MQTT_FeatherS2_N2_Humidity\"}" ({"topic":"openhab/links/MQTT_FeatherS2_N2_Humidity-mqtt:homie300:6a75cc6119:feathers2n2:bme280#humidity/removed","payload":"{\"channelUID\":\"mqtt:homie300:6a75cc6119:feathers2n2:bme280#humidity\",\"configuration\":{},\"itemName\":\"MQTT_FeatherS2_N2_Humidity\"}","type":"ItemChannelLinkRemovedEvent"})
		// 	"{\"channelUID\":\"mqtt:homie300:6a75cc6119:envirophat:bmp280#temperature\",\"configuration\":{},\"itemName\":\"MQTTEnvirophatAgent_Bmp280_Temperature\"}" ({"topic":"openhab/links/MQTTEnvirophatAgent_Bmp280_Temperature-mqtt:homie300:6a75cc6119:envirophat:bmp280#temperature/removed","payload":"{\"channelUID\":\"mqtt:homie300:6a75cc6119:envirophat:bmp280#temperature\",\"configuration\":{},\"itemName\":\"MQTTEnvirophatAgent_Bmp280_Temperature\"}","type":"ItemChannelLinkRemovedEvent"}`,
		// },
		// {
		// 	`{"topic":"openhab/addons/misc-homekit/installed","payload":"\"misc-homekit\"","type":"AddonEvent"}`,
		// },
	}

	for _, testItem := range testData {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			e, err := New(testItem.source)
			require.NoError(t, err)
			assert.Equal(t, testItem.event, e)
		})
	}
}

func TestErrorEventFromJSON(t *testing.T) {
	t.Parallel()
	testData := []struct {
		source string
	}{
		{`{["topic":"smarthome/items/TestSwitch/other","payload":"","type":"OtherEvent"]}`},
		{`{"topic":"smarthome/items/TestSwitch/other","payload":"","type":"OtherEvent"}]`},
		{`{"topic":"smarthome/items//command","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemCommandEvent"}`},
		{`{"topic":"smarthome/items/TestSwitch/command","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemCommandEvent"}]`},
		{`{"topic":"smarthome/items//state","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemStateEvent"}`},
		{`{"topic":"smarthome/items//state","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemStateEvent"}`},
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
			t.Parallel()
			_, err := New(testItem.source)
			t.Log(err)
			require.Error(t, err)
		})
	}
}
