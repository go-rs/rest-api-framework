/*!
 * rest-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package rest

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-rs/rest-api-framework/render"
)

type Task func()

/**
 * Context
 */
type Context struct {
	Request         *http.Request
	Response        http.ResponseWriter
	Query           url.Values
	Body            map[string]interface{}
	Params          map[string]string
	headers         map[string]string
	data            map[string]interface{}
	err             error
	status          int
	found           bool
	end             bool
	requestSent     bool
	preTasksCalled  bool
	postTasksCalled bool
	preSendTasks    []Task
	postSendTasks   []Task
}

/**
 * Initialization of context on every request
 */
func (ctx *Context) init() {
	ctx.headers = make(map[string]string)
	ctx.data = make(map[string]interface{})
	ctx.preSendTasks = make([]Task, 0)
	ctx.postSendTasks = make([]Task, 0)
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
func (ctx *Context) Get(key string) (val interface{}, exists bool) {
	val = ctx.data[key]
	exists = val != nil
	return
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
	ctx.headers[key] = val
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
	//ctx.SetHeader("Content-Type", "application/json;charset=UTF-8")
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
	//ctx.SetHeader("Content-Type", "text/plain;charset=UTF-8")
	ctx.send(body, err)
}

/**
 *
 */
func (ctx *Context) PreSend(task Task) {
	ctx.preSendTasks = append(ctx.preSendTasks, task)
}

/**
 *
 */
func (ctx *Context) PostSend(task Task) {
	ctx.postSendTasks = append(ctx.postSendTasks, task)
}

//////////////////////////////////////////////////
func compress(data []byte) (cdata []byte, err error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err = gz.Write(data)
	if err != nil {
		return
	}

	err = gz.Flush()
	if err != nil {
		return
	}

	err = gz.Close()
	if err != nil {
		return
	}

	cdata = b.Bytes()

	return
}

/**
 * Send data
 */
func (ctx *Context) send(data []byte, err error) {
	if ctx.end {
		return
	}

	if err != nil {
		ctx.err = err
		return
	}

	if !ctx.preTasksCalled {
		ctx.preTasksCalled = true
		for _, task := range ctx.preSendTasks {
			task()
		}
	}

	if !ctx.requestSent {
		ctx.requestSent = true

		for key, val := range ctx.headers {
			ctx.Response.Header().Set(key, val)
		}

		if strings.Contains(ctx.Request.Header.Get("Accept-Encoding"), "gzip") {
			data, err = compress(data)
			if err == nil {
				ctx.Response.Header().Set("Content-Encoding", "gzip")
			}
		}

		ctx.Response.WriteHeader(ctx.status)

		_, err = ctx.Response.Write(data)

		if err != nil {
			ctx.err = err
			return
		}
	}

	if !ctx.postTasksCalled {
		ctx.postTasksCalled = true
		for _, task := range ctx.postSendTasks {
			task()
		}
	}

	ctx.End()
}

/**
 * Unhandled Exception
 */
func (ctx *Context) unhandledException() {
	defer func() {
		err := recover()
		if err != nil {
			if !ctx.requestSent {
				ctx.Response.WriteHeader(http.StatusInternalServerError)
				ctx.Response.Header().Set("Content-Type", "text/plain;charset=UTF-8")
				_, _ = ctx.Response.Write([]byte("Internal Server Error"))
			}
		}
	}()

	err := ctx.GetError()

	if err != nil {
		msg := err.Error()
		ctx.Status(http.StatusInternalServerError)
		ctx.SetHeader("Content-Type", "text/plain;charset=UTF-8")
		if msg == "URL_NOT_FOUND" {
			ctx.Status(http.StatusNotFound)
		}
		ctx.Write([]byte(msg))
	}
}
