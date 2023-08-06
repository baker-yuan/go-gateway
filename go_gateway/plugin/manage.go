package plugin

import (
	pgk_plugin "github.com/baker-yuan/go-gateway/pkg/plugin"
	"github.com/baker-yuan/go-gateway/plugin/demo"
)

// PluginManager 插件管理器统一接口
type PluginManager interface {
	RegisterSchema(plg *pgk_plugin.PluginSchema)
}

type PluginManagerImpl struct {
	models map[string]*pgk_plugin.PluginSchema // 插件定义 PluginModel#name -> PluginModel
}

func NewPluginManager() PluginManager {
	pluginManager := &PluginManagerImpl{
		models: make(map[string]*pgk_plugin.PluginSchema),
	}

	// 注册插件
	pluginManager.RegisterSchema(demo.PluginSchema())

	return pluginManager
}

// RegisterSchema 注册插件定义信息
func (pm *PluginManagerImpl) RegisterSchema(plg *pgk_plugin.PluginSchema) {
	pm.models[plg.Name] = plg
}
