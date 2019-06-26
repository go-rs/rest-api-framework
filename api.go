/*!
 * rest-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package rest

import (
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/go-rs/rest-api-framework/utils"
)

type Handler func(ctx *Context)

/**
 * API - Application
 */
type API struct {
	prefix       string
	routes       []route
	interceptors []interceptor
	exceptions   []exception
	unhandled    Handler
}

/**
 * Route
 */
type route struct {
	method  string
	pattern string
	regex   *regexp.Regexp
	params  []string
	handle  Handler
}

/**
 * Request interceptor
 */
type interceptor struct {
	handle Handler
}

/**
 * Exception Route
 */
type exception struct {
	message string
	handle  Handler
}

/**
 * Common Route
 */
func (api *API) Route(method string, pattern string, handle Handler) {
	regex, params, err := utils.Compile(pattern)
	if err != nil {
		panic(err)
	}
	api.routes = append(api.routes, route{
		method:  method,
		pattern: pattern,
		regex:   regex,
		params:  params,
		handle:  handle,
	})
}

func (api *API) Use(handle Handler) {
	task := interceptor{
		handle: handle,
	}
	api.interceptors = append(api.interceptors, task)
}

func (api *API) All(pattern string, handle Handler) {
	api.Route("", pattern, handle)
}

func (api *API) Get(pattern string, handle Handler) {
	api.Route(http.MethodGet, pattern, handle)
}

func (api *API) Post(pattern string, handle Handler) {
	api.Route(http.MethodPost, pattern, handle)
}

func (api *API) Put(pattern string, handle Handler) {
	api.Route(http.MethodPut, pattern, handle)
}

func (api *API) Delete(pattern string, handle Handler) {
	api.Route(http.MethodDelete, pattern, handle)
}

func (api *API) Options(pattern string, handle Handler) {
	api.Route(http.MethodOptions, pattern, handle)
}

func (api *API) Head(pattern string, handle Handler) {
	api.Route(http.MethodHead, pattern, handle)
}

func (api *API) Patch(pattern string, handle Handler) {
	api.Route(http.MethodPatch, pattern, handle)
}

func (api *API) Exception(err string, handle Handler) {
	exp := exception{
		message: err,
		handle:  handle,
	}
	api.exceptions = append(api.exceptions, exp)
}

func (api *API) UnhandledException(handle Handler) {
	api.unhandled = handle
}

var (
	ErrNotFound          = errors.New("URL_NOT_FOUND")
	ErrUncaughtException = errors.New("UNCAUGHT_EXCEPTION")
)

/**
 * Required handle for http module
 */
func (api API) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	// STEP 1: initialize context
	ctx := Context{
		Request:  req,
		Response: res,
		Query:    req.URL.Query(),
	}

	ctx.init()
	defer ctx.destroy()

	defer func() {
		err := recover()
		if err != nil {
			log.Fatalln("uncaught exception - ", err)
			if !ctx.end {
				ctx.err = ErrUncaughtException
				ctx.unhandledException()
				return
			}
		}
	}()

	// STEP 2: execute all interceptors
	for _, task := range api.interceptors {
		if ctx.end || ctx.err != nil {
			break
		}

		task.handle(&ctx)
	}

	// STEP 3: check routes
	urlPath := []byte(req.URL.Path)
	for _, route := range api.routes {
		if ctx.end || ctx.err != nil {
			break
		}

		if (route.method == "" || route.method == req.Method) && route.regex.Match(urlPath) {
			ctx.found = route.method != "" //?
			ctx.Params = utils.Exec(route.regex, route.params, urlPath)
			route.handle(&ctx)
		}
	}

	// STEP 4: check handled exceptions
	for _, exp := range api.exceptions {
		if ctx.end || ctx.err == nil {
			break
		}

		if exp.message == ctx.err.Error() {
			exp.handle(&ctx)
		}
	}

	// STEP 5: unhandled exceptions
	if !ctx.end {
		if ctx.err == nil && !ctx.found {
			ctx.err = ErrNotFound
		}

		if api.unhandled != nil {
			api.unhandled(&ctx)
		}
	}

	// STEP 6: system handle
	if !ctx.end {
		ctx.unhandledException()
	}
}
