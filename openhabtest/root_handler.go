package openhabtest

import (
	"net/http"
	"strings"
)

type rootHandler struct {
	log    Logger
	routes []route
}

func newRootHandler(log Logger, routes []route) *rootHandler {
	return &rootHandler{
		log:    log,
		routes: routes,
	}
}

func (h *rootHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	h.log.Logf("%s: %s", req.Method, req.URL.String())

	parts := strings.Split(strings.TrimPrefix(req.URL.Path, "/"), "/")
	if len(parts) < 2 {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	if parts[0] != "rest" {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	for _, route := range h.routes {
		if route.prefix == parts[1] {
			route.handler.ServeHTTP(resp, req)
			return
		}
	}
	resp.WriteHeader(http.StatusNotFound)
}
