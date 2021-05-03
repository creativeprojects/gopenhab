package openhab

import (
	"context"
	"time"

	"github.com/creativeprojects/gopenhab/api"
)

type Item struct {
	name   string
	data   api.Item
	client *Client
}

func newItem(client *Client, name string) *Item {
	return &Item{
		name:   name,
		client: client,
	}
}

func (i *Item) set(data api.Item) *Item {
	i.data = data
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
	i.data = data
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

func (i *Item) Type() api.ItemType {
	itemType, _ := api.GetItemType(i.getData().Type)
	return itemType
}

// State always calls the api to return a fresh value from openHAB
func (i *Item) State() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	state, err := i.client.getString(ctx, "items/"+i.name+"/state")
	if err != nil {
		return "", err
	}
	return state, nil
}

func (i *Item) SendCommand(command string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.client.postString(ctx, "items/"+i.name, command)
	if err != nil {
		return err
	}
	return nil
}
