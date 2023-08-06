package http_context

import (
	gcontext "github.com/baker-yuan/go-gateway/pkg/context"
	"github.com/baker-yuan/go-gateway/pkg/util"
)

var (
	FilterSkillName = util.TypeNameOf((*HttpFilter)(nil))
)

type HttpFilter interface {
	DoHttpFilter(ctx IHttpContext, next gcontext.IChain) (err error)
}

func DoHttpFilter(httpFilter HttpFilter, ctx gcontext.GatewayContext, next gcontext.IChain) (err error) {
	httpContext, err := Assert(ctx)
	if err == nil {
		return httpFilter.DoHttpFilter(httpContext, next)
	}
	if next != nil {
		return next.DoChain(ctx)
	}
	return err
}
