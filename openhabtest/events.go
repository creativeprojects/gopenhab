package openhabtest

import (
	"net/http"
	"sync"
)

var (
	streamPrefix = []byte("event: message\ndata: ")
	streamSuffix = []byte("\n\n")
)

type eventsHandler struct {
	eventChan chan string
	closing   chan bool
	safe      sync.Mutex
}

func newEventsHandler(eventChan chan string, closing chan bool) *eventsHandler {
	return &eventsHandler{
		eventChan: eventChan,
		closing:   closing,
	}
}

func (h *eventsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.safe.Lock()
	defer h.safe.Unlock()

	resp.Header().Add("Content-Type", "text/event-stream")

	for {
		select {
		case <-h.closing:
			return
		case event := <-h.eventChan:
			resp.Write(streamPrefix)
			resp.Write([]byte(event))
			resp.Write(streamSuffix)
			if flusher, ok := resp.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}
