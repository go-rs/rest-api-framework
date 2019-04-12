/*!
 * utils-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package rest

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Context struct {
	req    *http.Request
	res    http.ResponseWriter
	Params *map[string]string
	data   map[string]interface{}
	err    error
	end    bool
}

func (ctx *Context) init() {
	ctx.data = make(map[string]interface{})
	//TODO: initialization
}

func (ctx *Context) Set(key string, val interface{}) {
	ctx.data[key] = val
}

func (ctx *Context) Get(key string) interface{} {
	return ctx.data[key]
}

func (ctx *Context) Status(code int) *Context {
	ctx.res.WriteHeader(code)
	return ctx
}

func (ctx *Context) SetHeader(key string, val string) *Context {
	ctx.res.Header().Set(key, val)
	return ctx
}

func (ctx *Context) Throw(err error) {
	ctx.err = err
}

func (ctx *Context) Write(data []byte) {
	if ctx.end == true {
		return
	}
	_, err := ctx.res.Write(data)
	if err != nil {
		ctx.err = errors.New("RESPONSE ERROR")
		return
	}

	ctx.end = true
}

func (ctx *Context) Send(data string) {
	ctx.Write([]byte(data))
}

func (ctx *Context) SendJSON(data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		ctx.err = errors.New("INVALID JSON")
		return
	}
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Write(body)
}
