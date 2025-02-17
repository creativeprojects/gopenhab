package openhabtest

import "sync"

type subscription struct {
	id       int
	topic    string
	callback func(message string)
}

// eventBus for the mock server is shamelessly copied from the client EventBus.
// I should probably merge them into a reusable pubsub component at some point.
type eventBus struct {
	subs       []subscription
	subLock    sync.Locker
	subIDCount int
}

func newEventBus() *eventBus {
	return &eventBus{
		subs:    make([]subscription, 0),
		subLock: &sync.Mutex{},
	}
}

// Subscribe returns an id for when you need to un-subscribe.
func (b *eventBus) Subscribe(topic string, callback func(message string)) int {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	b.subIDCount++
	sub := subscription{
		id:       b.subIDCount,
		topic:    topic,
		callback: callback,
	}
	b.subs = append(b.subs, sub)
	return b.subIDCount
}

// Unsubscribe keeps the order of the subscriptions.
// For that reason it is a relatively expensive operation
func (b *eventBus) Unsubscribe(subID int) {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	index := b.findID(subID)
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

// Publish event to all subscribers (sequentially)
func (b *eventBus) Publish(topic, message string) {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	for _, sub := range b.subs {
		if sub.topic == "" || topic == sub.topic {
			sub.callback(message)
		}
	}
}
