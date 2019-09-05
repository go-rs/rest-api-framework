package rest

import (
	"log"
	"net/http"
)

type Context interface {
	Request() *http.Request
	Params() map[string]string
	Status(int) Context
	Header(string, string) Context
	Throw(string, error)
	JSON(interface{})
	Raw(interface{})
}

type context struct {
	w http.ResponseWriter
	r *http.Request

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
	ctx.headers = nil
	ctx.params = nil
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

	for key, value := range ctx.headers {
		ctx.w.Header().Set(key, value)
	}

	if ctx.status > 0 {
		ctx.w.WriteHeader(ctx.status)
	}

	_, err = ctx.w.Write(body)
	if err != nil {
		log.Printf("write error: %v", err)
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

		ctx.Header("Content-Type", "text/plain;charset=UTF-8")

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
		log.Printf("runtime error: %v", err)
		if !ctx.end {
			http.Error(ctx.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
