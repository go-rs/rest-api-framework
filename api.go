/*!
 * rest-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package rest

import (
	"./utils"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

/**
 * API - Application
 */
type API struct {
	routes       []route
	interceptors []interceptor
	exceptions   []exception
	unhandled    func(ctx *Context)
}

/**
 * Route
 */
type route struct {
	method  string
	pattern string
	regex   *regexp.Regexp
	params  []string
	handle  func(ctx *Context)
}

/**
 * Request interceptor
 */
type interceptor struct {
	handle func(ctx *Context)
}

/**
 * Exception Route
 */
type exception struct {
	message string
	handle  func(ctx *Context)
}

/**
 * Common Route
 */
func (api *API) Route(method string, pattern string, handle func(ctx *Context)) {
	regex, params, err := utils.Compile(pattern)
	if err != nil {
		fmt.Println("Error in pattern", err)
		panic(1)
	}
	api.routes = append(api.routes, route{
		method:  method,
		pattern: pattern,
		regex:   regex,
		params:  params,
		handle:  handle,
	})
}

/**
 * Required handle for http module
 */
func (api API) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	urlPath := []byte(req.URL.Path)

	ctx := Context{
		Request:  req,
		Response: res,
	}

	// STEP 1: initialize context
	ctx.init()
	defer ctx.destroy()

	// STEP 2: execute all interceptors
	for _, task := range api.interceptors {
		if ctx.end || ctx.err != nil {
			break
		}

		task.handle(&ctx)
	}

	// STEP 3: check routes
	for _, route := range api.routes {
		if ctx.end || ctx.err != nil {
			break
		}

		if route.method == req.Method && route.regex.Match(urlPath) {
			ctx.found = true
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
			ctx.err = errors.New("URL_NOT_FOUND")
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

func (api *API) Use(handle func(ctx *Context)) {
	task := interceptor{
		handle: handle,
	}
	api.interceptors = append(api.interceptors, task)
}

func (api *API) GET(pattern string, handle func(ctx *Context)) {
	api.Route("GET", pattern, handle)
}

func (api *API) POST(pattern string, handle func(ctx *Context)) {
	api.Route("POST", pattern, handle)
}

func (api *API) PUT(pattern string, handle func(ctx *Context)) {
	api.Route("PUT", pattern, handle)
}

func (api *API) DELETE(pattern string, handle func(ctx *Context)) {
	api.Route("DELETE", pattern, handle)
}

func (api *API) OPTIONS(pattern string, handle func(ctx *Context)) {
	api.Route("OPTIONS", pattern, handle)
}

func (api *API) HEAD(pattern string, handle func(ctx *Context)) {
	api.Route("HEAD", pattern, handle)
}

func (api *API) PATCH(pattern string, handle func(ctx *Context)) {
	api.Route("PATCH", pattern, handle)
}

func (api *API) Exception(err string, handle func(ctx *Context)) {
	exp := exception{
		message: err,
		handle:  handle,
	}
	api.exceptions = append(api.exceptions, exp)
}

func (api *API) UnhandledException(handle func(ctx *Context)) {
	api.unhandled = handle
}
