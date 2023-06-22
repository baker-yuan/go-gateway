package ip_restriction

import (
	"encoding/json"

	eocontext "github.com/baker-yuan/go-gateway/context"
	http_service "github.com/baker-yuan/go-gateway/context/http-context"
	"github.com/baker-yuan/go-gateway/drivers"
	"github.com/baker-yuan/go-gateway/eosc"
)

var _ http_service.HttpFilter = (*IPHandler)(nil)
var _ eocontext.IFilter = (*IPHandler)(nil)

// IPHandler ip黑白名单处理器
type IPHandler struct {
	drivers.WorkerBase          //
	responseType       string   // 响应类型
	filter             IPFilter // 业务逻辑
}

// DoFilter 过滤器逻辑
func (I *IPHandler) DoFilter(ctx eocontext.EoContext, next eocontext.IChain) (err error) {
	return http_service.DoHttpFilter(I, ctx, next)
}

func (I *IPHandler) doRestriction(ctx http_service.IHttpContext) error {
	// 获取客户端真实IP
	realIP := ctx.Request().RealIp()
	if I.filter != nil {
		ok, err := I.filter(realIP)
		if !ok {
			return err
		}
	}
	return nil
}

// Start 生命周期-启动
func (I *IPHandler) Start() error {
	return nil
}

// Reset 生命周期-重置
func (I *IPHandler) Reset(conf interface{}, workers map[eosc.RequireId]eosc.IWorker) error {
	confObj, err := check(conf)
	if err != nil {
		return err
	}
	I.filter = confObj.genFilter()
	return nil
}

// Stop 生命周期-停止
func (I *IPHandler) Stop() error {
	return nil
}

// CheckSkill 生命周期-CheckSkill
func (I *IPHandler) CheckSkill(skill string) bool {
	return http_service.FilterSkillName == skill
}

// responseEncode 拒绝执行返回错误信息
func (I *IPHandler) responseEncode(origin string, statusCode int) string {
	if I.responseType == "json" {
		tmp := map[string]interface{}{
			"message":     origin,
			"status_code": statusCode,
		}
		newInfo, _ := json.Marshal(tmp)
		return string(newInfo)
	}
	return origin
}

// Destroy 生命周期-销毁
func (I *IPHandler) Destroy() {
	I.filter = nil
	I.responseType = ""
}

// DoHttpFilter 过滤器逻辑
func (I *IPHandler) DoHttpFilter(ctx http_service.IHttpContext, next eocontext.IChain) error {
	err := I.doRestriction(ctx)
	// 拒绝执行
	if err != nil {
		resp := ctx.Response()
		info := I.responseEncode(err.Error(), 403)
		resp.SetStatus(403, "403")
		resp.SetBody([]byte(info))
		return err
	}
	// 继续执行下一个插件
	if next != nil {
		return next.DoChain(ctx)
	}
	return nil
}
