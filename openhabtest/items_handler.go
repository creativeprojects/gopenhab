package openhabtest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
)

type itemsHandler struct {
	log         Logger
	items       map[string]api.Item
	itemsLocker sync.Mutex
	eventBus    *eventBus
}

func newItemsHandler(log Logger, bus *eventBus) *itemsHandler {
	return &itemsHandler{
		log:      log,
		items:    make(map[string]api.Item, 10),
		eventBus: bus,
	}
}

func (h *itemsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	encoder := json.NewEncoder(resp)

	if len(parts) == 2 && req.Method == http.MethodGet {
		// request is: get all items
		h.sendAllItems(encoder, resp)
		return
	}

	if len(parts) == 3 {
		if req.Method == http.MethodGet {
			// request is: get single item
			h.sendItem(parts[2], encoder, resp)
			return
		}

		if req.Method == http.MethodPost {
			// request is: send command
			h.receiveCommand(parts[2], encoder, resp, req)
			return
		}
	}

	if len(parts) == 4 && parts[3] == "state" {
		if req.Method == http.MethodGet {
			// request is: get item state
			h.sendItemState(parts[2], encoder, resp)
			return
		}

		if req.Method == http.MethodPut {
			// request is: set item state
			h.receiveState(parts[2], encoder, resp, req)
			return
		}
	}

	// fallback
	resp.WriteHeader(http.StatusNotFound)
}

func (h *itemsHandler) sendAllItems(encoder *json.Encoder, resp http.ResponseWriter) {
	data := h.getItems()
	err := encoder.Encode(&data)
	if err != nil {
		h.log.Logf("cannot encode data into JSON: %+v", data)
		resp.WriteHeader(http.StatusBadRequest)
	}
}

func (h *itemsHandler) sendItem(name string, encoder *json.Encoder, resp http.ResponseWriter) {
	data, ok := h.getItem(name)
	if ok {
		err := encoder.Encode(&data)
		if err != nil {
			h.log.Logf("cannot encode data into JSON: %+v", data)
			resp.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	// item not found
	resp.WriteHeader(http.StatusNotFound)
}

func (h *itemsHandler) receiveCommand(name string, encoder *json.Encoder, resp http.ResponseWriter, req *http.Request) {
	item, ok := h.getItem(name)
	if ok {
		state, err := io.ReadAll(req.Body)
		if err != nil || len(state) == 0 {
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
		oldState := item.State
		newState := string(state)
		item.State = newState
		h.setItem(item)
		resp.WriteHeader(http.StatusOK)

		if h.eventBus == nil {
			return
		}
		// now send the events to the bus
		topic, ev := EventString(event.NewItemReceivedCommand(name, "Test", newState))
		h.eventBus.Publish(topic, ev)
		topic, ev = EventString(event.NewItemReceivedState(name, "Test", newState))
		h.eventBus.Publish(topic, ev)
		if oldState != newState {
			topic, ev = EventString(event.NewItemStateChanged(name, "Test", oldState, newState))
			h.eventBus.Publish(topic, ev)
		}
		return
	}
	// item not found
	resp.WriteHeader(http.StatusNotFound)
}

func (h *itemsHandler) sendItemState(name string, encoder *json.Encoder, resp http.ResponseWriter) {
	data, ok := h.getItem(name)
	if ok {
		resp.Write([]byte(data.State))
		return
	}
	// item not found
	resp.WriteHeader(http.StatusNotFound)
}

func (h *itemsHandler) receiveState(name string, encoder *json.Encoder, resp http.ResponseWriter, req *http.Request) {
	item, ok := h.getItem(name)
	if ok {
		state, err := io.ReadAll(req.Body)
		if err != nil || len(state) == 0 {
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
		oldState := item.State
		newState := string(state)
		item.State = newState
		h.setItem(item)
		resp.WriteHeader(http.StatusAccepted)

		if h.eventBus == nil {
			return
		}
		// now send the events to the bus
		topic, ev := EventString(event.NewItemReceivedState(name, "Test", newState))
		h.eventBus.Publish(topic, ev)
		if oldState != newState {
			topic, ev = EventString(event.NewItemStateChanged(name, "Test", oldState, newState))
			h.eventBus.Publish(topic, ev)
		}
		return
	}
	// item not found
	resp.WriteHeader(http.StatusNotFound)
}

// setItem adds the new item, or replaces the existing one (with the same name)
func (h *itemsHandler) setItem(item api.Item) error {
	if item.Name == "" {
		return errors.New("missing item name")
	}
	h.itemsLocker.Lock()
	defer h.itemsLocker.Unlock()

	item.Editable = true
	if item.Tags == nil {
		item.Tags = []string{}
	}
	if item.GroupNames == nil {
		item.GroupNames = []string{}
	}
	if item.Members == nil {
		item.Members = []string{}
	}
	h.items[item.Name] = item
	return nil
}

// removeItem removes an existing item. It doesn't return an error if the item doesn't exist.
func (h *itemsHandler) removeItem(itemName string) error {
	if itemName == "" {
		return errors.New("missing item name")
	}
	h.itemsLocker.Lock()
	defer h.itemsLocker.Unlock()

	delete(h.items, itemName)
	return nil
}

func (h *itemsHandler) getItems() []api.Item {
	h.itemsLocker.Lock()
	defer h.itemsLocker.Unlock()

	all := make([]api.Item, len(h.items))
	i := 0
	for _, item := range h.items {
		all[i] = item
		i++
	}
	return all
}

func (h *itemsHandler) getItem(name string) (api.Item, bool) {
	h.itemsLocker.Lock()
	defer h.itemsLocker.Unlock()

	item, ok := h.items[name]
	return item, ok
}
