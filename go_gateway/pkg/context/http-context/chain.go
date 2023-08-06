package http_context

import (
	gcontext "github.com/baker-yuan/go-gateway/pkg/context"
)

// Assert EoContext是否是IHttpContext
func Assert(ctx gcontext.GatewayContext) (IHttpContext, error) {
	var httpContext IHttpContext
	err := ctx.Assert(&httpContext)
	return httpContext, err
}
