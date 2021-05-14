package openhab

import (
	"context"
	"sync"
	"time"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
)

// Item represents an item in openHAB
type Item struct {
	name        string
	data        api.Item
	state       StateValue
	mainType    ItemType
	subType     string
	client      *Client
	apiLocker   sync.Mutex
	stateLocker sync.Mutex
}

func newItem(client *Client, name string) *Item {
	return &Item{
		name:   name,
		state:  nil,
		client: client,
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
	ctx, cancel := context.WithTimeout(context.Background(), i.client.config.TimeoutHTTP)
	defer cancel()

	err := i.client.getJSON(ctx, "items/"+i.name, &data)
	if err != nil {
		return err
	}
	i.set(data)
	return nil
}

// func (i *Item) hasData() bool {
// 	return i.data.Name != ""
// }

// func (i *Item) getData() api.Item {
// 	if !i.hasData() {
// 		i.load()
// 	}
// 	return i.data
// }

// Name returns the name of the item (an item name is unique in openHAB)
func (i *Item) Name() string {
	return i.name
}

// Type return the item type
func (i *Item) Type() ItemType {
	return i.mainType
}

// State returns an internal cached value if available,
// or calls the api to return a fresh value from openHAB if not
//
// State value is automatically refreshed from openHAB events,
// so you should always get an accurate value.
//
// Please note if you just sent a state change command,
// the new value might not be reflected instantly,
// but only after openHAB sent a state changed event.
func (i *Item) State() (StateValue, error) {
	internalState := i.getInternalStateValue()
	if internalState != nil {
		return internalState, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), i.client.config.TimeoutHTTP)
	defer cancel()

	state, err := i.client.getString(ctx, "items/"+i.name+"/state")
	if err != nil {
		return i.stateFromString(""), err
	}
	value := i.stateFromString(state)
	i.setInternalStateValue(value)
	return value, nil
}

// SendCommand sends a command to an item
func (i *Item) SendCommand(command StateValue) error {
	i.apiLocker.Lock()
	defer i.apiLocker.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), i.client.config.TimeoutHTTP)
	defer cancel()

	err := i.client.postString(ctx, "items/"+i.name, command.String())
	if err != nil {
		return err
	}
	return nil
}

// SendCommandWait sends a command to an item and wait until the event bus acknowledge receiving the state, or after a timeout
// It returns true if openHAB acknowledge it's setting the desired state to the item (even if it's the same value as before).
// It returns false in case the acknowledged value is different than the command, or after timeout
func (i *Item) SendCommandWait(command StateValue, timeout time.Duration) (bool, error) {
	stateChan := make(chan string)
	subID := i.client.subscribe(i.Name(), event.TypeItemState, func(e event.Event) {
		if ev, ok := e.(event.ItemReceivedState); ok {
			stateChan <- ev.State
		}
	})
	defer func() {
		i.client.unsubscribe(subID)
		close(stateChan)
	}()

	err := i.SendCommand(command)
	if err != nil {
		return false, err
	}

	select {
	case state := <-stateChan:
		return state == command.String(), nil
	case <-time.After(timeout):
		return false, nil
	}
}

func (i *Item) stateFromString(state string) StateValue {
	switch i.mainType {
	default:
		return StringState(state)
	case ItemTypeSwitch:
		return SwitchState(state)
	case ItemTypeNumber:
		return MustParseDecimalState(state)
	}
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
}

func (i *Item) setInternalState(state string) {
	i.setInternalStateValue(i.stateFromString(state))
}
