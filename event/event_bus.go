package event

import (
	"sync"
)

type PubSub interface {
	Subscribe(name string, eventType Type, callback func(e Event)) int
	SubscribeOnce(name string, eventType Type, callback func(e Event)) int
	Unsubscribe(subID int) int
	Publish(e Event) int
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
	return b.subscribe(name, eventType, false, callback)
}

// SubscribeOnce can only receive one event.
//
// name is the name of the item/thing/channel you want to follow.
// eventType is the type of event you want to follow.
// callback function is called when a matching event occurs.
func (b *eventBus) SubscribeOnce(name string, eventType Type, callback func(e Event)) int {
	return b.subscribe(name, eventType, true, callback)
}

func (b *eventBus) subscribe(name string, eventType Type, once bool, callback func(e Event)) int {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	b.subIdCount++
	sub := subscription{
		id:        b.subIdCount,
		name:      name,
		eventType: eventType,
		callback:  callback,
		once:      once,
	}
	b.subs = append(b.subs, sub)
	return b.subIdCount
}

// Unsubscribe keeps the order of the subscriptions.
// For that reason it is a relatively expensive operation
// It returns the number of subscriptions removed
func (b *eventBus) Unsubscribe(subID int) int {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	return b.unsubscribe(subID)
}

// unsubscribe is not thread safe, it should be called from within a locked context
func (b *eventBus) unsubscribe(subID int) int {
	found := 0
	if index := b.findID(subID); index > -1 {
		found = 1
		b.subs = append(b.subs[:index], b.subs[index+1:]...)
	}
	return found
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
// It returns the number of subscribers that received the event
func (b *eventBus) Publish(event Event) int {
	b.subLock.Lock()
	defer b.subLock.Unlock()

	unsubscribed := make([]int, 0)
	receivers := 0

	for _, sub := range b.subs {
		if sub.eventType != event.Type() {
			continue
		}
		if sub.name == "" || sub.eventType.Match(event.Topic(), sub.name) {
			if sub.once {
				unsubscribed = append(unsubscribed, sub.id)
			}
			receivers++
			if b.async {
				// run the callback in a goroutine
				b.wg.Add(1)
				go func(b *eventBus, sub subscription, e Event) {
					defer b.wg.Done()
					sub.callback(e)
				}(b, sub, event)
			} else {
				// run synchronously
				sub.callback(event)
			}
		}
	}

	// remove subscriptions that are only valid once
	for _, id := range unsubscribed {
		b.unsubscribe(id)
	}
	return receivers
}

// Wait for all the subscribers to finish their tasks
func (b *eventBus) Wait() {
	b.wg.Wait()
}

// Verify interface
var _ PubSub = &eventBus{}
