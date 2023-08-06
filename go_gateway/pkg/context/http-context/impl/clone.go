package http_context

import (
	"context"
	"fmt"
	"net"
	"time"

	gcontext "github.com/baker-yuan/go-gateway/pkg/context"
	http_service "github.com/baker-yuan/go-gateway/pkg/context/http-context"
	fasthttp_client "github.com/baker-yuan/go-gateway/pkg/fasthttp-client"
	"github.com/baker-yuan/go-gateway/pkg/util"
	"github.com/valyala/fasthttp"
)

var _ http_service.IHttpContext = (*cloneContext)(nil)

// HttpContext fasthttpRequestCtx
type cloneContext struct {
	org             *HttpContext             //
	proxyRequest    ProxyRequest             //
	response        Response                 //
	proxyRequests   []http_service.IProxy    //
	ctx             context.Context          //
	completeHandler gcontext.CompleteHandler //
	finishHandler   gcontext.FinishHandler   //
	app             gcontext.IService        //
	balance         gcontext.BalanceHandler  //
	labels          map[string]string        //
	responseError   error                    //
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
func (ctx *cloneContext) SendTo(scheme string, node gcontext.IInstance, timeout time.Duration) error {
	host := node.Addr()
	request := ctx.proxyRequest.Request()

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
	return fmt.Errorf("not suport:%s", util.TypeNameOf(i))
}

func (ctx *cloneContext) GetComplete() gcontext.CompleteHandler {
	return ctx.completeHandler
}

func (ctx *cloneContext) SetCompleteHandler(handler gcontext.CompleteHandler) {
	ctx.completeHandler = handler
}

func (ctx *cloneContext) GetFinish() gcontext.FinishHandler {
	return ctx.finishHandler
}

func (ctx *cloneContext) SetFinish(handler gcontext.FinishHandler) {
	ctx.finishHandler = handler
}

func (ctx *cloneContext) GetBalance() gcontext.BalanceHandler {
	return ctx.balance
}

func (ctx *cloneContext) SetBalance(handler gcontext.BalanceHandler) {
	ctx.balance = handler
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

func (ctx *cloneContext) Clone() (gcontext.GatewayContext, error) {
	return nil, fmt.Errorf("%s %w", "HttpContext", gcontext.ErrEoCtxUnCloneable)
}
