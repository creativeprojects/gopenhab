package openhabtest

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"sync"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
)

// Server is a mock openHAB instance to use in tests.
type Server struct {
	log         Logger
	server      *httptest.Server
	eventBus    *eventBus
	closeLocker sync.Mutex
	items       *itemsHandler
	done        chan bool
	closed      bool
}

// NewServer creates a new mock openHAB instance to use in tests
func NewServer(log Logger) *Server {
	if log == nil {
		log = dummyLogger{}
	}
	done := make(chan bool)
	bus := newEventBus()
	items := newItemsHandler(log)
	routes := []route{
		{"events", newEventsHandler(bus, done)},
		{"items", items},
	}

	server := httptest.NewServer(newRootHandler(log, routes))
	return &Server{
		log:      log,
		server:   server,
		eventBus: bus,
		items:    items,
		done:     done,
	}
}

// URL returns the local URL of a mock openHAB server to use inside a unit test.
// The URL returned has no trailing '/'
func (s *Server) URL() string {
	if s.server == nil {
		panic("no instance of http server")
	}
	return s.server.URL
}

// Close the mock openHAB server. The call will also close any long running request to the event bus API.
// The method can safely be called multiple times.
func (s *Server) Close() {
	s.closeLocker.Lock()
	defer s.closeLocker.Unlock()

	if s.closed {
		return
	}
	s.closed = true

	close(s.done)
	s.eventBus.Close()

	if s.server != nil {
		s.server.CloseClientConnections()
		s.server.Close()
		s.server = nil
	}
}

// RawEvent sends a raw JSON string event to the event bus. Example of a raw event:
//
// {"topic":"smarthome/items/LocalWeatherAndForecast_Current_Cloudiness/state","payload":"{\"type\":\"Quantity\",\"value\":\"20 %\"}","type":"ItemStateEvent"}
func (s *Server) RawEvent(event string) {
	s.eventBus.Publish("", event)
}

// Event sends a event.Event to the mock openHAB event bus
func (s *Server) Event(e event.Event) {
	if e == nil {
		return
	}
	switch ev := e.(type) {
	case event.ItemReceivedCommand:
		rawPayload, err := json.Marshal(api.EventCommand{
			Type:  ev.CommandType,
			Value: ev.Command,
		})
		if err != nil {
			panic(err)
		}
		rawEvent, err := json.Marshal(api.EventMessage{
			Topic:   ev.Topic(),
			Payload: string(rawPayload),
			Type:    api.EventItemCommand,
		})
		if err != nil {
			panic(err)
		}
		s.RawEvent(string(rawEvent))

	case event.ItemReceivedState:
		rawPayload, err := json.Marshal(api.EventState{
			Type:  ev.StateType,
			Value: ev.State,
		})
		if err != nil {
			panic(err)
		}
		rawEvent, err := json.Marshal(api.EventMessage{
			Topic:   ev.Topic(),
			Payload: string(rawPayload),
			Type:    api.EventItemState,
		})
		if err != nil {
			panic(err)
		}
		s.RawEvent(string(rawEvent))

	case event.ItemStateChanged:
		rawPayload, err := json.Marshal(api.EventStateChanged{
			Type:     ev.NewStateType,
			Value:    ev.NewState,
			OldType:  ev.PreviousStateType,
			OldValue: ev.PreviousState,
		})
		if err != nil {
			panic(err)
		}
		rawEvent, err := json.Marshal(api.EventMessage{
			Topic:   ev.Topic(),
			Payload: string(rawPayload),
			Type:    api.EventItemStateChanged,
		})
		if err != nil {
			panic(err)
		}
		s.RawEvent(string(rawEvent))

	default:
		panic(fmt.Sprintf("event type %d not handled", e.Type()))
	}
}

// SetItem adds the new item, or replaces the existing one (with the same name).
// If Link property is not set, it will be automatically set
func (s *Server) SetItem(item api.Item) error {
	if item.Link == "" {
		item.Link = s.URL() + "/rest/items/" + item.Name
	}
	return s.items.setItem(item)
}

// RemoveItem removes an existing item. It doesn't return an error if the item doesn't exist.
func (s *Server) RemoveItem(itemName string) error {
	return s.items.removeItem(itemName)
}
