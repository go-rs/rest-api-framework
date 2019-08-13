package main

import (
	"fmt"
	"github.com/go-rs/rest-api-framework"
	"net/http"
)

func main() {
	var api, handler = rest.New("/")

	api.Get("/", func() {
		fmt.Println("Hello, World!")
	})

	user := api.Group("/user")

	user.Get("/", func() {
		fmt.Println("Hello, User!")
	})

	_ = http.ListenAndServe(":8080", handler)

	fmt.Println(api)
}
