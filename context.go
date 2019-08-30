package rest

import (
	"net/http"
)

type Context interface {
	Request() *http.Request
	Status(int) Context
	Throw(string, error)
	JSON(interface{})
	Raw(interface{})
}

type context struct {
	w http.ResponseWriter
	r *http.Request

	// private
	end    bool
	status int
	code   string
	err    error
}

func (ctx *context) Request() *http.Request {
	return ctx.r
}

func (ctx *context) Status(status int) Context {
	ctx.status = status
	return ctx
}

func (ctx *context) Throw(code string, err error) {
	ctx.code = code
	ctx.err = err
}

// send JSON
func (ctx *context) JSON(data interface{}) {
	// set header
	// TODO:
	ctx.write([]byte(data.(string)))
}

// send Raw
func (ctx *context) Raw(data interface{}) {
	// set header
	// TODO:
	ctx.write([]byte(data.(string)))
}

// write bytes in response
func (ctx *context) write(body []byte) {
	var err error
	ctx.end = true

	if ctx.status > 0 {
		ctx.w.WriteHeader(ctx.status)
	}

	_, err = ctx.w.Write(body)
	if err != nil {
		// TODO: handle error
	}
}
