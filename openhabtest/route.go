package openhabtest

import "net/http"

type route struct {
	prefix  string
	handler http.Handler
}
