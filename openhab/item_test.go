package openhab

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type itemModel struct {
	Link       string   `json:"link"`
	State      string   `json:"state"`
	Editable   bool     `json:"editable"`
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Label      string   `json:"label"`
	Category   string   `json:"category"`
	Tags       []string `json:"tags"`
	GroupNames []string `json:"groupNames"`
}

func TestGetItemAPI(t *testing.T) {
	response := itemModel{
		Link:       "http://openhab:8080/rest/items/TestSwitch",
		State:      "OFF",
		Editable:   false,
		Type:       "Switch",
		Name:       "TestSwitch",
		Label:      "Test lights",
		Category:   "lightbulb",
		Tags:       []string{},
		GroupNames: []string{},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		if req.Method == http.MethodGet {
			if req.URL.Path == "/rest/items/TestSwitch" {
				encoder := json.NewEncoder(w)
				err := encoder.Encode(response)
				assert.NoError(t, err)
				return
			}
			if req.URL.Path == "/rest/items/TestSwitch/state" {
				w.Write([]byte(response.State))
				return
			}
		}
		if req.Method == http.MethodPost {
			if req.URL.Path == "/rest/items/TestSwitch" {
				value, err := io.ReadAll(req.Body)
				assert.NoError(t, err)
				response.State = string(value)
				return
			}
		}
		http.NotFound(w, req)
	}))

	client := NewClient(Config{
		URL: server.URL,
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
		item := newSwitchItem(client, "TestSwitch")
		assert.Equal(t, ItemTypeSwitch, item.Type())
		state, err := item.State()
		assert.NoError(t, err)
		assert.Equal(t, SwitchOFF, state)
	})

	t.Run("TestSendCommand", func(t *testing.T) {
		item := newSwitchItem(client, "TestSwitch")
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
		item := newSwitchItem(client, "TestSwitch")
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
					ev := event.NewItemReceivedState("smarthome/items/TestSwitch/state")
					ev.State = SwitchON.String()
					item.client.eventBus.Publish(ev)
				}(i)
				ok, err := item.SendCommandWait(SwitchON, 100*time.Millisecond)
				assert.NoError(t, err)
				assert.True(t, ok)
				wg.Done()
			}(i)
		}

		wg.Wait()
	})

	// Test is FAILING
	t.Run("TestSendCommandTimeout", func(t *testing.T) {
		item := newSwitchItem(client, "TestSwitch")
		state, err := item.State()
		require.NoError(t, err)
		assert.Equal(t, SwitchOFF, state)

		ok, err := item.SendCommandWait(SwitchON, 100*time.Millisecond)
		assert.NoError(t, err)
		assert.False(t, ok)
	})
}

func newSwitchItem(client *Client, name string) *Item {
	item := newItem(client, name)
	// the item needs a type so it can work properly
	item.set(api.Item{Type: "Switch", State: "OFF"})
	return item
}
