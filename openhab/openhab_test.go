package openhab

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
	"github.com/creativeprojects/gopenhab/openhabtest"
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

		`event: alive
data: {"type":"ALIVE","interval":10}`,
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

func TestRuleID(t *testing.T) {
	client := NewClient(Config{URL: "http://localhost"})
	id := client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {},
	)
	assert.NotEmpty(t, id)
}

func TestGivenRuleID(t *testing.T) {
	client := NewClient(Config{URL: "http://localhost"})
	id := client.AddRule(
		RuleData{ID: "rule-ID"},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {},
	)
	assert.Equal(t, "rule-ID", id)
}

func TestAddRuleWithError(t *testing.T) {
	client := NewClient(Config{URL: "http://localhost"})
	id := client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {},
		OnTimeCron("0 0 0 ? * * *"), // 7 fields instead of 6
	)
	assert.NotEmpty(t, id)
	client.activateRules()
}

func TestDeleteNoRule(t *testing.T) {
	client := NewClient(Config{URL: "http://localhost"})
	deleted := client.DeleteRule("no rule")
	assert.Equal(t, 0, deleted)
}

func TestDeleteOneRule(t *testing.T) {
	client := NewClient(Config{URL: "http://localhost"})
	id := client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {},
	)
	assert.Equal(t, 1, len(client.rules))
	deleted := client.DeleteRule(id)
	assert.Equal(t, 1, deleted)
	assert.Equal(t, 0, len(client.rules))
}

func TestDeleteTwoRules(t *testing.T) {
	client := NewClient(Config{URL: "http://localhost"})
	id := client.AddRule(
		RuleData{ID: "rule-ID"},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {},
	)
	client.AddRule(
		RuleData{ID: "rule-ID"},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {},
	)
	client.AddRule(
		RuleData{ID: "another rule"},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {},
	)
	assert.Equal(t, 3, len(client.rules))
	deleted := client.DeleteRule(id)
	assert.Equal(t, 2, deleted)
	assert.Equal(t, 1, len(client.rules))
}

