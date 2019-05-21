/*!
 * rest-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */

package rest

import (
	"testing"
)

var api API
var ns Namespace
var handle Handler

func TestNamespace_Set(t *testing.T) {
	ns.Set("/test", &api)

	if ns.prefix != "/test" {
		t.Error("Prefix is not set.")
	}
}

func validateRoute(fun string, method string, url string, t *testing.T) {
	flag := true
	for _, route := range api.routes {
		if route.method == method && route.pattern == url {
			flag = false
			break
		}
	}

	if flag {
		t.Error("Namespace " + fun + " is not working properly")
	}
}

func TestNamespace_Use(t *testing.T) {
	ns.Use(handle)

	validateRoute("Use", "", "/test/*", t)
}

func TestNamespace_All(t *testing.T) {
	ns.All("/:uid", handle)

	validateRoute("All", "", "/test/:uid", t)
}

func TestNamespace_Get(t *testing.T) {
	ns.Get("/:uid", handle)

	validateRoute("Get", "GET", "/test/:uid", t)
}

func TestNamespace_Post(t *testing.T) {
	ns.Post("/:uid", handle)

	validateRoute("Post", "POST", "/test/:uid", t)
}

func TestNamespace_Put(t *testing.T) {
	ns.Put("/:uid", handle)

	validateRoute("Put", "PUT", "/test/:uid", t)
}

func TestNamespace_Delete(t *testing.T) {
	ns.Delete("/:uid", handle)

	validateRoute("Delete", "DELETE", "/test/:uid", t)
}

func TestNamespace_Options(t *testing.T) {
	ns.Options("/:uid", handle)

	validateRoute("Options", "OPTIONS", "/test/:uid", t)
}

func TestNamespace_Head(t *testing.T) {
	ns.Head("/:uid", handle)

	validateRoute("Head", "HEAD", "/test/:uid", t)
}

func TestNamespace_Patch(t *testing.T) {
	ns.Patch("/:uid", handle)

	validateRoute("Patch", "PATCH", "/test/:uid", t)
}

func TestNamespace_Exception(t *testing.T) {
	ns.Exception("UID_NOT_FOUND", handle)

	flag := true
	for _, route := range api.exceptions {
		if route.message == "UID_NOT_FOUND" {
			flag = false
			break
		}
	}

	if flag {
		t.Error("Namespace Exception is not working properly")
	}
}
