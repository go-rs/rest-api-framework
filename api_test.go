/*!
 * rest-a-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package rest

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a API
var h Handler

func validateRoute(fun string, method string, url string, t *testing.T) {
	flag := true
	for _, route := range a.routes {
		if route.method == method && route.pattern == url {
			flag = false
			break
		}
	}

	if flag {
		t.Error("API: " + fun + " is not working properly")
	}
}

func TestAPI_Route(t *testing.T) {
	a.Route("GET", "/greeting", handle)

	validateRoute("Route", "GET", "/greeting", t)
}

func TestAPI_Use(t *testing.T) {
	a.Use(handle)

	if len(a.interceptors) == 0 {
		t.Error("API: Use is not working properly")
	}
}

func TestAPI_All(t *testing.T) {
	a.All("/:uid", handle)

	validateRoute("All", "", "/:uid", t)
}

func TestAPI_Get(t *testing.T) {
	a.Get("/:uid", handle)

	validateRoute("Get", "GET", "/:uid", t)
}

func TestAPI_Post(t *testing.T) {
	a.Post("/:uid", handle)

	validateRoute("Post", "POST", "/:uid", t)
}

func TestAPI_Put(t *testing.T) {
	a.Put("/:uid", handle)

	validateRoute("Put", "PUT", "/:uid", t)
}

func TestAPI_Delete(t *testing.T) {
	a.Delete("/:uid", handle)

	validateRoute("Delete", "DELETE", "/:uid", t)
}

func TestAPI_Options(t *testing.T) {
	a.Options("/:uid", handle)

	validateRoute("Optioa", "OPTIONS", "/:uid", t)
}

func TestAPI_Head(t *testing.T) {
	a.Head("/:uid", handle)

	validateRoute("Head", "HEAD", "/:uid", t)
}

func TestAPI_Patch(t *testing.T) {
	a.Patch("/:uid", handle)

	validateRoute("Patch", "PATCH", "/:uid", t)
}

func TestAPI_Exception(t *testing.T) {
	a.Exception("UID_NOT_FOUND", handle)

	flag := true
	for _, route := range a.exceptions {
		if route.message == "UID_NOT_FOUND" {
			flag = false
			break
		}
	}

	if flag {
		t.Error("API: Exception is not working properly")
	}
}

func TestAPI_ServeHTTP(t *testing.T) {
	var _api API

	_api.Get("/", func(ctx *Context) {
		ctx.JSON(`{"message": "Hello World!"}`)
	})

	dummy := httptest.NewServer(_api)
	defer dummy.Close()

	res, err := http.Get(dummy.URL)

	if err != nil {
		t.Error("ServeHTTP error")
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error("ServeHTTP error")
	}

	if string(greeting) != `{"message": "Hello World!"}` {
		t.Error("Response does not match")
	}
}
