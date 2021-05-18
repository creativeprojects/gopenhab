package event

import (
	"sync"
)

type PubSub interface {
	Subscribe(name string, eventType Type, callback func(e Event)) int
	Unsubscribe(subId int)
	Publish(e Event)
	Wait()
}

type EventBus struct {
	subs       []subscription
	subLock    sync.Locker
	subIdCount int
	wg         sync.WaitGroup
}

func NewEventBus() *EventBus {
	return &EventBus{
		subs:    make([]subscription, 0),
		subLock: &sync.Mutex{},
	}
}

// Subscribe returns an id for when you need to un-subscribe.
//
// name is the name of the item/thing/channel you want to follow.
// eventType is the type of event you want to follow.
// callback function is called when a matching event occurs.
func (b *EventBus) Subscribe(name string, eventType Type, callback func(e Event)) int {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	b.subIdCount++
	sub := subscription{
		id:        b.subIdCount,
		name:      name,
		eventType: eventType,
		callback:  callback,
	}
	b.subs = append(b.subs, sub)
	return b.subIdCount
}

// Unsubscribe keeps the order of the subscriptions.
// For that reason it is a relatively expensive operation
func (b *EventBus) Unsubscribe(subId int) {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	index := b.findID(subId)
	if index > -1 {
		b.subs = append(b.subs[:index], b.subs[index+1:]...)
	}
}

// findID returns the index in the slice where the sub ID is found,
// it returns -1 if not found
func (b *EventBus) findID(id int) int {
	for index, sub := range b.subs {
		if sub.id == id {
			return index
		}
	}
	return -1
}

// Publish event to all subscribers (in a goroutine each)
func (b *EventBus) Publish(e Event) {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	for _, sub := range b.subs {
		if sub.eventType != e.Type() {
			continue
		}
		if sub.name == "" || sub.eventType.Match(e.Topic(), sub.name) {
			// run the callback in a goroutine
			b.wg.Add(1)
			go func(sub subscription, e Event) {
				defer b.wg.Done()
				sub.callback(e)
			}(sub, e)
		}
	}
}

// Wait for all the subscribers to finish their tasks
func (b *EventBus) Wait() {
	b.wg.Wait()
}

// Verify interface
var _ PubSub = &EventBus{}