func TestStartEvent(t *testing.T) {
	var call int32
	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	defer server.Close()
	client := NewClient(Config{
		URL: server.URL(),
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {
			defer wg.Done()
			atomic.AddInt32(&call, 1)
			ev, ok := e.(event.SystemEvent)
			assert.True(t, ok)
			assert.Equal(t, event.TypeClientStarted, ev.Type())
		},
		OnStart())

	go func() {
		client.Start()
	}()

	wg.Wait()
	client.Stop()
	assert.Equal(t, int32(1), atomic.LoadInt32(&call))
}

func TestConnectEvent(t *testing.T) {
	var call int32
	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	defer server.Close()
	client := NewClient(Config{
		URL: server.URL(),
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {
			defer wg.Done()
			atomic.AddInt32(&call, 1)
			ev, ok := e.(event.SystemEvent)
			assert.True(t, ok)
			assert.Equal(t, event.TypeClientConnected, ev.Type())
		},
		OnConnect())

	go func() {
		client.Start()
	}()

	time.Sleep(10 * time.Millisecond)
	// send a random event so the client detects the connection
	server.Event(event.NewItemReceivedState("item", "string", "state"))

	wg.Wait()
	client.Stop()
	assert.Equal(t, int32(1), atomic.LoadInt32(&call))
}

func TestDisconnectEvent(t *testing.T) {
	var call int32
	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	server.SetItem(api.Item{Name: "item"})

	client := NewClient(Config{
		URL: server.URL(),
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {
			defer wg.Done()
			atomic.AddInt32(&call, 1)
			ev, ok := e.(event.SystemEvent)
			assert.True(t, ok)
			assert.Equal(t, event.TypeClientDisconnected, ev.Type())
		},
		OnDisconnect())

	go func() {
		client.Start()
	}()

	time.Sleep(10 * time.Millisecond)
	// send a random event so the client detects the connection
	server.Event(event.NewItemReceivedState("item", "string", "state"))
	time.Sleep(10 * time.Millisecond)
	// then disconnect the server
	server.Close()

	wg.Wait()
	client.Stop()
	assert.Equal(t, int32(1), atomic.LoadInt32(&call))
}

func TestStopEvent(t *testing.T) {
	var call int32
	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	defer server.Close()
	client := NewClient(Config{
		URL: server.URL(),
	})
	client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {
			ev, ok := e.(event.SystemEvent)
			atomic.AddInt32(&call, 1)
			assert.True(t, ok)
			assert.Equal(t, event.TypeClientStopped, ev.Type())
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
	assert.Equal(t, int32(1), atomic.LoadInt32(&call))
}

func TestErrorEvent(t *testing.T) {
	var call int32
	client := NewClient(Config{URL: "http://localhost", TimeoutHTTP: 100 * time.Millisecond})

	wg := sync.WaitGroup{}
	wg.Add(1)

	// the rule can run twice by the time the client is stopped: make sure we only run it once
	once := sync.Once{}
	client.AddRule(
		RuleData{Name: "Test error rule"},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {
			once.Do(func() {
				defer wg.Done()
				atomic.AddInt32(&call, 1)
				ev, ok := e.(event.ErrorEvent)
				if !ok {
					t.Fatal("expected event to be of type ErrorEvent")
				}
				assert.Equal(t, event.TypeClientError, ev.Type())
				assert.NotEmpty(t, ev.Error())
				t.Log(ev.Error())
			})
		},
		OnError(),
	)

	go func() {
		client.Start()
	}()

	wg.Wait()
	client.Stop()
	assert.Equal(t, int32(1), atomic.LoadInt32(&call))
}

func TestOnDateTimeEvent(t *testing.T) {
	var call int32
	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	defer server.Close()
	client := NewClient(Config{
		URL: server.URL(),
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {
			defer wg.Done()
			atomic.AddInt32(&call, 1)
			ev, ok := e.(event.SystemEvent)
			assert.True(t, ok)
			assert.Equal(t, event.TypeTimeCron, ev.Type())
		},
		OnDateTime(time.Now().Add(time.Second)))

	go func() {
		client.Start()
	}()

	wg.Wait()
	client.Stop()
	assert.Equal(t, int32(1), atomic.LoadInt32(&call))
}

func TestDeleteDateTimeEvent(t *testing.T) {
	client := NewClient(Config{URL: "http://localhost"})
	id := client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {
			t.Error("event shouldn't have been fired")
		},
		OnDateTime(time.Now().Add(time.Second)),
	)
	go func() {
		client.Start()
	}()
	time.Sleep(100 * time.Millisecond)
	deleted := client.DeleteRule(id)
	assert.Equal(t, 1, deleted)

	time.Sleep(2 * time.Second)
	client.Stop()
}

func TestAddRuleOnceStarted(t *testing.T) {
	var call int32
	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	defer server.Close()
	client := NewClient(Config{
		URL: server.URL(),
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		client.Start()
	}()

	time.Sleep(10 * time.Millisecond)
	client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {
			ev, ok := e.(event.SystemEvent)
			assert.True(t, ok)
			assert.Equal(t, event.TypeClientStopped, ev.Type())
			atomic.AddInt32(&call, 1)
		},
		OnStop())

	// stop the server in 20ms
	time.AfterFunc(20*time.Millisecond, func() {
		client.Stop()
	})

	wg.Wait()
	assert.Equal(t, int32(1), atomic.LoadInt32(&call))
}

func TestDispathEventSaveStateFirst(t *testing.T) {
	var call int32
	wg := sync.WaitGroup{}
	server := openhabtest.NewServer(openhabtest.Config{SendEventsFromAPI: true, Log: t})
	defer server.Close()
	err := server.SetItem(api.Item{
		Name:  "item",
		State: "FIRST",
		Type:  "String",
	})
	require.NoError(t, err)

	client := NewClient(Config{URL: server.URL()})
	wg.Add(1)
	client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {
			defer wg.Done()
			atomic.AddInt32(&call, 1)
			ev, ok := e.(event.ItemReceivedState)
			require.True(t, ok)
			item, err := client.GetItem("item")
			require.NoError(t, err)
			state, err := item.State()
			require.NoError(t, err)
			assert.Equal(t, state.String(), ev.State)
		},
		OnItemReceivedState("item", nil),
	)

	wg.Add(1)
	go func() {
		client.Start()
	}()

	time.Sleep(10 * time.Millisecond)
	// Manual test using client's internal method
	client.dispatchRawEvent(`{"topic":"smarthome/items/item/state","payload":"{\"type\":\"String\",\"value\":\"SECOND\"}","type":"ItemStateEvent"}`)
	// send received event from the server
	server.Event(event.NewItemReceivedState("item", "String", "THIRD"))
	wg.Wait()
	client.Stop()
	assert.Equal(t, int32(2), atomic.LoadInt32(&call))
}

func TestCancelEvent(t *testing.T) {
	var call int32
	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	defer server.Close()
	client := NewClient(Config{
		URL: server.URL(),
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	client.AddRule(
		RuleData{},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {
			defer wg.Done()
			<-ctx.Done()
			atomic.AddInt32(&call, 1)

		},
		OnStart())

	go func() {
		client.Start()
	}()

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 1, len(client.rules))
	client.rules[0].cancel()

	wg.Wait()
	client.Stop()
	assert.Equal(t, int32(1), atomic.LoadInt32(&call))
}

func TestTimeoutEvent(t *testing.T) {
	var call int32
	server := openhabtest.NewServer(openhabtest.Config{Log: t})
	defer server.Close()
	client := NewClient(Config{
		URL: server.URL(),
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	client.AddRule(
		RuleData{
			Timeout: 100 * time.Millisecond,
		},
		func(ctx context.Context, client *Client, ruleData RuleData, e event.Event) {
			defer wg.Done()
			<-ctx.Done()
			atomic.AddInt32(&call, 1)

		},
		OnStart())

	go func() {
		client.Start()
	}()

	wg.Wait()
	client.Stop()
	assert.Equal(t, int32(1), atomic.LoadInt32(&call))
}
