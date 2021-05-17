package openhabtest

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"sync"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
)

// Server is a mock openHAB instance to use in tests
//
// Please note you should only connect ONE client per server.
// The long running event bus request can only be polled by one client at a time.
// It's very unlikely to need more than one client in a unit test anyway.
type Server struct {
	log         Logger
	server      *httptest.Server
	eventChan   chan string
	closing     chan bool
	closeLocker sync.Mutex
}

// NewServer creates a new mock openHAB instance to use in tests
func NewServer(log Logger) *Server {
	if log == nil {
		log = dummyLogger{}
	}
	eventChan := make(chan string)
	closing := make(chan bool)
	routes := []route{
		{"events", newEventsHandler(eventChan, closing)},
	}

	server := httptest.NewServer(newRootHandler(log, routes))
	return &Server{
		log:       log,
		server:    server,
		eventChan: eventChan,
		closing:   closing,
	}
}

// URL returns the local URL of a mock openHAB server to use inside a unit test
func (s *Server) URL() string {
	if s.server == nil {
		panic("no instance of http server")
	}
	return s.server.URL
}

// Close the mock openHAB server. The call will also close any long running request to the event bus API.
func (s *Server) Close() {
	s.closeLocker.Lock()
	defer s.closeLocker.Unlock()

	if s.closing != nil {
		close(s.closing)
		s.closing = nil
	}
	if s.server != nil {
		s.server.Close()
		s.server = nil
	}
}

// RawEvent sends a raw JSON string event to the event bus. Example of a raw event:
//
// {"topic":"smarthome/items/LocalWeatherAndForecast_Current_Cloudiness/state","payload":"{\"type\":\"Quantity\",\"value\":\"20 %\"}","type":"ItemStateEvent"}
func (s *Server) RawEvent(event string) {
	s.eventChan <- event
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
