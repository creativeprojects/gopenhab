package openhab

import (
	"context"
	"time"

	"github.com/creativeprojects/gopenhab/api"
)

type Items struct {
	client *Client
	cache  map[string]*Item
}

func newItems(client *Client) *Items {
	return &Items{
		client: client,
		cache:  make(map[string]*Item),
	}
}

func (items *Items) GetItem(name string) (*Item, error) {
	if len(items.cache) == 0 {
		// load them all now
		all := make([]api.Item, 0)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := items.client.getJSON(ctx, "items", &all)
		if err != nil {
			return nil, err
		}
		items.cache = make(map[string]*Item, len(all))
		for _, item := range all {
			items.cache[item.Name] = newItem(items.client, item.Name).set(item)
		}
	}
	return items.cache[name], nil
}
