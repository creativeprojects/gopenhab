package openhabtest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/creativeprojects/gopenhab/api"
)

type itemsHandler struct {
	log         Logger
	items       map[string]api.Item
	itemsLocker sync.Mutex
}

func newItemsHandler(log Logger) *itemsHandler {
	return &itemsHandler{
		log:   log,
		items: make(map[string]api.Item, 10),
	}
}

func (h *itemsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	encoder := json.NewEncoder(resp)

	if len(parts) == 2 && req.Method == http.MethodGet {
		// download all items
		data := h.getItems()
		err := encoder.Encode(&data)
		if err != nil {
			h.log.Logf("cannot encode data into JSON: %+v", data)
			resp.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	if len(parts) == 3 {
		if req.Method == http.MethodGet {
			// get single item
			data, ok := h.getItem(parts[2])
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
			return
		}

		if req.Method == http.MethodPost {
			// send command
			item, ok := h.getItem(parts[2])
			if ok {
				state, err := io.ReadAll(req.Body)
				if err != nil || len(state) == 0 {
					resp.WriteHeader(http.StatusBadRequest)
					return
				}
				item.State = string(state)
				h.setItem(item)
				resp.WriteHeader(http.StatusOK)
				return
			}
			// item not found
			resp.WriteHeader(http.StatusNotFound)
			return
		}
	}

	if len(parts) == 4 && parts[3] == "state" {
		if req.Method == http.MethodGet {
			// get item state
			data, ok := h.getItem(parts[2])
			if ok {
				resp.Write([]byte(data.State))
				return
			}
			// item not found
			resp.WriteHeader(http.StatusNotFound)
			return
		}

		if req.Method == http.MethodPut {
			// set item state
			item, ok := h.getItem(parts[2])
			if ok {
				state, err := io.ReadAll(req.Body)
				if err != nil || len(state) == 0 {
					resp.WriteHeader(http.StatusBadRequest)
					return
				}
				item.State = string(state)
				h.setItem(item)
				resp.WriteHeader(http.StatusAccepted)
				return
			}
			// item not found
			resp.WriteHeader(http.StatusNotFound)
			return
		}
	}

	// fallback
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
