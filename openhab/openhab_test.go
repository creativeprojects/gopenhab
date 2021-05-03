package openhab

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	err := client.Subscribe("")
	assert.NoError(t, err)
}
