// go-rs/rest-api-framework
// Copyright(c) 2019-2022 Roshan Gade. All rights reserved.
// MIT Licensed
package rest

import (
	"net/http"
)

type Handler func(Context)

type ErrorHandler func(error, Context)

type Router interface {
	Use(Handler)
	Router(string) Router
	Get(string, Handler)
	Post(string, Handler)
	Put(string, Handler)
	Delete(string, Handler)
	CatchError(string, ErrorHandler)
}

type API interface {
	Router
	UncaughtException(ErrorHandler)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type api struct {
	prefix         string
	router         *router
	requestHandler *requestHandler
}

func (a *api) Router(prefix string) Router {
	var router Router = &api{
		prefix: trim(a.prefix + prefix),
		router: a.router,
	}
	return router
}

func (a *api) Use(task Handler) {
	a.router.middleware(a.prefix, task)
}

func (a *api) Get(pattern string, task Handler) {
	a.router.route(http.MethodGet, a.prefix+pattern, task)
}

func (a *api) Post(pattern string, task Handler) {
	a.router.route(http.MethodPost, a.prefix+pattern, task)
}

func (a *api) Put(pattern string, task Handler) {
	a.router.route(http.MethodPut, a.prefix+pattern, task)
}

func (a *api) Delete(pattern string, task Handler) {
	a.router.route(http.MethodDelete, a.prefix+pattern, task)
}

func (a *api) CatchError(code string, task ErrorHandler) {
	a.router.exception(code, task)
}

func (a *api) UncaughtException(task ErrorHandler) {
	a.router.unhandledException(task)
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.requestHandler.serveHTTP(w, r)
}
