package pre_handle

import (
	"net/http"

	"github.com/baker-yuan/go-gateway/common"
	"github.com/baker-yuan/go-gateway/http_rule"
	"github.com/valyala/fasthttp"
)

func HTTPAccessModeMiddleware(ctx *fasthttp.RequestCtx, next func(error)) {
	// 获取配置
	rule, _ := http_rule.GetRule(string(ctx.Method()), string(ctx.Path()))
	if rule == nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}
	// 设置到上线文中
	ctx.SetUserValue(common.CtxRuleKey, rule)
	// 转发到下一个
	next(nil)
}
