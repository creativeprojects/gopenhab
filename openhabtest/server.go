package openhabtest

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

type Server struct {
	server *httptest.Server
}

func NewServer() *Server {
	server := httptest.NewServer(newRootHandler())
	return &Server{
		server: server,
	}
}

func (s Server) URL() string {
	if s.server == nil {
		panic("no instance of http server")
	}
	return s.server.URL
}

type rootHandler struct {
	eventsHandler eventsHandler
}

func newRootHandler() rootHandler {
	events := newEventsHandler()
	return rootHandler{
		eventsHandler: events,
	}
}

func (h rootHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	parts := strings.Split(req.URL.Path, "/")
	if len(parts) < 2 {
		resp.WriteHeader(http.StatusNotFound)
	}
	if parts[0] != "rest" {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	switch parts[1] {
	case "events":
		h.eventsHandler.ServeHTTP(resp, req)
	default:
		resp.WriteHeader(http.StatusNotFound)
	}
}
