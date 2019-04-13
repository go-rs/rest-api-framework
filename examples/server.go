package main

import (
	".."
	"errors"
	"fmt"
	"net/http"
)

func main() {
	var api rest.API

	// request interceptor / middleware
	// method-override
	// body-parser : json, raw, form-data, etc
	// security
	api.Use(func(ctx *rest.Context) {
		ctx.Set("authtoken", "roshangade")
	})

	// routes
	api.GET("/", func(ctx *rest.Context) {
		ctx.Text("Hello World!")
	})

	api.GET("/foo", func(ctx *rest.Context) {
		ctx.Throw(errors.New("UNAUTHORIZED"))
	})

	api.GET("/strut", func(ctx *rest.Context) {
		type user struct {
			Name string
		}

		u := user{
			Name: "roshan",
		}

		ctx.JSON(u)
	})

	api.GET("/json", func(ctx *rest.Context) {
		ctx.JSON(`{"test": "JSON Response"}`)
	})

	api.GET("/:bar", func(ctx *rest.Context) {
		fmt.Println("authtoken", ctx.Get("authtoken"))
		ctx.JSON(ctx.Params)
	})

	// error handler
	api.Exception("UNAUTHORIZED", func(ctx *rest.Context) {
		ctx.Status(401).Text("You are unauthorized")
	})

	api.UnhandledException(func(ctx *rest.Context) {
		fmt.Println("---> ", ctx.GetError())
	})

	fmt.Println("Starting server.")

	http.ListenAndServe(":8080", api)
}
