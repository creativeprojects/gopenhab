package openhabtest

import (
	"net/http/httptest"
	"sync"
)

// Server is a mock openHAB instance to use in tests
//
// Please note you should only connect ONE client per server.
// The long running event bus request can only be polled by one client at a time.
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

// SendRawEvent sends a raw JSON string event to the event bus. Example of a raw event:
//
// {"topic":"smarthome/items/LocalWeatherAndForecast_Current_Cloudiness/state","payload":"{\"type\":\"Quantity\",\"value\":\"20 %\"}","type":"ItemStateEvent"}
func (s *Server) SendRawEvent(event string) {
	s.eventChan <- event
}
