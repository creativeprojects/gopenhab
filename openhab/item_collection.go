package openhab

import (
	"context"
	"fmt"
	"sync"

	"github.com/creativeprojects/gopenhab/api"
)

// itemCollection represents the collection of items in openHAB
type itemCollection struct {
	client      *Client
	cache       map[string]*Item
	cacheLocker sync.Mutex
}

func newItems(client *Client) *itemCollection {
	return &itemCollection{
		client: client,
		cache:  nil,
	}
}

// getItem returns an openHAB item from its name.
// The very first call will try to load the items collection from openHAB.
func (items *itemCollection) getItem(name string) (*Item, error) {
	items.cacheLocker.Lock()
	defer items.cacheLocker.Unlock()

	if items.cache == nil {
		// load them all now
		err := items.loadCache()
		if err != nil {
			return nil, err
		}
	}
	// try to get the item from the cache
	if item, ok := items.cache[name]; ok {
		items.client.addCounter(MetricItemCacheHit, 1, MetricItemName, name)
		return item, nil
	}
	// try to call the API to get the item
	item := newItem(items.client, name)
	if err := item.load(); err == nil {
		items.cache[name] = item
		return item, nil
	}
	// item wasn't found
	items.client.addCounter(MetricItemNotFound, 1, MetricItemName, name)
	return nil, fmt.Errorf("item %q %w", name, ErrNotFound)
}

func (items *itemCollection) removeItem(name string) {
	items.cacheLocker.Lock()
	defer items.cacheLocker.Unlock()

	delete(items.cache, name)
}

// getMembersOf returns a list of items member of the group
func (items *itemCollection) getMembersOf(groupName string) ([]*Item, error) {
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

// refreshCache reloads the items from openHAB and updates the cache.
// This method is thread safe.
func (items *itemCollection) refreshCache() error {
	items.cacheLocker.Lock()
	defer items.cacheLocker.Unlock()

	return items.loadCache()
}

// loadCache loads all items into the cache.
// This method is NOT using the cacheLocker: it is the responsibility of the caller to do so.
func (items *itemCollection) loadCache() error {
	all, err := items.load()
	if err != nil {
		return err
	}

	items.cache = make(map[string]*Item, len(all))
	for _, item := range all {
		items.cache[item.Name] = newItem(items.client, item.Name).set(item)
	}
	items.client.setGauge(MetricItemsCacheSize, int64(len(items.cache)), "", "")
	return nil
}

// load all items from the API
func (items *itemCollection) load() ([]api.Item, error) {
	all := make([]api.Item, 0)
	ctx, cancel := context.WithTimeout(context.Background(), items.client.config.TimeoutHTTP)
	defer cancel()
	err := items.client.getJSON(ctx, "items", &all)
	if err != nil {
		return nil, err
	}
	return all, nil
}
