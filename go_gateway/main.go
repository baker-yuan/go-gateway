package main

import (
	"github.com/baker-yuan/go-gateway/http_rule"
	"github.com/baker-yuan/go-gateway/middleware/grpc_forward"
	"github.com/baker-yuan/go-gateway/middleware/pre_handle"
	"github.com/valyala/fasthttp"
	"github.com/vincentLiuxiang/lu"
)

func main() {

	if err := http_rule.LoadHttpRule(); err != nil {
		panic(err)
	}

	app := lu.New()

	app.Use("/", pre_handle.HTTPAccessModeMiddleware)
	app.Use("/", grpc_forward.GrpcForwardModeMiddleware)

	app.Get("/hello", func(ctx *fasthttp.RequestCtx, next func(error)) {
		ctx.SetBody([]byte("ok..."))
	})

	server := &fasthttp.Server{
		Handler:     app.Handler,
		Concurrency: 1024 * 1024,
	}
	server.ListenAndServe(":8080")
}
