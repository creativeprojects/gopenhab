package openhab

import (
	"testing"

	"github.com/creativeprojects/gopenhab/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBaseTriggerSubscribe(t *testing.T) {
	t.Parallel()
	const id = 2
	call := 0

	run := func(ev event.Event) {
		call++
	}

	match := func(ev event.Event) bool {
		return ev.Type() == event.TypeClientConnected
	}

	var subscribedCallback func(e event.Event)

	client := newMockSubscriber(t)
	client.On("subscribe", "thing", event.TypeClientConnected, mock.Anything).
		Return(func(name string, eventType event.Type, callback func(e event.Event)) int {
			subscribedCallback = callback
			return id
		})

	trigger := &baseTrigger{}
	subID := trigger.subscribe(client, "thing", event.TypeClientConnected, run, match)
	assert.Equal(t, id, subID)

	// first one is not the right event
	subscribedCallback(event.NewSystemEvent(event.TypeClientDisconnected))
	assert.Equal(t, 0, call)

	// second one is the right event
	subscribedCallback(event.NewSystemEvent(event.TypeClientConnected))
	assert.Equal(t, 1, call)
}

func TestBaseTriggerSubscribeNilCallback(t *testing.T) {
	t.Parallel()
	const id = 2

	match := func(ev event.Event) bool {
		return ev.Type() == event.TypeClientConnected
	}

	var subscribedCallback func(e event.Event)

	client := newMockSubscriber(t)
	client.On("subscribe", "thing", event.TypeClientConnected, mock.Anything).
		Return(func(name string, eventType event.Type, callback func(e event.Event)) int {
			subscribedCallback = callback
			return id
		})

	trigger := &baseTrigger{}
	subID := trigger.subscribe(client, "thing", event.TypeClientConnected, nil, match)
	assert.Equal(t, id, subID)

	// test that sending events in not throwing a panic
	subscribedCallback(event.NewSystemEvent(event.TypeClientDisconnected))
	subscribedCallback(event.NewSystemEvent(event.TypeClientConnected))
}
