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
		ctx.JSON(true)
	})

	api.Get("/", func(ctx rest.Context) {
		ctx.JSON("{\"message\": \"Hello, World!\"}")
	})

	//==========================Error Handling==================================

	api.Get("/error", func(ctx rest.Context) {
		ctx.Status(500).Throw("TEST_ERROR", errors.New("custom error"))
	})

	api.Exception("TEST_ERROR", func(err error, ctx rest.Context) {
		ctx.Text(err.Error())
	})

	//============================Extended routes====================================

	user := api.Group("/user")

	user.Use(func(ctx rest.Context) {
		fmt.Println("/user/* middleware")
	})

	user.Get("/", func(ctx rest.Context) {
		ctx.XML(`
			<person id="13">
				  <name>
					  <first>John</first>
					  <last>Doe</last>
				  </name>
				  <age>42</age>
				  <Married>false</Married>
				  <City>Hanga Roa</City>
				  <State>Easter Island</State>
			  </person>
		`)
	})

	http.ListenAndServe(":8080", api)
}
