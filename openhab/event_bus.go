package openhab

import (
	"sync"

	"github.com/creativeprojects/gopenhab/event"
)

type subscription struct {
	id        int
	topic     string
	eventType event.Type
	callback  func(e event.Event)
}

type eventBus struct {
	subs       []subscription
	subLock    sync.Locker
	subIdCount int
}

func newEventBus() eventBus {
	return eventBus{
		subs:    make([]subscription, 0),
		subLock: &sync.Mutex{},
	}
}

// subscribe returns an id for when you need to un-subscribe
func (b *eventBus) subscribe(topic string, eventType event.Type, callback func(e event.Event)) int {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	b.subIdCount++
	sub := subscription{
		id:        b.subIdCount,
		topic:     topic,
		eventType: eventType,
		callback:  callback,
	}
	b.subs = append(b.subs, sub)
	return b.subIdCount
}

// unsubscribe keeps the order of the subscriptions.
// For that reason it is a relatively expensive operation
func (b *eventBus) unsubscribe(subId int) {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	index := b.findID(subId)
	if index > -1 {
		b.subs = append(b.subs[:index], b.subs[index+1:]...)
	}
}

// findID returns the index in the slice where the sub ID is found,
// it returns -1 if not found
func (b *eventBus) findID(id int) int {
	for index, sub := range b.subs {
		if sub.id == id {
			return index
		}
	}
	return -1
}

// publish event
func (b *eventBus) publish(e event.Event) {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	for _, sub := range b.subs {
		if sub.eventType != e.Type() {
			continue
		}
		// if e.Topic() == "" || strings.HasPrefix(e.Topic(), sub.topic) {
		if e.Topic() == "" || e.Topic() == sub.topic {
			// run the callback in a goroutine
			go func(sub subscription, e event.Event) {
				defer preventPanic()
				sub.callback(e)
			}(sub, e)
		}
	}
}
