package openhabtest

import (
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanReceiveNotFoundStatus(t *testing.T) {
	server := NewServer(nil)
	defer server.Close()

	resp, err := http.Get(server.URL() + "/something")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCanReceiveNotFoundRestStatus(t *testing.T) {
	server := NewServer(nil)
	defer server.Close()

	resp, err := http.Get(server.URL() + "/rest/something")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCanReceiveRawEvents(t *testing.T) {
	rawEvents := []string{
		`{"topic":"smarthome/things/openweathermap:weather-api:aa/status","payload":"{\"status\":\"ONLINE\",\"statusDetail\":\"NONE\"}","type":"ThingStatusInfoEvent"}`,
		`{"topic":"smarthome/items/LocalWeatherAndForecast_Current_Cloudiness/statechanged","payload":"{\"type\":\"Quantity\",\"value\":\"20 %\",\"oldType\":\"Quantity\",\"oldValue\":\"75 %\"}","type":"ItemStateChangedEvent"}`,
		`{"topic":"smarthome/items/LocalWeatherAndForecast_Current_Cloudiness/state","payload":"{\"type\":\"Quantity\",\"value\":\"20 %\"}","type":"ItemStateEvent"}`,
	}
	server := NewServer(t)
	defer server.Close()

	wg := sync.WaitGroup{}

	// request and read from the client
	wg.Add(1)
	go func() {
		defer wg.Done()

		resp, err := http.Get(server.URL() + "/rest/events")
		require.NoError(t, err)
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		expected := ""
		// build expected data
		for _, rawEvent := range rawEvents {
			expected += "event: message\ndata: " + rawEvent + "\n\n"
		}
		assert.Equal(t, expected, string(data))
	}()

	// send some messages
	for _, rawEvent := range rawEvents {
		time.Sleep(10 * time.Millisecond)
		server.SendRawEvent(rawEvent)
	}

	// stop the server in 50ms
	time.AfterFunc(50*time.Millisecond, func() {
		server.Close()
	})

	wg.Wait()
}
