package plugin_manager

import (
	pb_router "github.com/baker-yuan/go-gateway/pb/router"
	gcontext "github.com/baker-yuan/go-gateway/pkg/context"
)

// PluginObj 对一个路由对应对插件集合进行封装
type PluginObj struct {
	fs   gcontext.Filters             // 拦截器
	conf map[string]*pb_router.Plugin // 插件配置
}

func NewPluginObj(filters gcontext.Filters, conf map[string]*pb_router.Plugin) *PluginObj {
	obj := &PluginObj{
		fs:   filters,
		conf: conf,
	}
	return obj
}

func (p *PluginObj) Chain(ctx gcontext.GatewayContext, append ...gcontext.IFilter) error {
	return gcontext.DoChain(ctx, p.fs, append...)
}

func (p *PluginObj) Destroy() {
	handler := p.fs
	if handler != nil {
		handler.Destroy()
	}
}
