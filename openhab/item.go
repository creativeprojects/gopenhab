package openhab

import (
	"context"
	"time"

	"github.com/creativeprojects/gopenhab/api"
)

type Item struct {
	name     string
	data     api.Item
	mainType ItemType
	subType  string
	client   *Client
}

func newItem(client *Client, name string) *Item {
	return &Item{
		name:   name,
		client: client,
	}
}

func (i *Item) set(data api.Item) *Item {
	i.data = data
	i.mainType, i.subType = getItemType(i.data.Type)
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

// State always calls the api to return a fresh value from openHAB
func (i *Item) State() (StateValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	state, err := i.client.getString(ctx, "items/"+i.name+"/state")
	if err != nil {
		return i.stateFromString(""), err
	}
	return i.stateFromString(state), nil
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
