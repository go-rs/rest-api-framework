// go-rs/rest-api-framework
// Copyright(c) 2019 Roshan Gade. All rights reserved.
// MIT Licensed

package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var ctx Context

func TestContext_Set(t *testing.T) {
	ctx.init()
	ctx.Set("foo", "bar")

	if ctx.data["foo"] != "bar" {
		t.Error("Set is not working")
	}
}

func TestContext_Get1(t *testing.T) {
	bar, exists := ctx.Get("foo")

	if !exists && bar != "bar" {
		t.Error("Get is not working")
	}
}

func TestContext_Get2(t *testing.T) {
	_, exists := ctx.Get("bar")

	if exists {
		t.Error("Get is not working")
	}
}

func TestContext_Status(t *testing.T) {
	_ctx := ctx.Status(200)

	if ctx.status != 200 {
		t.Error("Status is not working")
	}

	if _ctx != &ctx {
		t.Error("Status is not returning ctx object")
	}
}

func TestContext_Throw(t *testing.T) {
	err := errors.New("test error")
	ctx.ThrowWithError("TEST_ERROR", err)

	if ctx.err != err {
		t.Error("Throw is not working")
	}

	ctx.err = nil
}

func TestContext_GetError(t *testing.T) {
	err := errors.New("test error")
	ctx.ThrowWithError("TEST_ERROR", err)

	if ctx.GetError() != err {
		t.Error("GetError is not working")
	}

	ctx.err = nil
}

func TestContext_End(t *testing.T) {
	ctx.End()

	if !ctx.end {
		t.Error("End is not working")
	}

	ctx.end = false
}

func TestContext_Write1(t *testing.T) {
	ctx.init()
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	ctx.Response = httptest.NewRecorder()
	data := []byte("Hello World")
	ctx.Write(data)

	if ctx.err != nil {
		t.Error("Write is not working")
	}

	if !ctx.end {
		t.Error("Write is no executing successfully")
	}
	ctx.destroy()
}

func TestContext_Write2(t *testing.T) {
	ctx.init()
	ctx.Response = httptest.NewRecorder()
	ctx.end = true
	data := []byte("Hello World")
	ctx.Write(data)

	if ctx.err != nil {
		t.Error("Write is not working")
	}

	if !ctx.end {
		t.Error("Write is no executing successfully")
	}
	ctx.destroy()
}

func TestContext_JSON1(t *testing.T) {
	ctx.init()
	ctx.Response = httptest.NewRecorder()
	ctx.JSON([]string{"Hello", "World"})

	if ctx.err != nil {
		t.Error("JSON is not working")
	}

	if !ctx.end {
		t.Error("JSON is no executing successfully")
	}
	ctx.destroy()
}

func TestContext_JSON2(t *testing.T) {
	ctx.init()
	ctx.Response = httptest.NewRecorder()
	ctx.JSON("Hello World")

	if ctx.err.Error() != "INVALID_JSON_RESPONSE" {
		t.Error("JSON is not working")
	}

	if !ctx.end {
		t.Error("JSON is no executing successfully")
	}
	ctx.destroy()
}

func TestContext_Text(t *testing.T) {
	ctx.init()
	ctx.Response = httptest.NewRecorder()
	ctx.Text("Hello World")

	if ctx.err != nil {
		t.Error("Text is not working")
	}

	if !ctx.end {
		t.Error("Text is no executing successfully")
	}
	ctx.destroy()
}
