package openhabtest

import (
	"errors"
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
	err      error // contains a list of errors that happened during events
}

func newEventsHandler(bus *eventBus, done <-chan bool) *eventsHandler {
	return &eventsHandler{
		eventBus: bus,
		done:     done,
	}
}

func (h *eventsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/event-stream")

	subID := h.eventBus.Subscribe("", func(message string) {
		var err error
		_, err = resp.Write(streamPrefix)
		h.err = errors.Join(h.err, err)
		_, err = resp.Write([]byte(message))
		h.err = errors.Join(h.err, err)
		_, err = resp.Write(streamSuffix)
		h.err = errors.Join(h.err, err)

		if flusher, ok := resp.(http.Flusher); ok {
			flusher.Flush()
		}
	})
	defer h.eventBus.Unsubscribe(subID)

	<-h.done
}

func (h *eventsHandler) AsyncServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/event-stream")

	sendEvent := func(message string) {
		var err error
		_, err = resp.Write(streamPrefix)
		h.err = errors.Join(h.err, err)
		_, err = resp.Write([]byte(message))
		h.err = errors.Join(h.err, err)
		_, err = resp.Write(streamSuffix)
		h.err = errors.Join(h.err, err)
		if flusher, ok := resp.(http.Flusher); ok {
			flusher.Flush()
		}
	}

	wg := sync.WaitGroup{}
	wg.Wait()
	event := make(chan string)
	subID := h.eventBus.Subscribe("", func(message string) {
		event <- message
	})
	defer h.eventBus.Unsubscribe(subID)

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
