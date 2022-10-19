package openhabtest

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/creativeprojects/gopenhab/api"
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
	server := NewServer(Config{})
	defer server.Close()

	for _, url := range urls {
		resp, err := http.Get(server.URL() + url)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	}
}

func TestCanLoadIndexV2(t *testing.T) {
	server := NewServer(Config{Version: V2})
	defer server.Close()

	resp, err := http.Get(server.URL() + "/rest/")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestCanLoadIndexV3(t *testing.T) {
	server := NewServer(Config{Version: V3})
	defer server.Close()

	resp, err := http.Get(server.URL() + "/rest/")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestCanReceiveRawEvents(t *testing.T) {
	rawEvents := []string{
		`{"topic":"smarthome/things/openweathermap:weather-api:aa/status","payload":"{\"status\":\"ONLINE\",\"statusDetail\":\"NONE\"}","type":"ThingStatusInfoEvent"}`,
		`{"topic":"smarthome/items/LocalWeatherAndForecast_Current_Cloudiness/statechanged","payload":"{\"type\":\"Quantity\",\"value\":\"20 %\",\"oldType\":\"Quantity\",\"oldValue\":\"75 %\"}","type":"ItemStateChangedEvent"}`,
		`{"topic":"smarthome/items/LocalWeatherAndForecast_Current_Cloudiness/state","payload":"{\"type\":\"Quantity\",\"value\":\"20 %\"}","type":"ItemStateEvent"}`,
	}
	server := NewServer(Config{Log: t})
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
		time.Sleep(20 * time.Millisecond)
		server.RawEvent("", rawEvent)
	}

	server.Close()

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
			event.NewItemStateChanged("TestSwitch", "OnOff", "OFF", "OnOff", "ON"),
			`{"topic":"smarthome/items/TestSwitch/statechanged","payload":"{\"type\":\"OnOff\",\"value\":\"ON\",\"oldType\":\"OnOff\",\"oldValue\":\"OFF\"}","type":"ItemStateChangedEvent"}`,
		},
	}
	server := NewServer(Config{Log: t})
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

	server.Close()

	wg.Wait()
}

func TestEmptyItems(t *testing.T) {
	server := NewServer(Config{Log: t})
	defer server.Close()

	resp, err := http.Get(server.URL() + "/rest/items")
	require.NoError(t, err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "[]\n", string(data))
}

