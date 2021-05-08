package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOneSubscriberNoTopic(t *testing.T) {
	done := make(chan Event)
	eventBus := NewEventBus()

	eventBus.Subscribe("", ClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	eventBus.Subscribe("", ClientConnected, func(e Event) {
		done <- e
	})

	eventBus.Publish(newMockEvent("", ClientConnected))

	select {
	case e := <-done:
		assert.IsType(t, mockEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}
}

func TestOneSubscriberWithTopic(t *testing.T) {
	done := make(chan Event, 1)
	eventBus := NewEventBus()

	eventBus.Subscribe("topic", ItemState, func(e Event) {
		panic("This function should not have been called!")
	})
	eventBus.Subscribe("topic2", ItemState, func(e Event) {
		done <- e
	})

	eventBus.Publish(newMockEvent("topic2", ItemState))

	select {
	case e := <-done:
		assert.IsType(t, mockEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}
}

func TestTwoSubscribers(t *testing.T) {
	first := false
	second := false
	done := make(chan Event, 2)
	eventBus := NewEventBus()

	eventBus.Subscribe("topic", ItemState, func(e Event) {
		if first {
			t.Error("function called more than once")
		}
		first = true
		done <- e
	})
	eventBus.Subscribe("", ItemState, func(e Event) {
		if second {
			t.Error("function called more than once")
		}
		second = true
		done <- e
	})

	eventBus.Publish(newMockEvent("topic", ItemState))

	for {
		select {
		case e := <-done:
			assert.IsType(t, mockEvent{}, e)
			if first && second {
				return
			}
		case <-time.After(100 * time.Millisecond):
			// fail the test after 100ms
			t.Fatal("timeout!")
		}
	}
}

func TestUnsubscribe(t *testing.T) {
	done := make(chan Event)
	eventBus := NewEventBus()

	eventBus.Subscribe("", ClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	sub := eventBus.Subscribe("", ClientConnected, func(e Event) {
		done <- e
	})

	eventBus.Publish(newMockEvent("", ClientConnected))

	select {
	case e := <-done:
		assert.IsType(t, mockEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}

	eventBus.Unsubscribe(sub)

	// republish an event
	eventBus.Publish(newMockEvent("", ClientConnected))

	select {
	case <-done:
		t.Fatal("subscription should have been cancelled")
	case <-time.After(100 * time.Millisecond):
		// success
	}
}

func TestUnsubscribeUnknownID(t *testing.T) {
	done := make(chan Event)
	eventBus := NewEventBus()

	eventBus.Subscribe("", ClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	sub := eventBus.Subscribe("", ClientConnected, func(e Event) {
		done <- e
	})

	eventBus.Publish(newMockEvent("", ClientConnected))

	select {
	case e := <-done:
		assert.IsType(t, mockEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}

	eventBus.Unsubscribe(sub + 10)

	// republish an event
	eventBus.Publish(newMockEvent("", ClientConnected))

	select {
	case e := <-done:
		assert.IsType(t, mockEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}
}
