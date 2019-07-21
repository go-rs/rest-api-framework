// go-rs/rest-api-framework
// Copyright(c) 2019 Roshan Gade. All rights reserved.
// MIT Licensed

package rest

import (
	"net/http"

	"github.com/go-rs/rest-api-framework/utils"
)

// Namespace is used to extend routes with a specific prefix
type Namespace struct {
	prefix string
	api    *API
}

// Set method is used to set prefix with API
func Extend(prefix string, api *API) *Namespace {
	prefix = utils.TrimSuffix(prefix, "/")
	return &Namespace{
		prefix: prefix,
		api:    api,
	}
}

// Use method is used to set interceptors/middlewares for all requests,
// which are going to travel through prefix
func (n *Namespace) Use(handle Handler) {
	n.api.Route("", n.prefix+"/*", handle)
}

// All method is used to execute for all methods,
// but with some more prefix as compare to Use method
func (n *Namespace) All(pattern string, handle Handler) {
	n.api.Route("", n.prefix+pattern, handle)
}

// Get method is used for GET http method with specific pattern
func (n *Namespace) Get(pattern string, handle Handler) {
	n.api.Route(http.MethodGet, n.prefix+pattern, handle)
}

// Post method is used for POST http method with specific pattern
func (n *Namespace) Post(pattern string, handle Handler) {
	n.api.Route(http.MethodPost, n.prefix+pattern, handle)
}

// Put method is used for PUT http method with specific pattern
func (n *Namespace) Put(pattern string, handle Handler) {
	n.api.Route(http.MethodPut, n.prefix+pattern, handle)
}

// Delete method is used for DELETE http method with specific pattern
func (n *Namespace) Delete(pattern string, handle Handler) {
	n.api.Route(http.MethodDelete, n.prefix+pattern, handle)
}

// Options method is used for OPTIONS http method with specific pattern
func (n *Namespace) Options(pattern string, handle Handler) {
	n.api.Route(http.MethodOptions, n.prefix+pattern, handle)
}

// Head method is used for HEAD http method with specific pattern
func (n *Namespace) Head(pattern string, handle Handler) {
	n.api.Route(http.MethodHead, n.prefix+pattern, handle)
}

// Patch method is used for PATCH http method with specific pattern
func (n *Namespace) Patch(pattern string, handle Handler) {
	n.api.Route(http.MethodPatch, n.prefix+pattern, handle)
}

// OnError method is an exactly same with API.OnError
func (n *Namespace) OnError(code string, handle Handler) {
	exp := exception{
		code:   code,
		handle: handle,
	}
	n.api.exceptions = append(n.api.exceptions, exp)
}
