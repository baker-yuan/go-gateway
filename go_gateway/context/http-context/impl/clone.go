package http_context

import (
	"context"
	"fmt"
	"net"
	"time"

	eoscContext "github.com/baker-yuan/go-gateway/context"
	http_service "github.com/baker-yuan/go-gateway/context/http-context"
	fasthttp_client "github.com/baker-yuan/go-gateway/node/fasthttp-client"
	"github.com/baker-yuan/go-gateway/utils/config"
	"github.com/valyala/fasthttp"
)

var _ http_service.IHttpContext = (*cloneContext)(nil)

// HttpContext fasthttpRequestCtx
type cloneContext struct {
	org                 *HttpContext                    //
	proxyRequest        ProxyRequest                    //
	response            Response                        //
	proxyRequests       []http_service.IProxy           //
	ctx                 context.Context                 //
	completeHandler     eoscContext.CompleteHandler     //
	finishHandler       eoscContext.FinishHandler       //
	app                 eoscContext.EoApp               //
	balance             eoscContext.BalanceHandler      //
	upstreamHostHandler eoscContext.UpstreamHostHandler //
	labels              map[string]string               //
	responseError       error                           //
}

// -------------------- go_gateway/context/http-context/context.go#IHttpContext实现 --------------------

func (ctx *cloneContext) Request() http_service.IRequestReader {
	return ctx.org.Request()
}

func (ctx *cloneContext) Proxy() http_service.IRequest {
	return &ctx.proxyRequest
}

func (ctx *cloneContext) Response() http_service.IResponse {
	return &ctx.response
}

// SendTo 转发请求到下游到http服务
func (ctx *cloneContext) SendTo(scheme string, node eoscContext.INode, timeout time.Duration) error {
	host := node.Addr()
	request := ctx.proxyRequest.Request()

	passHost, targetHost := ctx.GetUpstreamHostHandler().PassHost()
	switch passHost {
	case eoscContext.PassHost:
	case eoscContext.NodeHost:
		request.URI().SetHost(node.Addr())
	case eoscContext.ReWriteHost:
		request.URI().SetHost(targetHost)
	}

	beginTime := time.Now()
	ctx.responseError = fasthttp_client.ProxyTimeout(scheme, node, request, ctx.response.Response, timeout)
	agent := newRequestAgent(&ctx.proxyRequest, host, scheme, beginTime, time.Now())
	if ctx.responseError != nil {
		agent.setStatusCode(504)
	} else {
		agent.setStatusCode(ctx.response.Response.StatusCode())
	}

	agent.setResponseLength(ctx.response.Response.Header.ContentLength())

	ctx.proxyRequests = append(ctx.proxyRequests, agent)
	return ctx.responseError

}

func (ctx *cloneContext) Proxies() []http_service.IProxy {
	return ctx.proxyRequests
}

func (ctx *cloneContext) FastFinish() {
	ctx.ctx = nil
	ctx.app = nil
	ctx.balance = nil
	ctx.upstreamHostHandler = nil
	ctx.finishHandler = nil
	ctx.completeHandler = nil
	fasthttp.ReleaseRequest(ctx.proxyRequest.req)
	fasthttp.ReleaseResponse(ctx.response.Response)
	ctx.response.Finish()
	ctx.proxyRequest.Finish()
}

// -------------------- go_gateway/context/context.go#EoContext 实现 --------------------

// RequestId 请求ID
func (ctx *cloneContext) RequestId() string {
	return ctx.org.requestID
}

func (ctx *cloneContext) AcceptTime() time.Time {
	return ctx.org.AcceptTime()
}

func (ctx *cloneContext) Context() context.Context {
	return ctx.ctx
}

func (ctx *cloneContext) Value(key interface{}) interface{} {
	return ctx.ctx.Value(key)
}

func (ctx *cloneContext) WithValue(key, val interface{}) {
	ctx.ctx = context.WithValue(ctx.Context(), key, val)
}

func (ctx *cloneContext) Scheme() string {
	return ctx.org.Scheme()
}

func (ctx *cloneContext) Assert(i interface{}) error {
	if v, ok := i.(*http_service.IHttpContext); ok {
		*v = ctx
		return nil
	}
	return fmt.Errorf("not suport:%s", config.TypeNameOf(i))
}

func (ctx *cloneContext) GetComplete() eoscContext.CompleteHandler {
	return ctx.completeHandler
}

func (ctx *cloneContext) SetCompleteHandler(handler eoscContext.CompleteHandler) {
	ctx.completeHandler = handler
}

func (ctx *cloneContext) GetFinish() eoscContext.FinishHandler {
	return ctx.finishHandler
}

func (ctx *cloneContext) SetFinish(handler eoscContext.FinishHandler) {
	ctx.finishHandler = handler
}

func (ctx *cloneContext) GetBalance() eoscContext.BalanceHandler {
	return ctx.balance
}

func (ctx *cloneContext) SetBalance(handler eoscContext.BalanceHandler) {
	ctx.balance = handler
}

func (ctx *cloneContext) GetUpstreamHostHandler() eoscContext.UpstreamHostHandler {
	return ctx.upstreamHostHandler
}

func (ctx *cloneContext) SetUpstreamHostHandler(handler eoscContext.UpstreamHostHandler) {
	ctx.upstreamHostHandler = handler
}

func (ctx *cloneContext) SetLabel(name, value string) {
	ctx.labels[name] = value
}

func (ctx *cloneContext) GetLabel(name string) string {
	return ctx.labels[name]
}

func (ctx *cloneContext) Labels() map[string]string {
	return ctx.labels
}

func (ctx *cloneContext) RealIP() string {
	return ctx.org.RealIP()
}

func (ctx *cloneContext) LocalIP() net.IP {
	return ctx.org.LocalIP()
}

func (ctx *cloneContext) LocalAddr() net.Addr {
	return ctx.org.LocalAddr()
}

func (ctx *cloneContext) LocalPort() int {
	return ctx.org.LocalPort()
}

func (ctx *cloneContext) IsCloneable() bool {
	return false
}

func (ctx *cloneContext) Clone() (eoscContext.EoContext, error) {
	return nil, fmt.Errorf("%s %w", "HttpContext", eoscContext.ErrEoCtxUnCloneable)
}