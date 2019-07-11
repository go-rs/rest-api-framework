# REST API Framework
REST API framework for go lang

# Framework is under development
## Status: 
Released beta version
<br>
See examples 
  - Request Interceptors/Middleware
  - Routes with URL pattern 
  - Methods [GET, POST, PUT, DELETE, OPTIONS, HEAD, PATCH]
  - Extend routes with namespace
  - Error handler
  - HTTP, HTTPS support
  
## How to use?  
  
```
var api = rest.New("/v1")

// request interceptor / middleware
api.Use(func(ctx *rest.Context) {
  ctx.Set("authtoken", "roshangade")
})

// hooks support on request pre-send or post-send
// ctx.PreSend(func() error) OR ctx.PostSend(func() error)
api.Use(func(ctx *rest.Context) {
  s := time.Now().UnixNano()
  ctx.PreSend(func() {
    x := time.Now().UnixNano() - s
    ctx.SetHeader("X-Runtime", strconv.FormatInt(x/int64(time.Millisecond), 10))
    return nil
  })
})


// routes
api.Get("/", func(ctx *rest.Context) {
  ctx.Text("Hello World!")
})

api.Get("/foo", func(ctx *rest.Context) {
  ctx.Status(401).Throw("UNAUTHORIZED")
})

api.Get("/:bar", func(ctx *rest.Context) {
  fmt.Println("authtoken", ctx.Get("authtoken"))
  ctx.JSON(ctx.Params)
})

// error handler
api.OnError("UNAUTHORIZED", func(ctx *rest.Context) {
  ctx.Text("You are unauthorized")
})

fmt.Println("Starting server.")

http.ListenAndServe(":8080", api)
```

## Documentation
https://godoc.org/github.com/go-rs/rest-api-framework

##### Powered by
[![GoLand - JetBrains](https://raw.githubusercontent.com/go-rs/rest-api-framework/master/docs/powered-by/logo.svg?sanitize=true)](https://www.jetbrains.com/?from=Go+REST+Services)
