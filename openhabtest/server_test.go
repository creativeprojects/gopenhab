package openhabtest

import (
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanReceiveNotFoundStatus(t *testing.T) {
	urls := []string{
		"/something",
		"/rest/something",
		"/other/fail",
	}
	server := NewServer(nil)
	defer server.Close()

	for _, url := range urls {
		resp, err := http.Get(server.URL() + url)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	}
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
		server.RawEvent(rawEvent)
	}

	// stop the server in 50ms
	time.AfterFunc(50*time.Millisecond, func() {
		server.Close()
	})

	wg.Wait()
}

func TestCanEncodeEvents(t *testing.T) {
	events := []struct {
		e        event.Event
		expected string
	}{
		{
			event.NewItemReceivedCommand("TestSwitch", "OnOff", "ON"),
			`{"topic":"smarthome/items/TestSwitch/command","payload":"{\"type\":\"OnOff\",\"value\":\"ON\"}","type":"ItemCommandEvent"}`,
		},
		{
			event.NewItemReceivedState("TestSwitch", "OnOff", "ON"),
			`{"topic":"smarthome/items/TestSwitch/state","payload":"{\"type\":\"OnOff\",\"value\":\"ON\"}","type":"ItemStateEvent"}`,
		},
		{
			event.NewItemStateChanged("TestSwitch", "OnOff", "OFF", "ON"),
			`{"topic":"smarthome/items/TestSwitch/statechanged","payload":"{\"type\":\"OnOff\",\"value\":\"ON\",\"oldType\":\"OnOff\",\"oldValue\":\"OFF\"}","type":"ItemStateChangedEvent"}`,
		},
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
		for _, ev := range events {
			expected += "event: message\ndata: " + ev.expected + "\n\n"
		}
		assert.Equal(t, expected, string(data))
	}()

	// send some messages
	for _, ev := range events {
		time.Sleep(10 * time.Millisecond)
		server.Event(ev.e)
	}

	// wait a bit before stopping the server
	time.AfterFunc(20*time.Millisecond, func() {
		server.Close()
	})

	wg.Wait()
}
