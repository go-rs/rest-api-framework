package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var _api *api
var _list *list
var _handler *handler
var task = func(ctx Context) {}
var errTask = func(err error, ctx Context) {}

func init() {
	_list = new(list)

	_handler = &handler{
		list: _list,
	}

	_api = &api{
		prefix:  trim("/"),
		list:    _list,
		handler: _handler,
	}
}

func TestApi_Use(t *testing.T) {
	_api.Use(task)

	if !(len(_list.middlewares) == 1 && len(_list.routes) == 0 && len(_list.exceptions) == 0 && _list.uncaughtException == nil) {
		t.Error("api.Use should add middleware in the list")
	}

	if _list.middlewares[0].pattern == nil || _list.middlewares[0].task == nil {
		t.Error("api.Use should add type of middleware with pattern and task")
	}
}

func TestApi_Get(t *testing.T) {
	_api.Get("/", task)

	if !(len(_list.middlewares) == 1 && len(_list.routes) == 1 && len(_list.exceptions) == 0 && _list.uncaughtException == nil) {
		t.Error("api.Get should add route in the list")
	}

	if _list.routes[0].method != http.MethodGet || _list.routes[0].pattern == nil || _list.routes[0].task == nil {
		t.Error("api.Get should add type of route with method, pattern and task")
	}
}

func TestApi_Post(t *testing.T) {
	_api.Post("/", task)

	if !(len(_list.middlewares) == 1 && len(_list.routes) == 2 && len(_list.exceptions) == 0 && _list.uncaughtException == nil) {
		t.Error("api.Post should add route in the list")
	}

	if _list.routes[1].method != http.MethodPost || _list.routes[1].pattern == nil || _list.routes[1].task == nil {
		t.Error("api.Post should add type of route with method, pattern and task")
	}
}

func TestApi_Put(t *testing.T) {
	_api.Put("/", task)

	if !(len(_list.middlewares) == 1 && len(_list.routes) == 3 && len(_list.exceptions) == 0 && _list.uncaughtException == nil) {
		t.Error("api.Put should add route in the list")
	}

	if _list.routes[2].method != http.MethodPut || _list.routes[2].pattern == nil || _list.routes[2].task == nil {
		t.Error("api.Put should add type of route with method, pattern and task")
	}
}

func TestApi_Delete(t *testing.T) {
	_api.Delete("/", task)

	if !(len(_list.middlewares) == 1 && len(_list.routes) == 4 && len(_list.exceptions) == 0 && _list.uncaughtException == nil) {
		t.Error("api.Delete should add route in the list")
	}

	if _list.routes[3].method != http.MethodDelete || _list.routes[3].pattern == nil || _list.routes[3].task == nil {
		t.Error("api.Delete should add type of route with method, pattern and task")
	}
}

func TestApi_Group(t *testing.T) {
	user := _api.Group("/user")

	if user == nil {
		t.Error("api.Group should return api reference")
	}
}

func TestApi_Exception(t *testing.T) {
	_api.Exception(ErrCodeNotFound, errTask)

	if !(len(_list.middlewares) == 1 && len(_list.routes) == 4 && len(_list.exceptions) == 1 && _list.uncaughtException == nil) {
		t.Error("api.Exception should add exception handler in the list")
	}

	if _list.exceptions[0].code != ErrCodeNotFound || _list.exceptions[0].task == nil {
		t.Error("api.Exception should add type of exception with code and task")
	}
}

func TestApi_UncaughtException(t *testing.T) {
	_api.UncaughtException(errTask)

	if !(len(_list.middlewares) == 1 && len(_list.routes) == 4 && len(_list.exceptions) == 1 && _list.uncaughtException != nil) {
		t.Error("api.UncaughtException should add uncaught exception handler in the list")
	}

	if _list.uncaughtException == nil {
		t.Error("api.UncaughtException should add type of uncaughtException with task")
	}
}

func TestApi_ServeHTTP(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	_api.ServeHTTP(w, r)

	if w.Result().StatusCode != 200 {
		t.Error("api.ServeHTTP should respond on every request")
	}
}
