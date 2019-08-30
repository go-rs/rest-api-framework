package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-rs/rest-api-framework"
)

func main() {
	var api = rest.New("/")

	api.Use(func(ctx rest.Context) {
		fmt.Println("/* middleware")
	})

	api.Get("/", func(ctx rest.Context) {
		ctx.JSON("{\"message\": \"Hello, World!\"}")
	})

	api.Get("/error", func(ctx rest.Context) {
		ctx.Status(500).Throw("TEST_ERROR", errors.New("custom error"))
	})

	api.Exception("TEST_ERROR", func(err error, ctx rest.Context) {
		ctx.Raw(err.Error())
	})

	user := api.Group("/user")

	user.Use(func(context rest.Context) {
		fmt.Println("/user/* middleware")
	})

	user.Get("/", func(ctx rest.Context) {
		ctx.Raw("Hello, user!")
	})

	http.ListenAndServe(":8080", api)
}
