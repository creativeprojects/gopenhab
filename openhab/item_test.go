package openhab

import (
	"context"
	"strconv"
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
	// don't run parallel (sub-tests are in order)
	item1 := api.Item{
		State:      "OFF",
		Type:       "Switch",
		Name:       "TestSwitch",
		Label:      "Test lights",
		Category:   "lightbulb",
		Tags:       []string{},
		GroupNames: []string{},
	}
	item2 := api.Item{
		State:      "20.2",
		Type:       "Number",
		Name:       "temperature",
		Label:      "House Temperature",
		Category:   "temperature",
		Tags:       []string{},
		GroupNames: []string{},
	}

	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	defer server.Close()

	require.NoError(t, server.SetItem(item1))
	require.NoError(t, server.SetItem(item2))

	client := NewClient(Config{
		URL: server.URL(),
	})

	t.Run("TestLoadItemNotFound", func(t *testing.T) {
		item := newItem(client, "UnknownItem")
		err := item.load(context.Background())
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("TestLoadItem", func(t *testing.T) {
		item := newItem(client, "TestSwitch")
		err := item.load(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "TestSwitch", item.Name())
		assert.Equal(t, ItemTypeSwitch, item.Type())
	})

	t.Run("TestGetItemStateNotFound", func(t *testing.T) {
		item := newItem(client, "UnknownItem")
		_, err := item.State()
		assert.ErrorIs(t, err, ErrNotFound)
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
		item.setInternalState(SwitchON)

		state, err := item.State()
		assert.NoError(t, err)
		assert.Equal(t, SwitchON, state)

		// reset to OFF
		err = item.SendCommand(SwitchOFF)
		assert.NoError(t, err)

		// fake receiving the item state event
		item.setInternalState(SwitchOFF)

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
				wg.Add(1) //nolint:staticcheck
				go func(i int) {
					time.Sleep(time.Duration(i+1) * time.Millisecond)
					ev := event.NewItemReceivedState("TestSwitch", "OnOff", SwitchON.String())
					item.client.userEventBus.Publish(ev)
					wg.Done()
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
		timeout := 100 * time.Millisecond
		item := newTestItem(client, "TestSwitch", "Switch", "OFF")
		state, err := item.State()
		require.NoError(t, err)
		assert.Equal(t, SwitchOFF, state)

		start := time.Now()
		ok, err := item.SendCommandWait(SwitchOFF, timeout)
		require.NoError(t, err)
		assert.False(t, ok)
		// there's a problem with this test
		assert.GreaterOrEqual(t, time.Since(start), timeout)
		t.Log(time.Since(start))
	})

	t.Run("TestMultipleSendCommandWait", func(t *testing.T) {
		count := 10
		initialState := 20.2
		item := newTestItem(client, "temperature", "Number", "20.2")
		state, err := item.State()
		require.NoError(t, err)
		assert.Equal(t, NewDecimalState(initialState, ""), state)

		wg := sync.WaitGroup{}

		// SendCommandWait shouldn't hit the timeout
		timeout := time.AfterFunc(5*time.Second, func() {
			t.Errorf("SendCommandWait is blocked: %+v", item.client.userEventBus.Subscriptions())
		})
		defer timeout.Stop()

		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(i int) {
				newValue := initialState - float64(i)*0.1
				// the event bus is not connected so we send an event manually
				wg.Add(1)
				go func(i int) {
					time.Sleep(time.Duration(i) * time.Millisecond)
					ev := event.NewItemReceivedState("temperature", "Number", strconv.FormatFloat(newValue, 'f', 1, 64))
					item.client.userEventBus.Publish(ev)
					wg.Done()
				}(i)
				_, err := item.SendCommandWait(NewDecimalState(newValue, ""), 10*time.Second)
				assert.NoError(t, err)
				wg.Done()
			}(i)
		}

		wg.Wait()
	})

	assert.NoError(t, server.EventsErr())
	assert.NoError(t, server.ItemsErr())
}

func newTestItem(client *Client, name, itemType, state string) *Item {
	item := newItem(client, name)
	// the item needs a type so it can work properly
	item.set(api.Item{Type: itemType, State: state})
	return item
}

// Unused for now
// func newTestGroupItem(client *Client, name, groupItemType, state string) *Item {
// 	item := newItem(client, name)
// 	// the item needs a type so it can work properly
// 	item.set(api.Item{Type: groupItemType, State: state})
// 	return item
// }
