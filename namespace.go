/*!
 * rest-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package rest

/**
 * Namespace - Application
 */
type Namespace struct {
	prefix string
	api    *API
}

//TODO: error handling on unset api
func (n *Namespace) Set(prefix string, api *API) {
	n.prefix = prefix
	n.api = api
}

func (n *Namespace) Use(handle Handler) {
	n.api.Route("", n.prefix+"/*", handle)
}

func (n *Namespace) All(pattern string, handle Handler) {
	n.api.Route("", n.prefix+pattern, handle)
}

func (n *Namespace) Get(pattern string, handle Handler) {
	n.api.Route("GET", n.prefix+pattern, handle)
}

func (n *Namespace) Post(pattern string, handle Handler) {
	n.api.Route("POST", n.prefix+pattern, handle)
}

func (n *Namespace) Put(pattern string, handle Handler) {
	n.api.Route("PUT", n.prefix+pattern, handle)
}

func (n *Namespace) Delete(pattern string, handle Handler) {
	n.api.Route("DELETE", n.prefix+pattern, handle)
}

func (n *Namespace) Options(pattern string, handle Handler) {
	n.api.Route("OPTIONS", n.prefix+pattern, handle)
}

func (n *Namespace) Head(pattern string, handle Handler) {
	n.api.Route("HEAD", n.prefix+pattern, handle)
}

func (n *Namespace) Patch(pattern string, handle Handler) {
	n.api.Route("PATCH", n.prefix+pattern, handle)
}

func (n *Namespace) Exception(err string, handle Handler) {
	exp := exception{
		message: err,
		handle:  handle,
	}
	n.api.exceptions = append(n.api.exceptions, exp)
}
