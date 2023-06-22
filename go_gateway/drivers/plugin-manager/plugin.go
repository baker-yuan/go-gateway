package plugin_manager

import (
	"fmt"

	"github.com/baker-yuan/go-gateway/eosc"
	"github.com/baker-yuan/go-gateway/log"
)

type Plugins []*Plugin

type Plugin struct {
	*PluginConfig
	drive eosc.IExtenderDriver
}

func (p *PluginManager) newPlugin(conf *PluginConfig) (*Plugin, error) {
	// 获取IExtenderDriver
	driver, err := p.getExtenderDriver(conf)
	if err != nil {
		return nil, err
	}

	// 全局插件&配置为空
	if conf.Status == StatusGlobal && conf.Config == nil {
		return nil, ErrorGlobalPluginMastConfig
	}

	// 全局插件
	if conf.Status == StatusGlobal {
		// 获取配置
		v, err := toConfig(conf.Config, driver.ConfigType())
		if err != nil {
			log.Info("global plugin:", conf.Name, "config:", err)
			return nil, fmt.Errorf("%s:%w", conf.Name, ErrorGlobalPluginConfigInvalid)
		}
		// 检查配置
		if dc, ok := driver.(eosc.IExtenderConfigChecker); ok {
			errCheck := dc.Check(v, nil)
			if errCheck != nil {
				return nil, errCheck
			}
		}
	}

	// 创建Plugin
	return &Plugin{
		PluginConfig: conf,
		drive:        driver,
	}, nil

}

func (p *PluginManager) getExtenderDriver(config *PluginConfig) (eosc.IExtenderDriver, error) {
	// 获取IExtenderDriverFactory
	driverFactory, has := p.extenderDrivers.GetDriver(config.ID)
	if !has {
		return nil, fmt.Errorf("id:%w", ErrorDriverNotExit)
	}
	// 获取IExtenderDriver
	return driverFactory.Create("plugin@setting", config.Name, config.Name, config.Name, config.InitConfig)
}
