package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOneSubscriberNoTopic(t *testing.T) {
	done := make(chan Event)
	eventBus := NewEventBus()

	eventBus.Subscribe("", TypeClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	eventBus.Subscribe("", TypeClientConnected, func(e Event) {
		done <- e
	})

	eventBus.Publish(newMockEvent("", TypeClientConnected))

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

	eventBus.Subscribe("topic", TypeItemState, func(e Event) {
		panic("This function should not have been called!")
	})
	eventBus.Subscribe("topic2", TypeItemState, func(e Event) {
		done <- e
	})

	eventBus.Publish(newMockEvent("topic2", TypeItemState))

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
	done1 := make(chan Event)
	done2 := make(chan Event)
	eventBus := NewEventBus()

	eventBus.Subscribe("topic", TypeItemState, func(e Event) {
		if first {
			t.Error("function called more than once")
		}
		done1 <- e
	})
	eventBus.Subscribe("", TypeItemState, func(e Event) {
		if second {
			t.Error("function called more than once")
		}
		done2 <- e
	})

	eventBus.Publish(newMockEvent("topic", TypeItemState))

	for {
		select {
		case e := <-done1:
			assert.IsType(t, mockEvent{}, e)
			if second {
				return
			}
			first = true

		case e := <-done2:
			assert.IsType(t, mockEvent{}, e)
			if first {
				return
			}
			second = true

		case <-time.After(100 * time.Millisecond):
			// fail the test after 100ms
			t.Fatal("timeout!")
		}
	}
}

func TestUnsubscribe(t *testing.T) {
	done := make(chan Event)
	eventBus := NewEventBus()

	eventBus.Subscribe("", TypeClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	sub := eventBus.Subscribe("", TypeClientConnected, func(e Event) {
		done <- e
	})

	eventBus.Publish(newMockEvent("", TypeClientConnected))

	select {
	case e := <-done:
		assert.IsType(t, mockEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}

	eventBus.Unsubscribe(sub)

	// republish an event
	eventBus.Publish(newMockEvent("", TypeClientConnected))

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

	eventBus.Subscribe("", TypeClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	sub := eventBus.Subscribe("", TypeClientConnected, func(e Event) {
		done <- e
	})

	eventBus.Publish(newMockEvent("", TypeClientConnected))

	select {
	case e := <-done:
		assert.IsType(t, mockEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}

	eventBus.Unsubscribe(sub + 10)

	// republish an event
	eventBus.Publish(newMockEvent("", TypeClientConnected))

	select {
	case e := <-done:
		assert.IsType(t, mockEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}
}
