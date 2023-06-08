package http_context

import gcontext "github.com/baker-yuan/go-gateway/context"

// Assert EoContext是否是IHttpContext
func Assert(ctx gcontext.EoContext) (IHttpContext, error) {
	var httpContext IHttpContext
	err := ctx.Assert(&httpContext)
	return httpContext, err
}

// WebsocketAssert EoContext是否是IWebsocketContext
func WebsocketAssert(ctx gcontext.EoContext) (IWebsocketContext, error) {
	var websocketContext IWebsocketContext
	err := ctx.Assert(&websocketContext)
	return websocketContext, err
}
