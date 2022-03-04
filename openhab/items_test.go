package openhab

import (
	"sync"
	"testing"
	"time"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/openhabtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanLoadItemCreatedAfterLoading(t *testing.T) {
	wg := sync.WaitGroup{}
	server := openhabtest.NewServer(openhabtest.Config{
		SendEventsFromAPI: false,
		Log:               t,
	})
	err := server.SetItem(api.Item{
		Name:  "item1",
		State: "FIRST",
		Type:  "String",
	})
	require.NoError(t, err)

	client := NewClient(Config{URL: server.URL()})

	wg.Add(1)
	go func() {
		defer wg.Done()
		client.Start()
	}()

	time.Sleep(10 * time.Millisecond)
	// load existing item at the time of loading
	item1, err := client.items.getItem("item1")
	require.NoError(t, err)
	assert.NotNil(t, item1)
	assert.Equal(t, "FIRST", item1.state.String())

	// create a new item
	err = server.SetItem(api.Item{
		Name:  "item2",
		State: "SECOND",
		Type:  "String",
	})
	require.NoError(t, err)

	// load the new item
	item2, err := client.items.getItem("item2")
	require.NoError(t, err)
	assert.NotNil(t, item2)
	assert.Equal(t, "SECOND", item2.state.String())

	client.Stop()

	wg.Wait()
	server.Close()
}

func TestLoadingUnknownItemReturnsError(t *testing.T) {
	wg := sync.WaitGroup{}
	server := openhabtest.NewServer(openhabtest.Config{
		SendEventsFromAPI: false,
		Log:               t,
	})
	err := server.SetItem(api.Item{
		Name:  "item1",
		State: "FIRST",
		Type:  "String",
	})
	require.NoError(t, err)

	client := NewClient(Config{URL: server.URL()})

	wg.Add(1)
	go func() {
		defer wg.Done()
		client.Start()
	}()

	time.Sleep(10 * time.Millisecond)
	// load a non-existing item
	item2, err := client.items.getItem("item2")
	require.Error(t, err)
	assert.Nil(t, item2)

	client.Stop()

	wg.Wait()
	server.Close()
}
