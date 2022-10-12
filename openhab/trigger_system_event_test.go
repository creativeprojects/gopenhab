package openhab

import (
	"testing"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
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
