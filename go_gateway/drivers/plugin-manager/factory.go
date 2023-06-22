package plugin_manager

import (
	"fmt"

	"github.com/baker-yuan/go-gateway/common/bean"
	"github.com/baker-yuan/go-gateway/eosc"
	"github.com/baker-yuan/go-gateway/plugin"
	"github.com/baker-yuan/go-gateway/setting"
)

var (
	// 插件管理器，单例
	singleton *PluginManager
	_         eosc.ISetting = singleton
)

func Init() {
	fmt.Println("[info] 创建插件管理器 IPluginManager go_gateway/drivers/plugin-manager/factory.go init...")
	// 往容器中注入插件管理器
	singleton = NewPluginManager()
	var i plugin.IPluginManager = singleton
	bean.Injection(&i)
}

// Register 注册插件回调，这里注册了一个Setting
func Register(register eosc.IExtenderDriverRegister) {
	_ = setting.RegisterSetting("plugin", singleton)
}
