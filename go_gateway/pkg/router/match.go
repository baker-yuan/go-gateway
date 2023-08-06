package router

import (
	gcontext "github.com/baker-yuan/go-gateway/pkg/context"
)

// IMatcher 路由匹配
type IMatcher interface {
	// Match 路由匹配
	// @port 		端口
	// @request 入参
	// @return 处理请求的函数 是否匹配成功
	Match(port int, request interface{}) (IRouterHandler, bool)
}

// IRouterHandler 处理请求
type IRouterHandler interface {
	ServeHTTP(ctx gcontext.GatewayContext)
}
