// go-rs/rest-api-framework
// Copyright(c) 2019 Roshan Gade. All rights reserved.
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

type list struct {
	middlewares       []middleware
	routes            []route
	exceptions        []exception
	uncaughtException ErrorHandler
}

func (l *list) middleware(str string, task Handler) {
	p := &pattern{
		value: trim(str) + "/*",
	}
	if err := p.compile(); err != nil {
		log.Fatalf("Failed to compile `%s` due to %v", p.value, err)
	}
	l.middlewares = append(l.middlewares, middleware{pattern: p, task: task})
}

func (l *list) route(method string, str string, task Handler) {
	p := &pattern{
		value: trim(str),
	}
	if err := p.compile(); err != nil {
		log.Fatalf("Failed to compile `%s` due to %v", p.value, err)
	}
	l.routes = append(l.routes, route{
		method:  method,
		pattern: p,
		task:    task,
	})
}

func (l *list) exception(code string, task ErrorHandler) {
	l.exceptions = append(l.exceptions, exception{code: code, task: task})
}

func (l *list) unhandledException(task ErrorHandler) {
	l.uncaughtException = task
}
