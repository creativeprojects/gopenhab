package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitItemTopic(t *testing.T) {
	t.Parallel()
	testData := []struct {
		topic          string
		item           string
		triggeringItem string
		eventType      string
	}{
		{"something...", "", "", ""},
		{"something/not/too/bad", "", "", ""},
		{"smarthome/items/item", "", "", ""},
		{"smarthome/things/item/event", "", "", ""},
		{"smarthome/items/item/event", "item", "", "event"},
		{"smarthome/items/item/trigger/event", "item", "trigger", "event"},
		{"smarthome/items/item/trigger/event/toomany", "", "", ""},
		{"openhab/items/item/event", "item", "", "event"},
		{"openhab/items/item/trigger/event", "item", "trigger", "event"},
	}

	for _, testItem := range testData {
		t.Run(testItem.topic, func(t *testing.T) {
			t.Parallel()
			item, triggeringItem, eventType := splitItemTopic(testItem.topic)
			assert.Equal(t, testItem.item, item)
			assert.Equal(t, testItem.triggeringItem, triggeringItem)
			assert.Equal(t, testItem.eventType, eventType)
		})
	}
}
