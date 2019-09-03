package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-rs/rest-api-framework"
)

func main() {
	var api = rest.New("/")

	api.Use(func(ctx rest.Context) {
		fmt.Println("/* middleware")
	})

	api.Get("/page/:id", func(ctx rest.Context) {
		// way to reproduce uncaught exception
		zero, _ := strconv.ParseInt(ctx.Params()["id"], 10, 32)
		x := 10 / zero
		ctx.Raw("This will never respond, if value is zero - " + string(x))
	})

	api.UncaughtException(func(e error, ctx rest.Context) {
		log.Println("ERROR: ", e.Error())
		//zero, _ := strconv.ParseInt(ctx.Params()["id"], 10, 32)
		//_ = 10 / zero
		ctx.Raw("Uncaught exception is handled by user")
	})

	http.ListenAndServe(":8080", api)
}
