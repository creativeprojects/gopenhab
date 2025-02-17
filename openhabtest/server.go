package openhabtest

import (
	"net/http/httptest"
	"sync"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
)

// Server is a mock openHAB instance to use in tests.
type Server struct {
	log           Logger
	version       Version
	server        *httptest.Server
	eventBus      *eventBus
	closeLocker   sync.Mutex
	itemsHandler  *itemsHandler
	done          chan bool
	closed        bool
	eventsHandler *eventsHandler
}

// NewServer creates a new mock openHAB instance to use in tests
func NewServer(config Config) *Server {
	if config.Log == nil {
		config.Log = dummyLogger{}
	}
	done := make(chan bool)
	bus := newEventBus()
	autoBus := bus
	if !config.SendEventsFromAPI {
		// don't send the events automatically => we don't send the instance of the events bus to handlers
		autoBus = nil
	}
	eventsHandler := newEventsHandler(bus, done)
	itemsHandler := newItemsHandler(config.Log, autoBus, config.Version)
	routes := []route{
		{"events", eventsHandler},
		{"items", itemsHandler},
	}

	server := httptest.NewServer(newRootHandler(config.Log, routes, config.Version))
	return &Server{
		log:           config.Log,
		version:       config.Version,
		server:        server,
		eventBus:      bus,
		itemsHandler:  itemsHandler,
		done:          done,
		eventsHandler: eventsHandler,
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

// EventsErr returns an error if any happened from the event endpoints.
//
// A non-nil error returned by EventsErr implements the Unwrap() []error method.
func (s *Server) EventsErr() error {
	return s.eventsHandler.err
}

// ItemsErr returns an error if any happened from the item endpoints.
//
// A non-nil error returned by ItemsErr implements the Unwrap() []error method.
func (s *Server) ItemsErr() error {
	return s.itemsHandler.err
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

	if s.server != nil {
		s.server.Close()
		s.server = nil
	}
}

// RawEvent sends a raw JSON string event to the event bus. Example of a raw event:
//
//	{"topic":"smarthome/items/LocalWeatherAndForecast_Current_Cloudiness/state","payload":"{\"type\":\"Quantity\",\"value\":\"20 %\"}","type":"ItemStateEvent"}
//
// A topic parameter is needed for subscriber topic filtering, and to avoid decoding the event string unnecessarily.
func (s *Server) RawEvent(topic, event string) {
	s.eventBus.Publish(topic, event)
}

// Event sends a event.Event to the mock openHAB event bus
func (s *Server) Event(e event.Event) {
	if e == nil {
		return
	}
	topic, ev := EventString(e, topicPrefix(s.version))
	if topic != "" && ev != "" {
		s.log.Logf("sending event %s on topic %s", ev, topic)
		s.RawEvent(topic, ev)
	}
}

// SetItem adds the new item, or replaces the existing one (with the same name).
// If Link property is not set, it will be automatically set
func (s *Server) SetItem(item api.Item) error {
	if item.Link == "" {
		item.Link = s.URL() + "/rest/items/" + item.Name
	}
	return s.itemsHandler.setItem(item)
}

// RemoveItem removes an existing item. It doesn't return an error if the item doesn't exist.
func (s *Server) RemoveItem(itemName string) error {
	return s.itemsHandler.removeItem(itemName)
}
