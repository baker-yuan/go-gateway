package plugin_manager

import (
	pb_router "github.com/baker-yuan/go-gateway/pb/router"
	gcontext "github.com/baker-yuan/go-gateway/pkg/context"
	pkg_plugin "github.com/baker-yuan/go-gateway/pkg/plugin"
	"github.com/baker-yuan/go-gateway/plugin/demo"
)

// IPluginManager 插件管理器统一接口
type IPluginManager interface {
	RegisterSchema(plg *pkg_plugin.PluginSchema)                        // 注册插件定义
	CreateRequest(conf map[string]*pb_router.Plugin) gcontext.IChainPro // 获取插件
}

type PluginManager struct {
	schemas map[string]*pkg_plugin.PluginSchema // 插件定义 PluginModel#name -> PluginModel
}

func (m *PluginManager) CreateRequest(conf map[string]*pb_router.Plugin) gcontext.IChainPro {
	filters := make([]gcontext.IFilter, 0, len(conf))
	// for schema := range m.schemas {
	//
	// }
	filters = append(filters, demo.DemoPlugin{})

	return NewPluginObj(filters, conf)
}

func NewPluginManager() IPluginManager {
	pluginManager := &PluginManager{
		schemas: make(map[string]*pkg_plugin.PluginSchema),
	}

	// 注册插件
	pluginManager.RegisterSchema(demo.PluginSchema())

	return pluginManager
}

// RegisterSchema 注册插件定义信息
func (m *PluginManager) RegisterSchema(plg *pkg_plugin.PluginSchema) {
	m.schemas[plg.Name] = plg
}
