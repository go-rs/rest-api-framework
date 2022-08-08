// go-rs/rest-api-framework
// Copyright(c) 2019-2022 Roshan Gade. All rights reserved.
// MIT Licensed
package rest

func New(prefix string) API {
	var _router = new(router)

	var _requestHandler = &requestHandler{
		router: _router,
	}

	var api = &api{
		prefix:         trim(prefix),
		router:         _router,
		requestHandler: _requestHandler,
	}

	return api
}
