package http_complete

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/baker-yuan/go-gateway/pkg/context"
	http_service "github.com/baker-yuan/go-gateway/pkg/context/http-context"
	"github.com/baker-yuan/go-gateway/pkg/log"
	"github.com/baker-yuan/go-gateway/pkg/model/ctx_key"
	"github.com/baker-yuan/go-gateway/pkg/model/router"
)

var (
	ErrorTimeoutComplete = errors.New("complete timeout")
)

// HttpComplete 调用下游服务
type HttpComplete struct {
}

func NewHttpComplete() *HttpComplete {
	return &HttpComplete{}
}

// Complete 发送http请求到下游服务
func (h *HttpComplete) Complete(org context.GatewayContext) error {
	ctx, err := http_service.Assert(org)
	if err != nil {
		return err
	}
	// 设置响应开始时间
	proxyTime := time.Now()

	defer func() {
		// 设置原始响应状态码
		ctx.Response().SetProxyStatus(ctx.Response().StatusCode(), "")
		// 设置上游响应总时间, 单位为毫秒
		ctx.Response().SetResponseTime(time.Since(proxyTime))
		ctx.SetLabel("handler", "proxy")
	}()

	balance := ctx.GetBalance()

	scheme := balance.Scheme()
	switch strings.ToLower(scheme) {
	case "", "tcp":
		scheme = "http"
	case "tsl", "ssl", "https":
		scheme = "https"

	}
	timeOut := balance.TimeOut()

	// 重试
	retryValue := ctx.Value(ctx_key.CtxKeyRetry)
	retry, ok := retryValue.(int)
	if !ok {
		retry = router.DefaultRetry
	}
	// 请求超时
	timeoutValue := ctx.Value(ctx_key.CtxKeyTimeout)
	timeout, ok := timeoutValue.(time.Duration)
	if !ok {
		timeout = router.DefaultTimeout
	}

	var lastErr error
	for index := 0; index <= retry; index++ {
		if timeout > 0 && time.Since(proxyTime) > timeout {
			return ErrorTimeoutComplete
		}
		// 负载均衡
		node, _, err := balance.Select(ctx)
		if err != nil {
			log.Error("select error: ", err)
			ctx.Response().SetStatus(501, "501")
			ctx.Response().SetBody([]byte(err.Error()))
			return err
		}
		// 发送请求到下游服务
		lastErr = ctx.SendTo(scheme, node, timeOut)
		if lastErr == nil {
			return nil
		}
		log.Error("http upstream send error: ", lastErr)
	}

	return lastErr
}

// NoServiceCompleteHandler 查找下游服务失败
type NoServiceCompleteHandler struct {
	status int
	header map[string]string
	body   string
}

func NewNoServiceCompleteHandler(status int, header map[string]string, body string) *NoServiceCompleteHandler {
	return &NoServiceCompleteHandler{status: status, header: header, body: body}
}

func (n *NoServiceCompleteHandler) Complete(org context.GatewayContext) error {
	ctx, err := http_service.Assert(org)
	if err != nil {
		return err
	}
	// 设置响应开始时间
	proxyTime := time.Now()

	defer func() {
		// 设置原始响应状态码
		ctx.Response().SetProxyStatus(ctx.Response().StatusCode(), "")
		// 设置上游响应总时间, 单位为毫秒
		ctx.Response().SetResponseTime(time.Since(proxyTime))
		ctx.SetLabel("handler", "proxy")
	}()
	for key, value := range n.header {
		ctx.Response().SetHeader(key, value)
	}
	ctx.Response().SetBody([]byte(n.body))
	ctx.Response().SetStatus(n.status, strconv.Itoa(n.status))
	return nil
}

// HttpCompleteCaller 把调用下游请求的操作也封装为一个拦截器
type HttpCompleteCaller struct {
}

func NewHttpCompleteCaller() *HttpCompleteCaller {
	return &HttpCompleteCaller{}
}

func (h *HttpCompleteCaller) DoFilter(ctx context.GatewayContext, next context.IChain) (err error) {
	return ctx.GetComplete().Complete(ctx)
}

func (h *HttpCompleteCaller) Destroy() {
}