func TestListItem(t *testing.T) {
	server := NewServer(Config{Log: t})
	defer server.Close()

	server.SetItem(api.Item{
		Name:  "TestItem",
		Type:  "Switch",
		State: "OFF",
	})

	resp, err := http.Get(server.URL() + "/rest/items")
	require.NoError(t, err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Truef(t, strings.HasPrefix(string(data), `[{"name":"TestItem","label":"","link":"http://`), "unexpected JSON string: %s", string(data))
}

func TestListGroupItem(t *testing.T) {
	server := NewServer(Config{Log: t})
	defer server.Close()

	server.SetItem(api.Item{
		Name:      "TestItem",
		GroupType: "Number",
		Type:      "Group",
		State:     "OFF",
	})

	resp, err := http.Get(server.URL() + "/rest/items")
	require.NoError(t, err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Truef(t, strings.HasPrefix(string(data), `[{"name":"TestItem","label":"","link":"http://`), "unexpected JSON string: %s", string(data))
}

func TestGetItemNotFound(t *testing.T) {
	server := NewServer(Config{Log: t})
	defer server.Close()

	resp, err := http.Get(server.URL() + "/rest/items/NotFound")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetItem(t *testing.T) {
	server := NewServer(Config{Log: t})
	defer server.Close()

	server.SetItem(api.Item{
		Name:  "TestItem",
		Type:  "Switch",
		State: "OFF",
	})

	resp, err := http.Get(server.URL() + "/rest/items/TestItem")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Truef(t, strings.HasPrefix(string(data), `{"name":"TestItem","label":"","link":"http://`), "unexpected JSON string: %s", string(data))
}

func TestGetItemStateNotFound(t *testing.T) {
	server := NewServer(Config{Log: t})
	defer server.Close()

	resp, err := http.Get(server.URL() + "/rest/items/NotFound/state")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetItemState(t *testing.T) {
	state := "20.1 °C"
	server := NewServer(Config{Log: t})
	defer server.Close()

	server.SetItem(api.Item{
		Name:  "TestItem",
		Type:  "Number:Temperature",
		State: state,
	})

	resp, err := http.Get(server.URL() + "/rest/items/TestItem/state")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, state, string(data))
}

func TestSetItemState(t *testing.T) {
	state := "20.1 °C"
	server := NewServer(Config{Log: t})
	defer server.Close()

	server.SetItem(api.Item{
		Name:  "TestItem",
		Type:  "Number:Temperature",
		State: state,
	})

	// set new state
	state = "20.49 °C"
	func() {
		data := bytes.NewBufferString(state)
		req, err := http.NewRequest(http.MethodPut, server.URL()+"/rest/items/TestItem/state", data)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	}()

	resp, err := http.Get(server.URL() + "/rest/items/TestItem/state")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, state, string(data))
}

func TestSendItemCommand(t *testing.T) {
	state := "20.1 °C"
	server := NewServer(Config{Log: t})
	defer server.Close()

	server.SetItem(api.Item{
		Name:  "TestItem",
		Type:  "Number:Temperature",
		State: state,
	})

	// set new state
	state = "20.49 °C"
	func() {
		data := bytes.NewBufferString(state)
		req, err := http.NewRequest(http.MethodPost, server.URL()+"/rest/items/TestItem", data)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}()

	resp, err := http.Get(server.URL() + "/rest/items/TestItem/state")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, state, string(data))
}

func TestSendItemCommandEvents(t *testing.T) {
	state := "20.1 °C"
	server := NewServer(Config{Log: t, SendEventsFromAPI: true})
	defer server.Close()

	server.SetItem(api.Item{
		Name:  "TestItem",
		Type:  "Number:Temperature",
		State: state,
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	// read events
	go func() {
		defer wg.Done()

		resp, err := http.Get(server.URL() + "/rest/events")
		require.NoError(t, err)
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// check all 3 events were sent
		assert.Contains(t, string(data), `{"topic":"smarthome/items/TestItem/command","payload":"{\"type\":\"Test\",\"value\":\"20.49 °C\"}","type":"ItemCommandEvent"}`)
		assert.Contains(t, string(data), `{"topic":"smarthome/items/TestItem/state","payload":"{\"type\":\"Test\",\"value\":\"20.49 °C\"}","type":"ItemStateEvent"}`)
		assert.Contains(t, string(data), `{"topic":"smarthome/items/TestItem/statechanged","payload":"{\"type\":\"Test\",\"value\":\"20.49 °C\",\"oldType\":\"Test\",\"oldValue\":\"20.1 °C\"}","type":"ItemStateChangedEvent"}`)
	}()

	// set new state
	state = "20.49 °C"
	func() {
		time.Sleep(10 * time.Millisecond)

		data := bytes.NewBufferString(state)
		req, err := http.NewRequest(http.MethodPost, server.URL()+"/rest/items/TestItem", data)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}()

	server.Close()
	wg.Wait()
}

func TestSetItemStateEvents(t *testing.T) {
	state := "20.1 °C"
	server := NewServer(Config{Log: t, SendEventsFromAPI: true})
	defer server.Close()

	server.SetItem(api.Item{
		Name:  "TestItem",
		Type:  "Number:Temperature",
		State: state,
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	// read events
	go func() {
		defer wg.Done()

		resp, err := http.Get(server.URL() + "/rest/events")
		require.NoError(t, err)
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// check both events were sent
		assert.Contains(t, string(data), `{"topic":"smarthome/items/TestItem/state","payload":"{\"type\":\"Test\",\"value\":\"20.49 °C\"}","type":"ItemStateEvent"}`)
		assert.Contains(t, string(data), `{"topic":"smarthome/items/TestItem/statechanged","payload":"{\"type\":\"Test\",\"value\":\"20.49 °C\",\"oldType\":\"Test\",\"oldValue\":\"20.1 °C\"}","type":"ItemStateChangedEvent"}`)
	}()

	// set new state
	state = "20.49 °C"
	func() {
		time.Sleep(10 * time.Millisecond)

		data := bytes.NewBufferString(state)
		req, err := http.NewRequest(http.MethodPut, server.URL()+"/rest/items/TestItem/state", data)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	}()

	server.Close()
	wg.Wait()
}
