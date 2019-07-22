// go-rs/rest-api-framework
// Copyright(c) 2019 Roshan Gade. All rights reserved.
// MIT Licensed

package rest

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-rs/rest-api-framework/utils"
)

// Handler function is used to perform a specified task on request.
// In handler function, you will get rest context object,
// which carries request, response writer objects and other methods too.
type Handler func(ctx *Context)

// API
// It provides all methods, which are required to setup all kind of routes.
// It manages request interceptor/middlewares. which are used to intercept requests calls before performing actual operation,
// such as authentication, method override, request logger, etc.
// Also, it handles errors, which are thrown by users
type API struct {
	prefix       string
	routes       []route
	interceptors []interceptor
	exceptions   []exception
	unhandled    Handler
}

// routes, which help to find and execute/perform the exact or matched url path
type route struct {
	method  string
	pattern string
	regex   *regexp.Regexp
	params  []string
	handle  Handler
}

// request interceptors, which help to intercept every request before executing the respective route
type interceptor struct {
	handle Handler
}

// user exceptions, which is a common way to handle an error thrown by user
type exception struct {
	code   string
	handle Handler
}

// Initialize an API with prefix value and return the API pointer
func New(prefix string) *API {
	prefix = utils.TrimSuffix(prefix, "/")
	return &API{
		prefix: prefix,
	}
}

// Route method is used to define specific routes with handler.
// You can use http method, declare patten and finally you have to pass handler
func (api *API) Route(method string, pattern string, handle Handler) {
	pattern = api.prefix + pattern
	regex, params, err := utils.Compile(pattern)
	if err != nil {
		panic(err)
	}
	api.routes = append(api.routes, route{
		method:  strings.ToUpper(method),
		pattern: pattern,
		regex:   regex,
		params:  params,
		handle:  handle,
	})
}

// Use method is use to declare interceptors/middlewares
// It executes on all type method request
func (api *API) Use(handle Handler) {
	task := interceptor{
		handle: handle,
	}
	api.interceptors = append(api.interceptors, task)
}

// All method is slightly similar to Use method,
// but in All method you can use pattern before intercepting any request
func (api *API) All(pattern string, handle Handler) {
	api.Route("", pattern, handle)
}

// Get method is used for GET http method with specific pattern
func (api *API) Get(pattern string, handle Handler) {
	api.Route(http.MethodGet, pattern, handle)
}

// Post method is used for POST http method with specific pattern
func (api *API) Post(pattern string, handle Handler) {
	api.Route(http.MethodPost, pattern, handle)
}

// Put method is used for PUT http method with specific pattern
func (api *API) Put(pattern string, handle Handler) {
	api.Route(http.MethodPut, pattern, handle)
}

// Delete method is used for DELETE http method with specific pattern
func (api *API) Delete(pattern string, handle Handler) {
	api.Route(http.MethodDelete, pattern, handle)
}

// Options method is used for OPTIONS http method with specific pattern
func (api *API) Options(pattern string, handle Handler) {
	api.Route(http.MethodOptions, pattern, handle)
}

// Head method is used for HEAD http method with specific pattern
func (api *API) Head(pattern string, handle Handler) {
	api.Route(http.MethodHead, pattern, handle)
}

// Patch method is used for PATCH http method with specific pattern
func (api *API) Patch(pattern string, handle Handler) {
	api.Route(http.MethodPatch, pattern, handle)
}

// OnError method is used to handle a custom errors thrown by users
func (api *API) OnError(code string, handle Handler) {
	exp := exception{
		code:   code,
		handle: handle,
	}
	api.exceptions = append(api.exceptions, exp)
}

// OnErrors method is used to handle a custom errors thrown by users
func (api *API) OnErrors(codes []string, handle Handler) {
	for _, code := range codes {
		exp := exception{
			code:   code,
			handle: handle,
		}
		api.exceptions = append(api.exceptions, exp)
	}
}

// UnhandledException method is used to handle all unhandled exceptions
func (api *API) UnhandledException(handle Handler) {
	api.unhandled = handle
}

// error variables to handle expected errors
var (
	ErrCodeNotFound     = "URL_NOT_FOUND"
	ErrCodeRuntimeError = "RUNTIME_ERROR"
)

// It's required handle for http module.
// Every request travels from this method.
func (api *API) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	// STEP 1: initialize context
	ctx := &Context{
		Request:  req,
		Response: res,
		Query:    req.URL.Query(),
	}

	ctx.init()
	defer ctx.destroy()

	// recovery/handle any runtime error
	defer func() {
		err := recover()
		if err != nil {
			if !ctx.end {
				ctx.code = ErrCodeRuntimeError
				ctx.err = fmt.Errorf("%v", err)
				ctx.unhandledException()
			}
			return
		}
	}()

	// On context done, stop execution
	go func() {
		c := req.Context()
		select {
		case <-c.Done():
			ctx.End()
		}
	}()

	// STEP 2: execute all interceptors
	for _, task := range api.interceptors {
		if ctx.shouldBreak() {
			break
		}

		task.handle(ctx)
	}

	// STEP 3: check routes
	urlPath := req.URL.Path
	for _, route := range api.routes {
		if ctx.shouldBreak() {
			break
		}

		if (route.method == "" || strings.EqualFold(route.method, req.Method)) && route.regex.MatchString(urlPath) {
			ctx.Params = utils.Exec(route.regex, route.params, urlPath)
			route.handle(ctx)
		}
	}

	// STEP 4: check handled exceptions
	for _, exp := range api.exceptions {
		if ctx.shouldBreak() {
			break
		}

		if strings.EqualFold(exp.code, ctx.code) {
			exp.handle(ctx)
		}
	}

	// STEP 5: unhandled exceptions
	if !ctx.end {
		// if no error and still not ended that means it NOT FOUND
		if ctx.code == "" {
			ctx.Throw(ErrCodeNotFound)
		}

		// if user has custom unhandled function, then execute it
		if api.unhandled != nil {
			api.unhandled(ctx)
		}
	}

	// STEP 6: system handle
	if !ctx.end {
		ctx.unhandledException()
	}
}
