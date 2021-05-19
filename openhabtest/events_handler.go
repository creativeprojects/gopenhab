package openhabtest

import (
	"net/http"
)

var (
	streamPrefix = []byte("event: message\ndata: ")
	streamSuffix = []byte("\n\n")
)

type eventsHandler struct {
	eventBus *eventBus
	done     <-chan bool
}

func newEventsHandler(bus *eventBus, done <-chan bool) *eventsHandler {
	return &eventsHandler{
		eventBus: bus,
		done:     done,
	}
}

func (h *eventsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/event-stream")

	event := make(chan string)
	subId := h.eventBus.Subscribe("", func(message string) {
		event <- message
	})
	defer h.eventBus.Unsubscribe(subId)

	for {
		select {
		case <-h.done:
			return
		case message := <-event:
			resp.Write(streamPrefix)
			resp.Write([]byte(message))
			resp.Write(streamSuffix)
			if flusher, ok := resp.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}
