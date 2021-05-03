package openhab

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/stretchr/testify/assert"
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

		state, err := item.State()
		assert.NoError(t, err)
		assert.Equal(t, SwitchON, state)

		// reset to OFF
		err = item.SendCommand(SwitchOFF)
		assert.NoError(t, err)

		state, err = item.State()
		assert.NoError(t, err)
		assert.Equal(t, SwitchOFF, state)
	})
}

func newSwitchItem(client *Client, name string) *Item {
	item := newItem(client, "TestSwitch")
	// the item needs a type so it can work properly
	item.set(api.Item{Type: "Switch"})
	return item
}
