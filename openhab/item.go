package openhab

import (
	"context"
	"sync"
	"time"

	"github.com/creativeprojects/gopenhab/api"
)

type Item struct {
	name            string
	data            api.Item
	state           StateValue
	mainType        ItemType
	subType         string
	client          *Client
	stateLocker     sync.Mutex
	stateChanged    chan StateValue
	listenersLocker sync.Mutex
	listeners       int
}

func newItem(client *Client, name string) *Item {
	return &Item{
		name:         name,
		state:        nil,
		client:       client,
		stateChanged: make(chan StateValue),
		listeners:    0,
	}
}

func (i *Item) set(data api.Item) *Item {
	i.data = data
	i.mainType, i.subType = getItemType(i.data.Type)
	i.setInternalState(data.State)
	return i
}

func (i *Item) load() error {
	data := api.Item{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.client.getJSON(ctx, "items/"+i.name, &data)
	if err != nil {
		return err
	}
	i.set(data)
	return nil
}

func (i *Item) hasData() bool {
	return i.data.Name != ""
}

func (i *Item) getData() api.Item {
	if !i.hasData() {
		i.load()
	}
	return i.data
}

func (i *Item) Name() string {
	return i.name
}

func (i *Item) Type() ItemType {
	return i.mainType
}

func (i *Item) stateFromString(state string) StateValue {
	switch i.mainType {
	default:
		return StringState(state)
	case ItemTypeSwitch:
		return SwitchState(state)
	}
}

// State returns an internal cached value if available,
// or calls the api to return a fresh value from openHAB if not
//
// State value is automatically refreshed from openHAB events,
// so you should always get an accurate value.
//
// Please note if you just sent a state change command,
// the new value won't be reflected instantly, but only after openHAB
// sent a state changed event.
//
// If you need the new value after sending a command,
// you might want to use WaitState(duration) instead
func (i *Item) State() (StateValue, error) {
	internalState := i.getInternalStateValue()
	if internalState != nil {
		return internalState, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	state, err := i.client.getString(ctx, "items/"+i.name+"/state")
	if err != nil {
		return i.stateFromString(""), err
	}
	value := i.stateFromString(state)
	i.setInternalStateValue(value)
	return value, nil
}

// WaitState waits until the state gets updated, for a maximum of duration
// If a state was received, it returns it along with true,
// If a state hasn't been received before duration, it returns the last known state along with false
func (i *Item) WaitState(duration time.Duration) (StateValue, bool) {
	i.beginListener()
	defer i.endListener()

	select {
	case state := <-i.stateChanged:
		return state, true
	case <-time.After(duration):
		return i.getInternalStateValue(), false
	}
}

func (i *Item) SendCommand(command StateValue) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.client.postString(ctx, "items/"+i.name, command.String())
	if err != nil {
		return err
	}
	return nil
}

// getInternalStateValue gets the internal state value: it does not trigger an API call to get the state.
func (i *Item) getInternalStateValue() StateValue {
	i.stateLocker.Lock()
	defer i.stateLocker.Unlock()

	return i.state
}

// setInternalStateValue sets the internal state value: it does not trigger an API call to set the state.
// this method should be used after state received or changed events
func (i *Item) setInternalStateValue(state StateValue) {
	i.stateLocker.Lock()
	defer i.stateLocker.Unlock()

	i.state = state
	go i.sendStateChanged(state)
}

func (i *Item) setInternalState(state string) {
	i.setInternalStateValue(i.stateFromString(state))
}

func (i *Item) beginListener() {
	i.listenersLocker.Lock()
	defer i.listenersLocker.Unlock()

	i.listeners++
}

func (i *Item) endListener() {
	i.listenersLocker.Lock()
	defer i.listenersLocker.Unlock()

	i.listeners--
}

func (i *Item) sendStateChanged(value StateValue) {
	i.listenersLocker.Lock()
	defer i.listenersLocker.Unlock()

	// send the message to each receiver
	for count := 0; count < i.listeners; count++ {
		i.stateChanged <- value
	}
}
