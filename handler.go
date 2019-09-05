package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type handler struct {
	list *list
}

// error variables to handle expected errors
var (
	ErrCodeNotFound     = "URL_NOT_FOUND"
	ErrCodeRuntimeError = "RUNTIME_ERROR"
)

func (h *handler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx = &context{
		w: w,
		r: r,
	}

	// initialize the context and also prepare destroy
	ctx.init()
	defer ctx.destroy()

	// recovery/handle any runtime error
	defer func() {
		err := recover()
		if err != nil {
			ctx.Throw(ErrCodeRuntimeError, fmt.Errorf("%v", err))
		}
		h.caughtExceptions(ctx)
	}()

	// on context done, stop execution
	go func() {
		_ctx := r.Context()
		select {
		case <-_ctx.Done():
			ctx.end = true
		}
	}()

	// required "/" to match pattern
	var uri = r.RequestURI
	if !strings.HasSuffix(uri, sep) {
		uri += sep
	}

	// STEP 1: middlewares
	for _, handle := range h.list.middlewares {
		if ctx.end || ctx.code != "" || ctx.err != nil {
			break
		}

		if handle.pattern.test(uri) {
			ctx.params = handle.pattern.match(uri)
			handle.task(ctx)
		}
	}

	// STEP 2: routes
	for _, handle := range h.list.routes {
		if ctx.end || ctx.code != "" || ctx.err != nil {
			break
		}

		if r.Method == handle.method && handle.pattern.test(uri) {
			ctx.params = handle.pattern.match(uri)
			handle.task(ctx)
		}
	}

	// if no error and still not ended that means its NOT FOUND
	if !ctx.end && ctx.code == "" && ctx.err == nil {
		ctx.Throw(ErrCodeNotFound, errors.New("404 page not found"))
	}

	// STEP 3: errors
	for _, handle := range h.list.exceptions {
		if ctx.end || ctx.code == "" {
			break
		}

		if ctx.code == handle.code {
			handle.task(ctx.err, ctx)
		}
	}
}

func (h *handler) caughtExceptions(ctx *context) {
	defer h.recover(ctx)
	if !ctx.end {
		if h.list.uncaughtException != nil {
			h.list.uncaughtException(ctx.err, ctx)
		} else {
			ctx.unhandledException()
		}
	}
}

func (h *handler) recover(ctx *context) {
	err := recover()
	if err != nil {
		ctx.unhandledException()
	}
}
