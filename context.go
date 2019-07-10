// go-rs/rest-api-framework
// Copyright(c) 2019 Roshan Gade. All rights reserved.
// MIT Licensed

package rest

import (
	"log"
	"net/http"
	"net/url"

	"github.com/go-rs/rest-api-framework/render"
)

// Task, which is used to perform a specific job,
// just before request completion or after request completion
type Task func()

// Context, which initializes at every request with pre-declared variables
// such as Request, Response, Query, Body, Params etc.
type Context struct {
	// available to users
	Request  *http.Request
	Response http.ResponseWriter
	Query    url.Values
	Body     interface{}
	Params   map[string]string

	// for internal use
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

// Initialization of context data on every request
func (ctx *Context) init() {
	ctx.headers = make(map[string]string)
	ctx.data = make(map[string]interface{})
	ctx.preSendTasks = make([]Task, 0)
	ctx.postSendTasks = make([]Task, 0)
	ctx.status = 200
}

// Destroy context data once request end
func (ctx *Context) destroy() {
	ctx.Request = nil
	ctx.Response = nil
	ctx.Query = nil
	ctx.Body = nil
	ctx.Params = nil
	ctx.headers = nil
	ctx.data = nil
	ctx.err = nil
	ctx.status = 0
	ctx.found = false
	ctx.end = false
	ctx.requestSent = false
	ctx.preTasksCalled = false
	ctx.postTasksCalled = false
	ctx.preSendTasks = nil
	ctx.postSendTasks = nil
}

// Set data
func (ctx *Context) Set(key string, val interface{}) {
	ctx.data[key] = val
}

// Get data
func (ctx *Context) Get(key string) (val interface{}, exists bool) {
	val, exists = ctx.data[key]
	return
}

// Delete data
func (ctx *Context) Delete(key string, val interface{}) {
	delete(ctx.data, key)
}

// Set response status
func (ctx *Context) Status(code int) *Context {
	ctx.status = code
	return ctx
}

// Set response header
func (ctx *Context) SetHeader(key string, val string) *Context {
	ctx.headers[key] = val
	return ctx
}

// Caught error in context on throw
func (ctx *Context) Throw(err error) {
	ctx.err = err
}

// Get error if any
func (ctx *Context) GetError() error {
	return ctx.err
}

// Marked the request is ended
func (ctx *Context) End() {
	ctx.end = true
}

// Send response in bytes
func (ctx *Context) Write(data []byte) {
	ctx.send(data, nil)
}

// Send JSON data in response
func (ctx *Context) JSON(data interface{}) {
	json := render.JSON{
		Body: data,
	}
	body, err := json.ToBytes(ctx.Response)
	ctx.send(body, err)
}

// Send text in response
func (ctx *Context) Text(data string) {
	txt := render.Text{
		Body: data,
	}
	body, err := txt.ToBytes(ctx.Response)
	ctx.send(body, err)
}

// Register pre-send hook
func (ctx *Context) PreSend(task Task) {
	ctx.preSendTasks = append(ctx.preSendTasks, task)
}

// Register pre-post hook
func (ctx *Context) PostSend(task Task) {
	ctx.postSendTasks = append(ctx.postSendTasks, task)
}

//////////////////////////////////////////////////
// Send data, which uses bytes or error if any
// Also, it calls pre-send and post-send registered hooks
func (ctx *Context) send(data []byte, err error) {
	if ctx.end {
		return
	}

	if err != nil {
		ctx.err = err
		ctx.unhandledException()
		return
	}

	// execute pre-send hooks
	if !ctx.preTasksCalled {
		ctx.preTasksCalled = true
		for _, task := range ctx.preSendTasks {
			task()
		}
	}

	// write data
	if !ctx.requestSent {
		ctx.requestSent = true

		for key, val := range ctx.headers {
			ctx.Response.Header().Set(key, val)
		}

		ctx.Response.WriteHeader(ctx.status)

		_, err = ctx.Response.Write(data)

		if err != nil {
			//TODO: debugger mode
			log.Println("Response Error: ", err)
		}
	}

	// execute post-send hooks
	if !ctx.postTasksCalled {
		ctx.postTasksCalled = true
		for _, task := range ctx.postSendTasks {
			task()
		}
	}

	ctx.End()
}

// Unhandled Exception
func (ctx *Context) unhandledException() {
	defer func() {
		err := recover()
		if err != nil {
			//TODO: debugger mode
			log.Println("Unhandled Error: ", err)
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
