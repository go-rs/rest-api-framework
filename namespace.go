// Copyright(c) 2019 Roshan Gade.  All rights reserved.
// MIT Licensed

package rest

import "net/http"

/**
 * Namespace - Application
 */
type Namespace struct {
	prefix string
	api    *API
}

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
	n.api.Route(http.MethodGet, n.prefix+pattern, handle)
}

func (n *Namespace) Post(pattern string, handle Handler) {
	n.api.Route(http.MethodPost, n.prefix+pattern, handle)
}

func (n *Namespace) Put(pattern string, handle Handler) {
	n.api.Route(http.MethodPut, n.prefix+pattern, handle)
}

func (n *Namespace) Delete(pattern string, handle Handler) {
	n.api.Route(http.MethodDelete, n.prefix+pattern, handle)
}

func (n *Namespace) Options(pattern string, handle Handler) {
	n.api.Route(http.MethodOptions, n.prefix+pattern, handle)
}

func (n *Namespace) Head(pattern string, handle Handler) {
	n.api.Route(http.MethodHead, n.prefix+pattern, handle)
}

func (n *Namespace) Patch(pattern string, handle Handler) {
	n.api.Route(http.MethodPatch, n.prefix+pattern, handle)
}

func (n *Namespace) Exception(err string, handle Handler) {
	exp := exception{
		message: err,
		handle:  handle,
	}
	n.api.exceptions = append(n.api.exceptions, exp)
}
