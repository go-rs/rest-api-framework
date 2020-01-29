// go-rs/rest-api-framework
// Copyright(c) 2019-2020 Roshan Gade. All rights reserved.
// MIT Licensed
package rest

import "log"

type middleware struct {
	pattern *pattern
	task    Handler
}

type route struct {
	method  string
	pattern *pattern
	task    Handler
}

type exception struct {
	code string
	task ErrorHandler
}

type router struct {
	middlewares       []middleware
	routes            []route
	exceptions        []exception
	uncaughtException ErrorHandler
}

func (r *router) middleware(str string, task Handler) {
	p := &pattern{
		value: trim(str) + "/*",
	}
	if err := p.compile(); err != nil {
		log.Fatalf("Failed to compile `%s` due to %v", p.value, err)
	}
	r.middlewares = append(r.middlewares, middleware{pattern: p, task: task})
}

func (r *router) route(method string, str string, task Handler) {
	p := &pattern{
		value: trim(str),
	}
	if err := p.compile(); err != nil {
		log.Fatalf("Failed to compile `%s` due to %v", p.value, err)
	}
	r.routes = append(r.routes, route{
		method:  method,
		pattern: p,
		task:    task,
	})
}

func (r *router) exception(code string, task ErrorHandler) {
	r.exceptions = append(r.exceptions, exception{code: code, task: task})
}

func (r *router) unhandledException(task ErrorHandler) {
	r.uncaughtException = task
}
