package openhab

import (
	"context"
	"sync"
	"time"

	"github.com/creativeprojects/gopenhab/api"
)

// Items represents the collection of items in openHAB
type Items struct {
	client      *Client
	cache       map[string]*Item
	cacheLocker sync.Mutex
}

func newItems(client *Client) *Items {
	return &Items{
		client: client,
		cache:  nil,
	}
}

// GetItem returns an openHAB item from its name.
// The very first call of GetItem will try to load the items collection from openHAB.
func (items *Items) GetItem(name string) (*Item, error) {
	items.cacheLocker.Lock()
	defer items.cacheLocker.Unlock()

	if items.cache == nil {
		// load them all now
		err := items.loadCache()
		if err != nil {
			return nil, err
		}
	}
	if item, ok := items.cache[name]; ok {
		return item, nil
	}
	return nil, ErrorNotFound
}

// loadCache loads all items into the cache.
// This method is NOT using the cacheLocker: it is the responsability of the caller to do so.
func (items *Items) loadCache() error {
	all, err := items.load()
	if err != nil {
		return err
	}

	items.cache = make(map[string]*Item, len(all))
	for _, item := range all {
		items.cache[item.Name] = newItem(items.client, item.Name).set(item)
	}
	return nil
}

func (items *Items) load() ([]api.Item, error) {
	all := make([]api.Item, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := items.client.getJSON(ctx, "items", &all)
	if err != nil {
		return nil, err
	}
	return all, nil
}
