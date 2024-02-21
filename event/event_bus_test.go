package event

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOneSubscriberNoTopic(t *testing.T) {
	call := false
	eventBus := NewEventBus(false)

	eventBus.Subscribe("", TypeClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	eventBus.Subscribe("", TypeClientConnected, func(e Event) {
		assert.IsType(t, fakeEvent{}, e)
		call = true
	})

	eventBus.Publish(newFakeEvent("", TypeClientConnected))
	assert.True(t, call)
	eventBus.Wait()
}

func TestOneSubscriberNoTopicAsync(t *testing.T) {
	done := make(chan Event)
	eventBus := NewEventBus(true)

	eventBus.Subscribe("", TypeClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	eventBus.Subscribe("", TypeClientConnected, func(e Event) {
		done <- e
	})

	eventBus.Publish(newFakeEvent("", TypeClientConnected))

	select {
	case e := <-done:
		assert.IsType(t, fakeEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}
	eventBus.Wait()
}

func TestOneSubscriberWithTopic(t *testing.T) {
	call := false
	eventBus := NewEventBus(false)

	eventBus.Subscribe("item", TypeItemState, func(e Event) {
		panic("This function should not have been called!")
	})
	eventBus.Subscribe("item2", TypeItemState, func(e Event) {
		call = true
	})

	eventBus.Publish(newFakeEvent("items/item2/state", TypeItemState))
	assert.True(t, call)
	eventBus.Wait()
}

func TestOneSubscriberWithTopicAsync(t *testing.T) {
	done := make(chan Event, 1)
	eventBus := NewEventBus(true)

	eventBus.Subscribe("item", TypeItemState, func(e Event) {
		panic("This function should not have been called!")
	})
	eventBus.Subscribe("item2", TypeItemState, func(e Event) {
		done <- e
	})

	eventBus.Publish(newFakeEvent("items/item2/state", TypeItemState))

	select {
	case e := <-done:
		assert.IsType(t, fakeEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}
	eventBus.Wait()
}

func TestTwoSubscribers(t *testing.T) {
	first := false
	second := false
	eventBus := NewEventBus(false)

	eventBus.Subscribe("item", TypeItemState, func(e Event) {
		assert.IsType(t, fakeEvent{}, e)
		if first {
			t.Error("function called more than once")
		}
		first = true
	})
	eventBus.Subscribe("item", TypeItemState, func(e Event) {
		assert.IsType(t, fakeEvent{}, e)
		if second {
			t.Error("function called more than once")
		}
		second = true
	})

	eventBus.Publish(newFakeEvent("items/item/state", TypeItemState))
	assert.True(t, first)
	assert.True(t, second)
}

func TestTwoSubscribersAsync(t *testing.T) {
	first := false
	second := false
	done1 := make(chan Event)
	done2 := make(chan Event)
	eventBus := NewEventBus(true)

	eventBus.Subscribe("item", TypeItemState, func(e Event) {
		if first {
			t.Error("function called more than once")
		}
		done1 <- e
	})
	eventBus.Subscribe("item", TypeItemState, func(e Event) {
		if second {
			t.Error("function called more than once")
		}
		done2 <- e
	})

	eventBus.Publish(newFakeEvent("items/item/state", TypeItemState))

	for {
		select {
		case e := <-done1:
			assert.IsType(t, fakeEvent{}, e)
			if second {
				return
			}
			first = true

		case e := <-done2:
			assert.IsType(t, fakeEvent{}, e)
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
	call := 0
	eventBus := NewEventBus(false)

	eventBus.Subscribe("", TypeClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	sub := eventBus.Subscribe("", TypeClientConnected, func(e Event) {
		call++
	})

	eventBus.Publish(newFakeEvent("", TypeClientConnected))
	assert.Equal(t, 1, call)
	eventBus.Wait()

	eventBus.Unsubscribe(sub)

	// republish an event
	eventBus.Publish(newFakeEvent("", TypeClientConnected))
	assert.Equal(t, 1, call)
	eventBus.Wait()
}

func TestUnsubscribeAsync(t *testing.T) {
	done := make(chan Event)
	eventBus := NewEventBus(true)

	eventBus.Subscribe("", TypeClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	sub := eventBus.Subscribe("", TypeClientConnected, func(e Event) {
		done <- e
	})

	eventBus.Publish(newFakeEvent("", TypeClientConnected))

	select {
	case e := <-done:
		assert.IsType(t, fakeEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}
	eventBus.Wait()

	eventBus.Unsubscribe(sub)

	// republish an event
	eventBus.Publish(newFakeEvent("", TypeClientConnected))

	select {
	case <-done:
		t.Fatal("subscription should have been cancelled")
	case <-time.After(100 * time.Millisecond):
		// success
	}
	eventBus.Wait()
}

func TestUnsubscribeUnknownID(t *testing.T) {
	call := 0
	eventBus := NewEventBus(false)

	eventBus.Subscribe("", TypeClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	sub := eventBus.Subscribe("", TypeClientConnected, func(e Event) {
		assert.IsType(t, fakeEvent{}, e)
		call++
	})

	eventBus.Publish(newFakeEvent("", TypeClientConnected))
	assert.Equal(t, 1, call)
	eventBus.Wait()

	eventBus.Unsubscribe(sub + 10)

	// republish an event
	eventBus.Publish(newFakeEvent("", TypeClientConnected))
	assert.Equal(t, 2, call)
	eventBus.Wait()
}

func TestUnsubscribeUnknownIDAsync(t *testing.T) {
	done := make(chan Event)
	eventBus := NewEventBus(true)

	eventBus.Subscribe("", TypeClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	sub := eventBus.Subscribe("", TypeClientConnected, func(e Event) {
		done <- e
	})

	eventBus.Publish(newFakeEvent("", TypeClientConnected))

	select {
	case e := <-done:
		assert.IsType(t, fakeEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}
	eventBus.Wait()

	eventBus.Unsubscribe(sub + 10)

	// republish an event
	eventBus.Publish(newFakeEvent("", TypeClientConnected))

	select {
	case e := <-done:
		assert.IsType(t, fakeEvent{}, e)
	case <-time.After(100 * time.Millisecond):
		// fail the test after 100ms
		t.Fatal("timeout!")
	}
	eventBus.Wait()
}

func TestUnsubscribeWhileRunningAsync(t *testing.T) {
	eventBus := NewEventBus(true)
	sub := eventBus.Subscribe("", TypeClientConnected, func(e Event) {
		time.Sleep(100 * time.Millisecond)
	})
	// publish an event
	published := eventBus.Publish(newFakeEvent("", TypeClientConnected))
	assert.Equal(t, 1, published)

	unsub := eventBus.Unsubscribe(sub)
	assert.Equal(t, 1, unsub)

	eventBus.Wait()
}

func TestSubscribeOnce(t *testing.T) {
	call := 0
	eventBus := NewEventBus(false)

	eventBus.SubscribeOnce("", TypeClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	eventBus.SubscribeOnce("", TypeClientConnected, func(e Event) {
		call++
	})
	eventBus.SubscribeOnce("", TypeClientConnected, func(e Event) {
		call++
	})

	published := eventBus.Publish(newFakeEvent("", TypeClientConnected))
	assert.Equal(t, 2, published)

	published = eventBus.Publish(newFakeEvent("", TypeClientConnected))
	assert.Equal(t, 0, published)

	eventBus.Wait()
	assert.Equal(t, 2, call)
}

func TestSubscribeOnceAsync(t *testing.T) {
	var call int32
	eventBus := NewEventBus(true)

	eventBus.SubscribeOnce("", TypeClientDisconnected, func(e Event) {
		panic("This function should not have been called!")
	})
	eventBus.SubscribeOnce("", TypeClientConnected, func(e Event) {
		atomic.AddInt32(&call, 1)
	})
	eventBus.SubscribeOnce("", TypeClientConnected, func(e Event) {
		atomic.AddInt32(&call, 1)
	})

	published := eventBus.Publish(newFakeEvent("", TypeClientConnected))
	assert.Equal(t, 2, published)

	published = eventBus.Publish(newFakeEvent("", TypeClientConnected))
	assert.Equal(t, 0, published)

	eventBus.Wait()
	assert.Equal(t, int32(2), atomic.LoadInt32(&call))
}
