package http_context

import (
	gcontext "github.com/baker-yuan/go-gateway/context"
	"github.com/baker-yuan/go-gateway/utils/config"
)

var (
	FilterSkillName = config.TypeNameOf((*HttpFilter)(nil))
)

type HttpFilter interface {
	DoHttpFilter(ctx IHttpContext, next gcontext.IChain) (err error)
}

func DoHttpFilter(httpFilter HttpFilter, ctx gcontext.EoContext, next gcontext.IChain) (err error) {
	httpContext, err := Assert(ctx)
	if err == nil {
		return httpFilter.DoHttpFilter(httpContext, next)
	}
	if next != nil {
		return next.DoChain(ctx)
	}
	return err
}

type WebsocketFilter interface {
	DoWebsocketFilter(ctx IWebsocketContext, next gcontext.IChain) error
}
