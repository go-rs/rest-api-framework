package main

import (
	"github.com/go-rs/rest-api-framework"
	"errors"
	"fmt"
	"github.com/go-rs/rest-api-framework/examples/user"
	"net/http"
)

func main() {
	var api rest.API

	user.APIs(&api)

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
	api.Exception("UNAUTHORIZED", func(ctx *rest.Context) {
		ctx.Status(401).JSON(`{"message": "You are unauthorized"}`)
	})

	fmt.Println("Starting server.")

	http.ListenAndServe(":8080", api)
}
