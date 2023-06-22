package plugin_manager

import (
	eoscContext "github.com/baker-yuan/go-gateway/context"
	"github.com/baker-yuan/go-gateway/plugin"
)

type PluginObj struct {
	id   string                    // 插件id
	fs   eoscContext.Filters       // 拦截器
	conf map[string]*plugin.Config // 插件配置
}

func NewPluginObj(filters eoscContext.Filters, id string, conf map[string]*plugin.Config) *PluginObj {
	obj := &PluginObj{
		fs:   filters,
		id:   id,
		conf: conf,
	}
	return obj
}

func (p *PluginObj) Chain(ctx eoscContext.EoContext, append ...eoscContext.IFilter) error {
	return eoscContext.DoChain(ctx, p.fs, append...)
}
func (p *PluginObj) Destroy() {
	handler := p.fs
	if handler != nil {
		handler.Destroy()
	}
}
