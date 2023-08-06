package demo

import (
	gcontext "github.com/baker-yuan/go-gateway/pkg/context"
	http_context "github.com/baker-yuan/go-gateway/pkg/context/http-context"
)

type DemoPlugin struct {
}

func (d DemoPlugin) DoFilter(ctx gcontext.GatewayContext, next gcontext.IChain) error {
	return nil
}

func (d DemoPlugin) Destroy() {
}

func (d DemoPlugin) DoHttpFilter(ctx http_context.IHttpContext, next gcontext.IChain) error {
	return nil
}
