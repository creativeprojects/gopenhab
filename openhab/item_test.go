package openhab

import (
	"sync"
	"testing"
	"time"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
	"github.com/creativeprojects/gopenhab/openhabtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemDecimalType(t *testing.T) {
	item := newTestItem(nil, "temperature", "Number", "20.2")
	assert.Equal(t, ItemType("Number"), item.Type())
	assert.Equal(t, DecimalState{20.2, ""}, item.state)
	assert.Equal(t, "20.2", item.state.String())
}

func TestItemDecimalTypeWithUnit(t *testing.T) {
	item := newTestItem(nil, "temperature", "Number:Temperature", "20.2 °C")
	assert.Equal(t, ItemType("Number"), item.Type())
	assert.Equal(t, DecimalState{20.2, "°C"}, item.state)
	assert.Equal(t, "20.2 °C", item.state.String())
}

func TestGetItemAPI(t *testing.T) {

	item := api.Item{
		State:      "OFF",
		Type:       "Switch",
		Name:       "TestSwitch",
		Label:      "Test lights",
		Category:   "lightbulb",
		Tags:       []string{},
		GroupNames: []string{},
	}

	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	defer server.Close()

	server.SetItem(item)

	client := NewClient(Config{
		URL: server.URL(),
	})

	t.Run("TestLoadItemNotFound", func(t *testing.T) {
		item := newItem(client, "UnknownItem")
		err := item.load()
		assert.ErrorIs(t, err, ErrorNotFound)
	})

	t.Run("TestLoadItem", func(t *testing.T) {
		item := newItem(client, "TestSwitch")
		err := item.load()
		assert.NoError(t, err)
		assert.Equal(t, "TestSwitch", item.Name())
		assert.Equal(t, ItemTypeSwitch, item.Type())
	})

	t.Run("TestGetItemStateNotFound", func(t *testing.T) {
		item := newItem(client, "UnknownItem")
		_, err := item.State()
		assert.ErrorIs(t, err, ErrorNotFound)
	})

	t.Run("TestGetItemState", func(t *testing.T) {
		item := newTestItem(client, "TestSwitch", "Switch", "OFF")
		assert.Equal(t, ItemTypeSwitch, item.Type())

		state, err := item.State()
		assert.NoError(t, err)
		assert.Equal(t, SwitchOFF, state)
	})

	t.Run("TestSendCommand", func(t *testing.T) {
		item := newTestItem(client, "TestSwitch", "Switch", "OFF")
		assert.Equal(t, ItemTypeSwitch, item.Type())

		err := item.SendCommand(SwitchON)
		assert.NoError(t, err)

		// fake receiving the item state event
		item.setInternalStateValue(SwitchON)

		state, err := item.State()
		assert.NoError(t, err)
		assert.Equal(t, SwitchON, state)

		// reset to OFF
		err = item.SendCommand(SwitchOFF)
		assert.NoError(t, err)

		// fake receiving the item state event
		item.setInternalStateValue(SwitchOFF)

		state, err = item.State()
		assert.NoError(t, err)
		assert.Equal(t, SwitchOFF, state)
	})

	t.Run("TestListeningForChanges", func(t *testing.T) {
		item := newTestItem(client, "TestSwitch", "Switch", "OFF")
		state, err := item.State()
		require.NoError(t, err)
		assert.Equal(t, SwitchOFF, state)

		wg := sync.WaitGroup{}

		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(i int) {
				// the event bus is not connected so we send an event manually
				go func(i int) {
					time.Sleep(time.Duration(i+1) * time.Millisecond)
					ev := event.NewItemReceivedState("TestSwitch", "OnOff", SwitchON.String())
					item.client.userEventBus.Publish(ev)
				}(i)
				ok, err := item.SendCommandWait(SwitchON, 100*time.Millisecond)
				assert.NoError(t, err)
				assert.True(t, ok)
				wg.Done()
			}(i)
		}

		wg.Wait()
	})

	t.Run("TestSendCommandTimeout", func(t *testing.T) {
		item := newTestItem(client, "TestSwitch", "Switch", "OFF")
		state, err := item.State()
		require.NoError(t, err)
		assert.Equal(t, SwitchOFF, state)

		ok, err := item.SendCommandWait(SwitchOFF, 100*time.Millisecond)
		assert.NoError(t, err)
		assert.False(t, ok)
	})
}

func newTestItem(client *Client, name, itemType, state string) *Item {
	item := newItem(client, name)
	// the item needs a type so it can work properly
	item.set(api.Item{Type: itemType, State: state})
	return item
}

func newTestGroupItem(client *Client, name, groupItemType, state string) *Item {
	item := newItem(client, name)
	// the item needs a type so it can work properly
	item.set(api.Item{Type: groupItemType, State: state})
	return item
}
