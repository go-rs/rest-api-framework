// go-rs/rest-api-framework
// Copyright(c) 2019-2020 Roshan Gade. All rights reserved.
// MIT Licensed
package rest

import (
	"log"
	"net/http"
	"sync"
)

type Context interface {
	Request() *http.Request
	Params() map[string]string
	Status(int) Context
	Header(string, string) Context
	Set(string, any)
	Get(string) (any, bool)
	Unset(string)
	Throw(string, error)
	JSON(any)
	XML(any)
	Text(string)
	Write([]byte)
	End()
}

const (
	ErrCodeInvalidJSON = "INVALID_JSON"
	ErrCodeInvalidXML  = "INVALID_XML"
)

const HeaderContentType = "Content-Type"

const (
	headerText = "text/plain"
	headerJSON = "application/json"
	headerXML  = "text/xml"
)

type context struct {
	w    http.ResponseWriter
	r    *http.Request
	meta sync.Map

	// for internal purpose
	params  map[string]string
	headers map[string]string
	end     bool
	status  int
	code    string
	err     error
}

func (ctx *context) init() {
	ctx.headers = make(map[string]string)
}

func (ctx *context) destroy() {
	ctx.w = nil
	ctx.r = nil
	ctx.meta = sync.Map{}
	ctx.params = nil
	ctx.headers = nil
	ctx.end = true
	ctx.status = 0
	ctx.code = ""
	ctx.err = nil
}

func (ctx *context) reset() {
	ctx.headers = make(map[string]string)
}

func (ctx *context) Request() *http.Request {
	return ctx.r
}

func (ctx *context) Params() map[string]string {
	return ctx.params
}

func (ctx *context) Status(status int) Context {
	ctx.status = status
	return ctx
}

func (ctx *context) Header(key string, value string) Context {
	ctx.headers[key] = value
	return ctx
}

func (ctx *context) Set(key string, value any) {
	ctx.meta.Store(key, value)
}

func (ctx *context) Get(key string) (any, bool) {
	return ctx.meta.Load(key)
}

func (ctx *context) Unset(key string) {
	ctx.meta.Delete(key)
}

func (ctx *context) Throw(code string, err error) {
	ctx.code = code
	ctx.err = err
}

// send JSON
func (ctx *context) JSON(data any) {
	body, err := jsonToBytes(data)
	if err != nil {
		ctx.Throw(ErrCodeInvalidJSON, err)
		return
	}
	if _, exists := ctx.headers[HeaderContentType]; !exists {
		ctx.Header(HeaderContentType, headerJSON)
	}
	ctx.write(body)
}

func (ctx *context) XML(data any) {
	body, err := xmlToBytes(data)
	if err != nil {
		ctx.Throw(ErrCodeInvalidXML, err)
		return
	}
	if _, exists := ctx.headers[HeaderContentType]; !exists {
		ctx.Header(HeaderContentType, headerXML)
	}
	ctx.write(body)
}

// send Raw
func (ctx *context) Text(data string) {
	if _, exists := ctx.headers[HeaderContentType]; !exists {
		ctx.Header(HeaderContentType, headerText)
	}
	ctx.write([]byte(data))
}

func (ctx *context) Write(data []byte) {
	ctx.write(data)
}

// send blank
func (ctx *context) End() {
	ctx.write(nil)
}

// write bytes in response
func (ctx *context) write(body []byte) {

	if ctx.end {
		log.Printf("http/request: trying to write response data on already ended request.")
		return
	}

	var err error
	ctx.end = true

	for key, value := range ctx.headers {
		ctx.w.Header().Set(key, value)
	}

	if ctx.status > 0 {
		ctx.w.WriteHeader(ctx.status)
	}

	_, err = ctx.w.Write(body)
	if err != nil {
		log.Printf("http/request: write error - %v", err)
	}
}

// unhandled exception
func (ctx *context) unhandledException() {
	defer ctx.recover()

	if ctx.end {
		return
	}

	if ctx.err != nil {
		ctx.reset()

		ctx.Header(HeaderContentType, headerText)

		if ctx.code == ErrCodeNotFound {
			ctx.Status(http.StatusNotFound)
		} else if ctx.status < 400 {
			ctx.Status(http.StatusInternalServerError)
		}

		ctx.write([]byte(ctx.err.Error()))
	}
}

// recover
func (ctx *context) recover() {
	err := recover()
	if err != nil {
		log.Printf("http/request: runtime error - %v", err)
		if !ctx.end {
			http.Error(ctx.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
