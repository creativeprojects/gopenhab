package openhab

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventMessages(t *testing.T) {
	events := []string{`event: message
data: {"topic":"smarthome/items/TestSwitch/command","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemCommandEvent"}`,

		`event: message
data: {"topic":"smarthome/items/TestSwitch/state","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemStateEvent"}`,

		`event: message
data: {"topic":"smarthome/items/TestSwitch/command","payload":"{\"type\":\"OnOff\",\"value\":\"ON\"}","type":"ItemCommandEvent"}`,

		`event: message
data: {"topic":"smarthome/items/TestSwitch/state","payload":"{\"type\":\"OnOff\",\"value\":\"ON\"}","type":"ItemStateEvent"}`,

		`event: message
data: {"topic":"smarthome/items/TestSwitch/statechanged","payload":"{\"type\":\"OnOff\",\"value\":\"ON\",\"oldType\":\"OnOff\",\"oldValue\":\"OFF\"}","type":"ItemStateChangedEvent"}`,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		if req.Method == http.MethodGet && req.URL.Path == "/rest/events" {
			for _, event := range events {
				w.Write([]byte(event + "\n\n"))
				time.Sleep(100 * time.Millisecond)
			}
			return
		}
		log.Printf("not found: %s\n", req.URL)
		http.NotFound(w, req)
	}))

	client := NewClient(Config{
		URL: server.URL,
	})
	err := client.listenEvents()
	assert.NoError(t, err)
}

func TestLoadEvents(t *testing.T) {
	testData := []struct {
		source string
		event  event.Event
	}{
		{`{"topic":"smarthome/items/TestSwitch/command","payload":"{\"type\":\"OnOff\",\"value\":\"OFF\"}","type":"ItemCommandEvent"}`,
			event.ItemReceivedCommand{CommandType: "OnOff", Command: "OFF"}},
	}

	for _, testItem := range testData {
		t.Run("", func(t *testing.T) {
			e, err := loadEvent(testItem.source)
			require.NoError(t, err)
			assert.Equal(t, testItem.event, e)
		})
	}
}

func TestDispatchEvents(t *testing.T) {
	testData := []struct {
		source string
		event  event.Event
	}{
		{"", event.ItemReceivedCommand{}},
	}

	client := NewClient(Config{URL: "http://localhost:8080"})

	for _, testItem := range testData {
		t.Run("", func(t *testing.T) {
			client.dispatchRawEvent(testItem.source)
		})
	}
}
