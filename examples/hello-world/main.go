package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-rs/rest-api-framework"
)

func main() {
	var api = rest.New("/")

	// global request middleware
	api.Use(func(ctx rest.Context) {
		fmt.Println("/* middleware")
	})

	// curl -X GET 'localhost:8080'
	api.Get("/", func(ctx rest.Context) {
		ctx.Status(200).JSON("{\"message\": \"Hello, World!\"}")
	})

	//============================Extended routes====================================

	user := api.Router("/user")

	// All /user/* request middleware
	user.Use(func(ctx rest.Context) {
		fmt.Println("/user/* middleware")
	})

	// curl -X GET 'localhost:8080/user'
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

	// User session
	session := user.Router("/session")

	// 	curl -X POST 'localhost:8080/user/session'
	session.Post("/", func(ctx rest.Context) {
		ctx.Text("Login successful")
	})

	// 	curl -X DELETE 'localhost:8080/user/session'
	session.Delete("/", func(ctx rest.Context) {
		ctx.Status(204).End()
	})

	//==========================Error Handling==================================

	api.Get("/test", func(ctx rest.Context) {
		ctx.Set("Hi", "There")
		ctx.Throw("CUSTOM_ERROR", errors.New("custom error"))
	})

	api.CatchError("CUSTOM_ERROR", func(err error, ctx rest.Context) {
		//ctx.Status(412).JSON("{\"a\": 1}")
		ctx.Status(412).JSON(ctx.Metadata())
	})

	http.ListenAndServe(":8080", api)
}
