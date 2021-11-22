package openhab

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/creativeprojects/gopenhab/openhabtest"
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
		t.Logf("not found: %s\n", req.URL)
		http.NotFound(w, req)
	}))

	client := NewClient(Config{
		URL: server.URL,
	})
	err := client.listenEvents()
	assert.NoError(t, err)
}

func TestStartEvent(t *testing.T) {
	called := false
	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	client := NewClient(Config{
		URL: server.URL(),
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	client.AddRule(
		RuleData{},
		func(client RuleClient, ruleData RuleData, e event.Event) {
			defer wg.Done()
			ev, ok := e.(event.SystemEvent)
			assert.True(t, ok)
			assert.Equal(t, event.TypeClientStarted, ev.Type())
			called = true
		},
		OnStart())

	go func() {
		client.Start()
	}()

	wg.Wait()
	client.Stop()
	assert.True(t, called)
}

func TestStopEvent(t *testing.T) {
	called := false
	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	client := NewClient(Config{
		URL: server.URL(),
	})
	client.AddRule(
		RuleData{},
		func(client RuleClient, ruleData RuleData, e event.Event) {
			ev, ok := e.(event.SystemEvent)
			assert.True(t, ok)
			assert.Equal(t, event.TypeClientStopped, ev.Type())
			called = true
		},
		OnStop())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		client.Start()
	}()

	// stop the server in 20ms
	time.AfterFunc(20*time.Millisecond, func() {
		client.Stop()
	})

	wg.Wait()
	assert.True(t, called)
}

func TestErrorEvent(t *testing.T) {
	client := NewClient(Config{URL: "http://localhost", TimeoutHTTP: 100 * time.Millisecond})

	wg := sync.WaitGroup{}
	wg.Add(1)
	client.AddRule(
		RuleData{Name: "Test error rule"},
		func(client RuleClient, ruleData RuleData, e event.Event) {
			defer wg.Done()
			ev, ok := e.(event.SystemEvent)
			if !ok {
				t.Fatal("expected event to be of type SystemEvent")
			}
			assert.Equal(t, event.TypeClientError, ev.Type())
		},
		OnError(),
	)

	go func() {
		client.Start()
	}()

	wg.Wait()
	client.Stop()
}
