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

type eventBus struct {
	async      bool
	subs       []subscription
	subLock    sync.Locker
	subIdCount int
	wg         sync.WaitGroup
}

func NewEventBus(async bool) *eventBus {
	return &eventBus{
		async:   async,
		subs:    make([]subscription, 0),
		subLock: &sync.Mutex{},
	}
}

// Subscribe returns an id for when you need to un-subscribe.
//
// name is the name of the item/thing/channel you want to follow.
// eventType is the type of event you want to follow.
// callback function is called when a matching event occurs.
func (b *eventBus) Subscribe(name string, eventType Type, callback func(e Event)) int {
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
func (b *eventBus) Unsubscribe(subId int) {
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

// Publish event to all subscribers (in a goroutine each)
func (b *eventBus) Publish(e Event) {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	for _, sub := range b.subs {
		if sub.eventType != e.Type() {
			continue
		}
		if sub.name == "" || sub.eventType.Match(e.Topic(), sub.name) {
			if b.async {
				// run the callback in a goroutine
				b.wg.Add(1)
				go func(b *eventBus, sub subscription, e Event) {
					defer b.wg.Done()
					sub.callback(e)
				}(b, sub, e)
			} else {
				// run synchronously
				sub.callback(e)
			}
		}
	}
}

// Wait for all the subscribers to finish their tasks
func (b *eventBus) Wait() {
	b.wg.Wait()
}

// Verify interface
var _ PubSub = &eventBus{}
