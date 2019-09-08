// go-rs/rest-api-framework
// Copyright(c) 2019 Roshan Gade. All rights reserved.
// MIT Licensed
package rest

import (
	"net/http"
)

type Handler func(Context)

type ErrorHandler func(error, Context)

type API interface {
	Use(Handler)
	Group(string) Group
	Get(string, Handler)
	Post(string, Handler)
	Put(string, Handler)
	Delete(string, Handler)
	Exception(string, ErrorHandler)
	UncaughtException(ErrorHandler)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Group interface {
	Use(Handler)
	Group(string) Group
	Get(string, Handler)
	Post(string, Handler)
	Put(string, Handler)
	Delete(string, Handler)
	Exception(string, ErrorHandler)
}

type api struct {
	prefix  string
	list    *list
	handler *handler
}

func (a *api) Group(prefix string) Group {
	var group Group = &api{
		prefix: trim(a.prefix + prefix),
		list:   a.list,
	}
	return group
}

func (a *api) Use(task Handler) {
	a.list.middleware(a.prefix, task)
}

func (a *api) Get(pattern string, task Handler) {
	a.list.route(http.MethodGet, a.prefix+pattern, task)
}

func (a *api) Post(pattern string, task Handler) {
	a.list.route(http.MethodPost, a.prefix+pattern, task)
}

func (a *api) Put(pattern string, task Handler) {
	a.list.route(http.MethodPut, a.prefix+pattern, task)
}

func (a *api) Delete(pattern string, task Handler) {
	a.list.route(http.MethodDelete, a.prefix+pattern, task)
}

func (a *api) Exception(code string, task ErrorHandler) {
	a.list.exception(code, task)
}

func (a *api) UncaughtException(task ErrorHandler) {
	a.list.unhandledException(task)
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.serveHTTP(w, r)
}
