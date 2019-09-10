package rest

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var ctx *context

func init() {
	ctx = &context{
		r: httptest.NewRequest(http.MethodGet, "/", nil),
		w: httptest.NewRecorder(),
	}
}

func TestContext_init(t *testing.T) {
	ctx.init()

	if ctx.headers == nil {
		t.Error("context.init() should initialize the context")
	}
}

func TestContext_reset(t *testing.T) {
	ctx.reset()

	if !reflect.DeepEqual(ctx.headers, make(map[string]string)) {
		t.Error("context.reset() should reset the context")
	}
}

func TestContext_Request(t *testing.T) {
	if ctx.Request() != ctx.r {
		t.Error("context.Request() should return a pointer of request")
	}
}

func TestContext_Status(t *testing.T) {
	status := 200
	c := ctx.Status(status)

	if c != ctx {
		t.Error("context.Status(int) should return pointer of context")
	}

	if ctx.status != status {
		t.Error("context.Status(int) should set a status")
	}
}

func TestContext_Params(t *testing.T) {
	ctx.params = make(map[string]string)
	params := ctx.Params()
	if !reflect.DeepEqual(ctx.params, params) {
		t.Error("context.Params() should return the object of request params")
	}
}

func TestContext_Throw(t *testing.T) {
	ctx.Status(504).Throw("TIMEOUT", errors.New("request timeout"))

	if ctx.code != "TIMEOUT" || !reflect.DeepEqual(ctx.err, errors.New("request timeout")) {
		t.Error("context.Throw(string, error) should throw an error with error code")
	}
}

func TestContext_JSON(t *testing.T) {
	res := "{\"message\":\"Hello, World!\"}"
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ctx = &context{
		r: r,
		w: w,
	}
	ctx.init()
	defer ctx.destroy()

	ctx.JSON(res)

	body, err := ioutil.ReadAll(w.Body)

	if err != nil || string(body) != res {
		t.Error("context.JSON() should write a response a JSON data in request")
	}

	if ctx.headers["Content-Type"] != headerJSON {
		t.Error("context.JSON() should set a response header with Content-Type: " + headerJSON)
	}
}

func TestContext_destroy(t *testing.T) {
	ctx.destroy()

	if ctx.w != nil || ctx.r != nil || ctx.headers != nil {
		t.Error("context.destroy() should destroy the context")
	}
}
