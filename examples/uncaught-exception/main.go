package main

import (
	"net/http"
	"strconv"

	"github.com/go-rs/rest-api-framework"
)

func main() {
	var api = rest.New("/")

	api.Get("/page/:id", func(ctx rest.Context) {
		// way to reproduce uncaught exception
		zero, _ := strconv.ParseInt(ctx.Params()["id"], 10, 32)
		x := 10 / zero
		ctx.Text("This will never respond, if value is zero - " + string(x))
	})

	api.OnUncaughtException(func(e error, ctx rest.Context) {
		ctx.Status(500).Text("Uncaught exception is handled by user")
	})

	http.ListenAndServe(":8080", api)
}
