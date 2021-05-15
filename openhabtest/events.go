package openhabtest

import "net/http"

type eventsHandler struct{}

func newEventsHandler() eventsHandler {
	return eventsHandler{}
}

func (h eventsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	//
}
