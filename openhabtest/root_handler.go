package openhabtest

import (
	"fmt"
	"net/http"
	"strings"
)

type rootHandler struct {
	log     Logger
	routes  []route
	version Version
}

func newRootHandler(log Logger, routes []route, version Version) *rootHandler {
	return &rootHandler{
		log:     log,
		routes:  routes,
		version: version,
	}
}

func (h *rootHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	h.log.Logf("%s: %s", req.Method, req.URL.String())

	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) < 1 {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	if parts[0] != "rest" {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	if len(parts) == 1 {
		h.sendIndex(resp, req)
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

func (h *rootHandler) sendIndex(resp http.ResponseWriter, req *http.Request) {
	baseURL := req.URL.Scheme + "://" + req.URL.Host + strings.TrimSuffix(req.URL.Path, "/")
	result := ""
	if h.version == V2 {
		result = fmt.Sprintf(`{"version":"3","links":[{"type":"uuid","url":"%s/uuid"}]}`, baseURL)
	} else {
		result = fmt.Sprintf(`{"version":"4","locale":"en_GB","measurementSystem":"SI","runtimeInfo":{"version":"3.0.4","buildString":"Release Build"},"links":[{"type":"uuid","url":"%s/uuid"}]}`, baseURL)
	}
	resp.Write([]byte(result))
}
