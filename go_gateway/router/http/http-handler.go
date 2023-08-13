package router_manager

import (
	"fmt"
	"time"

	"github.com/baker-yuan/go-gateway/pkg/context"
	http_context "github.com/baker-yuan/go-gateway/pkg/context/http-context"
	"github.com/baker-yuan/go-gateway/pkg/model/ctx_key"
	pkg_service "github.com/baker-yuan/go-gateway/pkg/service"
	http_complete "github.com/baker-yuan/go-gateway/router/http/http-complete"
)

// 发送http请求到下游服务
var completeCaller = http_complete.NewHttpCompleteCaller()

// httpHandler 处理http请求，实现接口IRouterHandler，一个路由对应一个httpHandler
type httpHandler struct {
	routerID        uint32                  // 路由ID
	serviceID       uint32                  // 服务ID
	disable         bool                    // 是否禁用路由
	retry           uint32                  // 超时重试次数
	timeout         time.Duration           // 超时时间，当为0时不设置超时，单位：ms
	service         pkg_service.IService    // 服务信息
	filters         context.IChainPro       // 拦击器链
	completeHandler context.CompleteHandler // 完成请求
	finisher        context.FinishHandler   // 资源清理
}

func (h *httpHandler) ServeHTTP(ctx context.GatewayContext) {
	httpContext, err := http_context.Assert(ctx)
	if err != nil {
		return
	}

	// 路由被禁用
	if h.disable {
		_ = httpContext.GetComplete().Complete(ctx)
		httpContext.FastFinish()
		return
	}

	// set retry timeout
	ctx.WithValue(ctx_key.CtxKeyRetry, h.retry)
	ctx.WithValue(ctx_key.CtxKeyTimeout, h.timeout)

	// Set Label
	ctx.SetLabel("api_id", fmt.Sprintf("%d", h.routerID))
	ctx.SetLabel("service_id", fmt.Sprintf("%d", h.serviceID))
	ctx.SetLabel("method", httpContext.Request().Method())
	ctx.SetLabel("path", httpContext.Request().URI().RequestURI())
	ctx.SetLabel("ip", httpContext.Request().RealIp())

	ctx.SetCompleteHandler(h.completeHandler)
	ctx.SetBalance(h.service)
	ctx.SetFinish(h.finisher)

	// 执行拦截器链
	_ = h.filters.Chain(ctx, completeCaller)
}
