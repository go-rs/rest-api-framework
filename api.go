/*!
 * utils-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package rest

import (
	"./utils"
	"fmt"
	"net/http"
	"regexp"
)

type API struct {
	routes       []route
	interceptors []interceptor
	errors       []errorHandler
}

type route struct {
	method  string
	pattern string
	regex   *regexp.Regexp
	params  []string
	handle  func(ctx *Context)
}

type interceptor struct {
	handle func(ctx *Context)
}

//TODO: check errors in language specifications
type errorHandler struct {
	code   string
	handle func(ctx *Context)
}

func (api *API) handler(method string, pattern string, handle func(ctx *Context)) {
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

func (api API) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("-------------------------------------")
	//fmt.Println("URI: ", req.RequestURI)
	//fmt.Println("Method: ", req.Method)
	//fmt.Println("Form: ", req.Form)
	//fmt.Println("RawPath: ", req.URL.RawPath)
	//fmt.Println("Header: ", req.Header)
	//fmt.Println("Path: ", req.URL.Path)
	//fmt.Println("RawQuery: ", req.URL.RawQuery)
	//fmt.Println("Opaque: ", req.URL.Opaque)
	//fmt.Println("-------------------------------------")
	urlPath := []byte(req.URL.Path)

	ctx := Context{
		req: req,
		res: res,
	}

	ctx.init()

	for _, route := range api.interceptors {
		if ctx.err != nil {
			break
		}

		if ctx.end == true {
			break
		}

		route.handle(&ctx)
	}

	for _, route := range api.routes {
		if ctx.err != nil {
			break
		}

		if ctx.end == true {
			break
		}

		if route.method == req.Method && route.regex.Match(urlPath) {
			ctx.Params = utils.Exec(route.regex, route.params, urlPath)
			route.handle(&ctx)
		}
	}

	//TODO: NOT FOUND, INTERNAL SERVER ERROR, ETC

	for _, route := range api.errors {
		if ctx.end == true {
			break
		}

		if ctx.err != nil && route.code == ctx.err.Error() {
			route.handle(&ctx)
		}
	}
}

func (api *API) Use(handle func(ctx *Context)) {
	api.interceptors = append(api.interceptors, interceptor{handle: handle})
}

func (api *API) GET(pattern string, handle func(ctx *Context)) {
	api.handler("GET", pattern, handle)
}

//TODO: POST, PUT, DELETE, OPTIONS, HEAD, PATCH

func (api *API) Error(code string, handle func(ctx *Context)) {
	api.errors = append(api.errors, errorHandler{
		code:   code,
		handle: handle,
	})
}
