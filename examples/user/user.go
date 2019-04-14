package user

import (
	"../.."
)

func APIs(api *rest.API) {

	var user rest.Namespace

	user.Set("/user", api)

	user.Use(func(ctx *rest.Context) {
		println("User middleware > /user/*")
	})

	user.Get("/:uid/profile", func(ctx *rest.Context) {
		ctx.JSON(`{"user": "profile"}`)
	})

	user.Get("/:uid", func(ctx *rest.Context) {
		ctx.JSON(ctx.Params)
	})
}
