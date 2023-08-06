package http_router

import (
	gcontext "github.com/baker-yuan/go-gateway/pkg/context"
	http_context "github.com/baker-yuan/go-gateway/pkg/context/http-context"
)

var defaultFinisher = &Finisher{}

type Finisher struct {
}

func (f *Finisher) Finish(org gcontext.GatewayContext) error {
	ctx, err := http_context.Assert(org)
	if err != nil {
		return err
	}
	ctx.FastFinish()
	return nil
}
