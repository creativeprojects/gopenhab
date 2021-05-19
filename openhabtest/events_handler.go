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

	subId := h.eventBus.Subscribe("", func(message string) {
		resp.Write(streamPrefix)
		resp.Write([]byte(message))
		resp.Write(streamSuffix)

		if flusher, ok := resp.(http.Flusher); ok {
			flusher.Flush()
		}
	})
	defer h.eventBus.Unsubscribe(subId)

	<-h.done
}

func (h *eventsHandler) AsyncServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/event-stream")

	sendEvent := func(message string) {
		resp.Write(streamPrefix)
		resp.Write([]byte(message))
		resp.Write(streamSuffix)
		if flusher, ok := resp.(http.Flusher); ok {
			flusher.Flush()
		}
	}

	wg := sync.WaitGroup{}
	wg.Wait()
	event := make(chan string)
	subId := h.eventBus.Subscribe("", func(message string) {
		event <- message
	})
	defer h.eventBus.Unsubscribe(subId)

	for {
		// first select: we wait for either data or the exit signal
		select {
		case message := <-event:
			sendEvent(message)
		case <-h.done:
			// got an exit signal: but now we need to drain the event channel before leaving
			select {
			case message := <-event:
				sendEvent(message)
			default:
				return
			}
		}
	}
}
