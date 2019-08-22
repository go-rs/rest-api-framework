package rest

import "net/http"

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	// private
	end bool
}

// send JSON
func (ctx *Context) JSON(data string) {
	ctx.write([]byte(data))
}

// send Raw
func (ctx *Context) Raw(data string) {
	ctx.write([]byte(data))
}

// write bytes in response
func (ctx *Context) write(data []byte) {
	ctx.end = true

	// TODO: handle error
	_, _ = ctx.Writer.Write(data)
}
