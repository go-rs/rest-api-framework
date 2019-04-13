/*!
 * rest-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package rest

import (
	"./render"
	"net/http"
)

/**
 * Context
 */
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Params   *map[string]string
	data     map[string]interface{}
	err      error
	status   int
	found    bool
	end      bool
}

/**
 * Initialization of context on every request
 */
func (ctx *Context) init() {
	ctx.data = make(map[string]interface{})
	ctx.status = 200
	ctx.found = false
	ctx.end = false
}

/**
 * Destroy context once request end
 */
func (ctx *Context) destroy() {
	ctx.Request = nil
	ctx.Response = nil
	ctx.Params = nil
	ctx.data = nil
	ctx.err = nil
	ctx.status = 0
	ctx.found = false
	ctx.end = false
}

/**
 * Set request data in context
 */
func (ctx *Context) Set(key string, val interface{}) {
	ctx.data[key] = val
}

/**
 * Get request data from context
 */
func (ctx *Context) Get(key string) interface{} {
	return ctx.data[key]
}

/**
 * Set Status
 */
func (ctx *Context) Status(code int) *Context {
	ctx.status = code
	return ctx
}

/**
 * Set Header
 */
func (ctx *Context) SetHeader(key string, val string) *Context {
	ctx.Response.Header().Set(key, val)
	return ctx
}

/**
 * Throw error
 */
func (ctx *Context) Throw(err error) {
	ctx.err = err
}

/**
 * Get error
 */
func (ctx *Context) GetError() error {
	return ctx.err
}

/**
 * End
 */
func (ctx *Context) End() {
	ctx.end = true
}

/**
 * Write Bytes
 */
func (ctx *Context) Write(data []byte) {
	ctx.send(data, nil)
}

/**
 * Write JSON
 */
func (ctx *Context) JSON(data interface{}) {
	json := render.JSON{
		Body: data,
	}
	body, err := json.Write(ctx.Response)
	ctx.send(body, err)
}

/**
 * Write Text
 */
func (ctx *Context) Text(data string) {
	txt := render.Text{
		Body: data,
	}
	body, err := txt.Write(ctx.Response)
	ctx.send(body, err)
}

//////////////////////////////////////////////////
/**
 * Send data
 */
func (ctx *Context) send(data []byte, err error) {
	if ctx.end && err != nil {
		return
	}
	ctx.Response.WriteHeader(ctx.status)
	_, err = ctx.Response.Write(data)
	if err != nil {
		ctx.err = err
		return
	}

	ctx.End()
}

/**
 * Unhandled Exception
 */
func (ctx *Context) unhandledException() {
	err := ctx.GetError().Error()
	ctx.Status(500)
	if err == "URL_NOT_FOUND" {
		ctx.Status(404)
	}
	ctx.Write([]byte(err))
}
