package rest

import (
	"net/http"
)

type Group interface {
	Use(Handler)
	Get(string, Handler)
	Post(string, Handler)
	Put(string, Handler)
	Delete(string, Handler)
	Exception(string, ErrorHandler)
}

type Handler func(*Context)

type ErrorHandler func(error, *Context)

type API struct {
	prefix  string
	list    *list
	handler *handler
}

func (api *API) Group(prefix string) Group {
	var group Group = &API{
		prefix: trim(api.prefix + prefix),
		list:   api.list,
	}
	return group
}

func (api *API) Use(task Handler) {
	api.list.middleware(api.prefix, task)
}

func (api *API) Get(pattern string, task Handler) {
	api.list.route(http.MethodGet, api.prefix+pattern, task)
}

func (api *API) Post(pattern string, task Handler) {
	api.list.route(http.MethodGet, api.prefix+pattern, task)
}

func (api *API) Put(pattern string, task Handler) {
	api.list.route(http.MethodGet, api.prefix+pattern, task)
}

func (api *API) Delete(pattern string, task Handler) {
	api.list.route(http.MethodGet, api.prefix+pattern, task)
}

func (api *API) Exception(code string, task ErrorHandler) {
	api.list.exception(code, task)
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.handler.serveHTTP(w, r)
}
