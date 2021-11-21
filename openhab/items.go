package openhab

import (
	"context"
	"fmt"
	"sync"

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

// getItem returns an openHAB item from its name.
// The very first call of GetItem will try to load the items collection from openHAB.
func (items *Items) getItem(name string) (*Item, error) {
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
	return nil, fmt.Errorf("item %q %w", name, ErrorNotFound)
}

// getMembersOf returns a list of items member of the group
func (items *Items) getMembersOf(groupName string) ([]*Item, error) {
	items.cacheLocker.Lock()
	defer items.cacheLocker.Unlock()

	if items.cache == nil {
		// load them all now
		err := items.loadCache()
		if err != nil {
			return nil, err
		}
	}

	members := []*Item{}
	for _, item := range items.cache {
		if item.IsMemberOf(groupName) {
			members = append(members, item)
		}
	}
	return members, nil
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

// load all items from the API
func (items *Items) load() ([]api.Item, error) {
	all := make([]api.Item, 0)
	ctx, cancel := context.WithTimeout(context.Background(), items.client.config.TimeoutHTTP)
	defer cancel()
	err := items.client.getJSON(ctx, "items", &all)
	if err != nil {
		return nil, err
	}
	return all, nil
}
