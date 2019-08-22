package main

import (
	"fmt"
	"net/http"

	"github.com/go-rs/rest-api-framework"
)

func main() {
	var api = rest.New("/")

	api.Use(func(ctx *rest.Context) {
		fmt.Println("/* middleware")
	})

	api.Get("/", func(ctx *rest.Context) {
		ctx.JSON("{\"message\": \"Hello, World!\"}")
	})

	user := api.Group("/user")

	user.Use(func(context *rest.Context) {
		fmt.Println("/user/* middleware")
	})

	user.Get("/", func(ctx *rest.Context) {
		ctx.Raw("Hello, user!")
	})

	http.ListenAndServe(":8080", api)
}
