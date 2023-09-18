package openhab

import (
	"context"
	"sync"
	"time"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
)

const (
	itemsPath = "items/"
)

// Item represents an item in openHAB
type Item struct {
	name        string
	isGroup     bool
	data        api.Item
	state       State
	mainType    ItemType
	subType     string
	client      *Client
	apiLocker   sync.Mutex
	stateLocker sync.Mutex
	updated     time.Time
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
	if i.data.Type != "" {
		i.mainType, i.subType = getItemType(i.data.Type)
		i.isGroup = false
	} else if i.data.GroupType != "" {
		i.mainType, i.subType = getItemType(i.data.GroupType)
		i.isGroup = true
	}
	i.setInternalStateString(data.State)
	return i
}

func (i *Item) load() error {
	data := api.Item{}
	ctx, cancel := context.WithTimeout(context.Background(), i.client.config.TimeoutHTTP)
	defer cancel()

	i.client.addCounter(MetricItemLoad, 1, MetricItemName, i.name)
	err := i.client.getJSON(ctx, itemsPath+i.name, &data)
	if err != nil {
		return err
	}
	i.set(data)
	return nil
}

// Name returns the name of the item (an item name is unique in openHAB)
func (i *Item) Name() string {
	return i.name
}

// Type return the item type
func (i *Item) Type() ItemType {
	return i.mainType
}

// IsGroup returns true if the item is a group of items
func (i *Item) IsGroup() bool {
	return i.isGroup
}

// IsMemberOf returns true if the item is a member of the group
func (i *Item) IsMemberOf(groupName string) bool {
	for _, group := range i.data.GroupNames {
		if group == groupName {
			return true
		}
	}
	return false
}

// Updated returns the last time the item state was updated (doesn't necessarily mean the state was changed)
func (i *Item) Updated() time.Time {
	return i.updated
}

// State returns an internal cached value if available,
// or calls the api to return a fresh value from openHAB if not
//
// State value is automatically refreshed from openHAB events,
// so you should always get an accurate value.
//
// Please note if you just sent a state change command,
// the new value might not be reflected instantly,
// but only after openHAB sent a state changed event back.
func (i *Item) State() (State, error) {
	internalState := i.getInternalState()
	if internalState != nil {
		return internalState, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), i.client.config.TimeoutHTTP)
	defer cancel()

	i.client.addCounter(MetricItemLoadState, 1, MetricItemName, i.name)
	state, err := i.client.getString(ctx, itemsPath+i.name+"/state")
	if err != nil {
		return i.stateFromString(""), err
	}
	value := i.stateFromString(state)
	i.setInternalState(value)
	return value, nil
}

// SendCommand sends a command to an item
func (i *Item) SendCommand(command State) error {
	i.apiLocker.Lock()
	defer i.apiLocker.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), i.client.config.TimeoutHTTP)
	defer cancel()

	i.client.addCounter(MetricItemSetState, 1, MetricItemName, i.name)
	err := i.client.postString(ctx, itemsPath+i.name, command.String())
	if err != nil {
		return err
	}
	return nil
}

// SendCommandWait sends a command to an item and wait until the event bus acknowledge receiving the state, or after a timeout
// It returns true if openHAB acknowledge it's setting the desired state to the item (even if it's the same value as before).
// It returns false in case the acknowledged value is different than the command, or after timeout
func (i *Item) SendCommandWait(command State, timeout time.Duration) (bool, error) {
	stateChan := make(chan string, 1)
	subID := i.client.subscribe(i.Name(), event.TypeItemState, func(e event.Event) {
		if ev, ok := e.(event.ItemReceivedState); ok {
			stateChan <- ev.State
		}
	})
	defer func() {
		i.client.unsubscribe(subID)
	}()

	err := i.SendCommand(command)
	if err != nil {
		return false, err
	}

	select {
	case state := <-stateChan:
		return command.Equal(state), nil
	case <-time.After(timeout):
		return false, nil
	}
}

func (i *Item) stateFromString(state string) State {
	switch i.mainType {
	default:
		return StringState(state)
	case ItemTypeSwitch:
		return SwitchState(state)
	case ItemTypeNumber:
		return MustParseDecimalState(state)
	case ItemTypeDateTime:
		return MustParseDateTimeState(state)
	}
}

// getInternalState gets the internal state value: it does not trigger an API call to get the state.
func (i *Item) getInternalState() State {
	i.stateLocker.Lock()
	defer i.stateLocker.Unlock()

	return i.state
}

// setInternalState sets the internal state value: it does not trigger an API call to set the state.
// this method should be used after state received or changed events
func (i *Item) setInternalState(state State) {
	i.stateLocker.Lock()
	defer i.stateLocker.Unlock()

	i.state = state
	i.updated = time.Now()
}

func (i *Item) setInternalStateString(state string) {
	i.setInternalState(i.stateFromString(state))
}
