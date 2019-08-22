package rest

import (
	"net/http"
	"strings"
)

type handler struct {
	list *list
}

func (h *handler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx = &Context{
		Writer:  w,
		Request: r,
	}

	// required "/" to match pattern
	var uri = r.RequestURI
	if !strings.HasSuffix(uri, sep) {
		uri += sep
	}

	// on context done, stop execution
	go func() {
		_ctx := r.Context()
		select {
		case <-_ctx.Done():
			//TODO: end context
		}
	}()

	// STEP 1: middlewares
	for _, handle := range h.list.middlewares {
		if handle.pattern.test(uri) {
			handle.task(ctx)
		}
	}

	// STEP 2: routes
	for _, handle := range h.list.routes {
		if r.Method == handle.method && handle.pattern.test(uri) {
			handle.task(ctx)
		}
	}
}
