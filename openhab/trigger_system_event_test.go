package openhab

import (
	"testing"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestMatchingSystemEvent(t *testing.T) {
	testEvents := []struct {
		e       event.Event
		trigger Trigger
		match   bool
	}{
		{
			event.NewSystemEvent(event.TypeClientConnectionStable),
			OnConnect(),
			false,
		},
		{
			event.NewSystemEvent(event.TypeClientConnected),
			OnConnect(),
			true,
		},
		{
			event.NewSystemEvent(event.TypeClientConnected),
			OnStableConnection(),
			false,
		},
		{
			event.NewSystemEvent(event.TypeClientConnectionStable),
			OnStableConnection(),
			true,
		},
		{
			event.NewSystemEvent(event.TypeClientConnected),
			OnDisconnect(),
			false,
		},
		{
			event.NewSystemEvent(event.TypeClientDisconnected),
			OnDisconnect(),
			true,
		},
		{
			event.NewSystemEvent(event.TypeClientStopped),
			OnStart(),
			false,
		},
		{
			event.NewSystemEvent(event.TypeClientStarted),
			OnStart(),
			true,
		},
		{
			event.NewSystemEvent(event.TypeClientStarted),
			OnStop(),
			false,
		},
		{
			event.NewSystemEvent(event.TypeClientStopped),
			OnStop(),
			true,
		},
		{
			event.NewSystemEvent(event.TypeClientStarted),
			OnError(),
			false,
		},
		{
			event.NewErrorEvent(nil),
			OnError(),
			true,
		},
		{
			event.NewSystemEvent(event.TypeClientStarted),
			OnAlive(),
			false,
		},
		{
			event.NewAliveEvent(),
			OnAlive(),
			true,
		},
	}

	for _, testEvent := range testEvents {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, testEvent.match, testEvent.trigger.match(testEvent.e))
		})
	}
}

func TestSystemEventSubscription(t *testing.T) {
	testFixtures := []struct {
		trigger   Trigger
		eventType event.Type
		event     event.Event
		calls     int
	}{
		{
			OnConnect(),
			event.TypeClientConnected,
			event.NewSystemEvent(event.TypeClientDisconnected),
			0,
		},
		{
			OnConnect(),
			event.TypeClientConnected,
			event.NewSystemEvent(event.TypeClientConnected),
			1,
		},
		{
			OnStableConnection(),
			event.TypeClientConnectionStable,
			event.NewSystemEvent(event.TypeClientDisconnected),
			0,
		},
		{
			OnStableConnection(),
			event.TypeClientConnectionStable,
			event.NewSystemEvent(event.TypeClientConnectionStable),
			1,
		},
		{
			OnDisconnect(),
			event.TypeClientDisconnected,
			event.NewSystemEvent(event.TypeClientConnected),
			0,
		},
		{
			OnDisconnect(),
			event.TypeClientDisconnected,
			event.NewSystemEvent(event.TypeClientDisconnected),
			1,
		},
		{
			OnStart(),
			event.TypeClientStarted,
			event.NewSystemEvent(event.TypeClientStopped),
			0,
		},
		{
			OnStart(),
			event.TypeClientStarted,
			event.NewSystemEvent(event.TypeClientStarted),
			1,
		},
		{
			OnStop(),
			event.TypeClientStopped,
			event.NewSystemEvent(event.TypeClientStarted),
			0,
		},
		{
			OnStop(),
			event.TypeClientStopped,
			event.NewSystemEvent(event.TypeClientStopped),
			1,
		},
		{
			OnError(),
			event.TypeClientError,
			event.NewSystemEvent(event.TypeClientStarted),
			0,
		},
		{
			OnError(),
			event.TypeClientError,
			event.NewSystemEvent(event.TypeClientError),
			1,
		},
		{
			OnAlive(),
			event.TypeServerAlive,
			event.NewSystemEvent(event.TypeClientStarted),
			0,
		},
		{
			OnAlive(),
			event.TypeServerAlive,
			event.NewSystemEvent(event.TypeServerAlive),
			1,
		},
	}

	for _, testFixture := range testFixtures {
		const subID = 11

		t.Run("", func(t *testing.T) {
			calls := 0

			run := func(ev event.Event) {
				calls++
			}

			var subscribedCallback func(e event.Event)

			client := newMockSubscriber(t)
			client.On("subscribe", "", mock.Anything, mock.Anything).
				Return(func(name string, eventType event.Type, callback func(e event.Event)) int {
					assert.Equal(t, testFixture.eventType, eventType, "incorrect event type")
					subscribedCallback = callback
					return subID
				})

			trigger := testFixture.trigger
			err := trigger.activate(client, run, RuleData{})
			assert.NoError(t, err)
			assert.Equal(t, 0, calls)
			assert.NotNil(t, subscribedCallback)

			subscribedCallback(testFixture.event)
			assert.Equal(t, testFixture.calls, calls)
		})
	}
}