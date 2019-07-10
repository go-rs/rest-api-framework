package main

import (
	"errors"
	"fmt"
	"net/http"

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
		ctx.Set("authtoken", "roshangade")
	})

	// routes
	api.Get("/", func(ctx *rest.Context) {
		ctx.JSON(`{"message": "Hello World!"}`)
	})

	api.Get("/foo", func(ctx *rest.Context) {
		ctx.Throw(errors.New("UNAUTHORIZED"))
	})

	// error handler
	api.OnError("UNAUTHORIZED", func(ctx *rest.Context) {
		ctx.Status(401).JSON(`{"message": "You are unauthorized"}`)
	})

	fmt.Println("Starting server.")

	http.ListenAndServe(":8080", api)
}
