package server

import (
	"context"

	v1 "github.com/baker-yuan/go-gateway/go-gateway-admin/api/admin/v1"
	"github.com/baker-yuan/go-gateway/go-gateway-admin/internal/conf"
	"github.com/baker-yuan/go-gateway/go-gateway-admin/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// MiddlewareCors 设置跨域请求头
func MiddlewareCors() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if ts, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := ts.(http.Transporter); ok {
					ht.ReplyHeader().Set("Access-Control-Allow-Origin", "*")
					ht.ReplyHeader().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,PATCH,DELETE")
					ht.ReplyHeader().Set("Access-Control-Allow-Credentials", "true")
					ht.ReplyHeader().Set("Access-Control-Allow-Headers", "Content-Type,"+
						"X-Requested-With,Access-Control-Allow-Credentials,User-Agent,Content-Length,Authorization")
				}
			}
			return handler(ctx, req)
		}
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.HttpRuleService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			MiddlewareCors(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGatewayAdminHTTPServer(srv, greeter)
	return srv
}
