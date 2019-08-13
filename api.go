package rest

import "net/http"

type API struct {
	prefix  string
	list    *list
	Handler Handler
}

func (api *API) Group(prefix string) *API {
	return &API{
		prefix: api.prefix + prefix,
		list:   api.list,
	}
}

func (api *API) Use(task handler) {
	api.list.middleware(task)
}

func (api *API) Get(pattern string, task handler) {
	api.list.route(http.MethodGet, api.prefix+pattern, task)
}

func (api *API) Exception(code string, task handler) {
	api.list.exception(code, task)
}
