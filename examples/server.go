package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-rs/rest-api-framework"
	"github.com/go-rs/rest-api-framework/examples/user"
)

func main() {

	var api = rest.New("/v1")

	user.Load(api)

	// request interceptor / middleware
	// body-parser : json, raw, form-data, etc
	// security
	api.Use(func(ctx *rest.Context) {
		// way to reproduce uncaught exception
		//zero, _ := strconv.ParseInt(ctx.Query.Get("zero"), 10, 32)
		//x := 10/zero
		//ctx.Set("x", x)
		ctx.Set("authtoken", "roshangade")
	})

	// calculate runtime
	api.Use(func(ctx *rest.Context) {
		s := time.Now().UnixNano()
		ctx.PreSend(func() {
			x := time.Now().UnixNano() - s
			ctx.SetHeader("X-Runtime", strconv.FormatInt(x/int64(time.Microsecond), 10))
		})
	})

	// routes
	api.Get("/", func(ctx *rest.Context) {
		ctx.JSON(`{"message": "Hello World!"}`)
	})

	api.Get("/foo", func(ctx *rest.Context) {
		ctx.Throw("UNAUTHORIZED")
	})

	// error handler
	api.OnError("UNAUTHORIZED", func(ctx *rest.Context) {
		ctx.Status(401).JSON(`{"message": "You are unauthorized"}`)
	})

	fmt.Println("Starting server.")

	server := http.Server{
		Addr:    ":8080",
		Handler: api,
	}

	err := server.ListenAndServe()
	fmt.Println(err)
}
