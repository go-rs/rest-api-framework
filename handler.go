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

	ctx.init()
	defer ctx.destroy()

	// recovery/handle any runtime error
	defer func() {
		err := recover()
		if err != nil {
			if !ctx.end {
				defer h.recover(ctx)
				ctx.code = ErrCodeRuntimeError
				ctx.err = fmt.Errorf("%v", err)
				if h.list.uncaughtException != nil {
					h.list.uncaughtException(ctx.err, ctx)
				} else {
					ctx.unhandledException()
				}
			}
			return
		}
	}()

	// on context done, stop execution
	go func() {
		_ctx := r.Context()
		select {
		case <-_ctx.Done():
			//TODO: end context
		}
	}()

	// required "/" to match pattern
	var uri = r.RequestURI
	if !strings.HasSuffix(uri, sep) {
		uri += sep
	}

	// STEP 1: middlewares
	for _, handle := range h.list.middlewares {
		if ctx.end || ctx.err != nil {
			break
		}

		if handle.pattern.test(uri) {
			handle.task(ctx)
		}
	}

	// STEP 2: routes
	for _, handle := range h.list.routes {
		if ctx.end || ctx.err != nil {
			break
		}

		if r.Method == handle.method && handle.pattern.test(uri) {
			handle.task(ctx)
		}
	}

	// STEP 3: errors
	for _, handle := range h.list.exceptions {
		if ctx.end || ctx.err == nil {
			break
		}

		if ctx.code == handle.code {
			handle.task(ctx.err, ctx)
		}
	}

	// STEP 5: unhandled exceptions
	if !ctx.end {
		// if no error and still not ended that means it NOT FOUND
		if ctx.code == "" {
			ctx.Throw(ErrCodeNotFound, errors.New("404 page not found"))
		}

		// if user has custom unhandled function, then execute it
		if h.list.uncaughtException != nil {
			h.list.uncaughtException(ctx.err, ctx)
		}
	}

	// STEP 6: system handle
	if !ctx.end {
		ctx.unhandledException()
	}
}

func (h *handler) recover(ctx *context) {
	err := recover()
	if err != nil {
		ctx.unhandledException()
	}
}
