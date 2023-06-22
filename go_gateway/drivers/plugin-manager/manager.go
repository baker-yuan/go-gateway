package plugin_manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/baker-yuan/go-gateway/common/bean"
	eocontext "github.com/baker-yuan/go-gateway/context"
	"github.com/baker-yuan/go-gateway/eosc"
	"github.com/baker-yuan/go-gateway/log"
	"github.com/baker-yuan/go-gateway/plugin"
)

var (
	errConfig                      = errors.New("invalid config")
	ErrorDriverNotExit             = errors.New("drive not exit")
	ErrorGlobalPluginMastConfig    = errors.New("global must have config")
	ErrorGlobalPluginConfigInvalid = errors.New("invalid global config")
)

// PluginManager 插件管理器
// 1、eosc/driver.go#ISetting
// 2、apinto/plugin/plugin.go#IPluginManager
type PluginManager struct {
	name            string                           // 名称
	extenderDrivers eosc.IExtenderDrivers            // 获取插件工厂（插件注册）
	workers         eosc.IWorkers                    // 插件获取
	plugins         Plugins                          // 插件集合
	pluginObjs      eosc.Untyped[string, *PluginObj] // 拦击器（成品）key=插件id value=PluginObj
	global          eocontext.IChainPro              // 全局插件
}

func NewPluginManager() *PluginManager {
	pm := &PluginManager{
		name:       "plugin",
		plugins:    make(Plugins, 0),
		pluginObjs: eosc.BuildUntyped[string, *PluginObj](),
	}
	log.Debug("autowired extenderDrivers")

	// 注入获取插件工厂
	// eosc/process-worker/process.go#NewProcessWorker
	bean.Autowired(&pm.extenderDrivers)
	fmt.Println("[info] 获取 IExtenderDrivers go_gateway/drivers/plugin-manager/manager.go newPluginManager...")

	// 注入插件
	// eosc/process-worker/service.go#NewWorkerServer
	bean.Autowired(&pm.workers) // 获取插件
	fmt.Println("[info] 获取 IWorkers go_gateway/drivers/plugin-manager/manager.go newPluginManager...")

	log.DebugF("autowired extenderDrivers = %p", pm.extenderDrivers)
	return pm
}

func (p *PluginManager) Global() eocontext.IChainPro {
	if p.global == nil {
		p.global = p.createChain("global", map[string]*plugin.Config{})
	}
	return p.global
}

func (p *PluginManager) Check(cfg interface{}) (profession, name, driver, desc string, err error) {
	err = eosc.ErrorUnsupportedKind
	return
}

func (p *PluginManager) AllWorkers() []string {
	return []string{"plugin@setting"}
}

func (p *PluginManager) Mode() eosc.SettingMode {
	return eosc.SettingModeSingleton
}

func (p *PluginManager) Set(conf interface{}) (err error) {
	err = p.Reset(conf)
	return
}

func (p *PluginManager) Get() interface{} {
	return p.plugins
}

func (p *PluginManager) ConfigType() reflect.Type {
	return reflect.TypeOf(new(PluginWorkerConfig))
}

func (p *PluginManager) CreateRequest(id string, conf map[string]*plugin.Config) eocontext.IChainPro {
	return p.createChain(id, conf)
}

func (p *PluginManager) GetConfigType(name string) (reflect.Type, bool) {
	log.Debug("plugin manager get config type:", p.plugins)
	for _, plg := range p.plugins {
		if name == plg.Name {
			return plg.drive.ConfigType(), true
		}
	}
	return nil, false
}

func (p *PluginManager) Reset(conf interface{}) error {
	plugins, err := p.check(conf)
	if err != nil {
		return err
	}
	p.plugins = plugins
	list := p.pluginObjs.List()

	// 遍历，全量更新
	for _, v := range list {
		v.fs = p.createFilters(v.conf)
	}

	return nil
}

func (p *PluginManager) createFilters(conf map[string]*plugin.Config) []eocontext.IFilter {
	filters := make([]eocontext.IFilter, 0, len(conf))
	plugins := p.plugins

	for _, plg := range plugins {

		if plg.Status == StatusDisable {
			// 禁用插件，跳过
			continue
		}

		c := plg.Config
		// 如果传入了配置，就用传入的配置
		if v, ok := conf[plg.Name]; ok {
			// 局部禁用
			if v.Disable {
				continue
			}
			// 替换为传入的配置
			if v.Config != nil {
				c = v.Config
			}
		} else if plg.Status != StatusGlobal {
			// 非全局插件跳过
			continue
		}

		// 配置反序列化
		confObj, err := toConfig(c, plg.drive.ConfigType())
		if err != nil {
			log.Error("plg manager: fail to createFilters filter,error is ", err)
			continue
		}
		// 创建插件
		worker, err := plg.drive.Create(fmt.Sprintf("%s@%s", plg.Name, p.name), plg.Name, confObj, nil)
		if err != nil {
			log.Error("plg manager: fail to createFilters filter,error is ", err)
			continue
		}
		// 类型转换
		fi, ok := worker.(eocontext.IFilter)
		if !ok {
			log.Error("extender ", plg.ID, " not plg for http-service.Filter")
			continue
		}
		// 加入集合
		filters = append(filters, fi)
	}
	return filters
}

func (p *PluginManager) createChain(id string, conf map[string]*plugin.Config) *PluginObj {
	chain := p.createFilters(conf)
	obj, has := p.pluginObjs.Get(id)
	if !has {
		obj = NewPluginObj(chain, id, conf)
		p.pluginObjs.Set(id, obj)
	} else {
		obj.fs = chain
	}
	log.Debug("create chain len: ", len(chain))
	return obj
}

func (p *PluginManager) check(conf interface{}) (Plugins, error) {
	cfg, ok := conf.(*PluginWorkerConfig)
	if !ok {
		return nil, errConfig
	}
	plugins := make(Plugins, 0, len(cfg.Plugins))
	for i, cf := range cfg.Plugins {
		log.DebugF("new plugin:%d=>%v", i, cf)
		newPlugin, err := p.newPlugin(cf)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, newPlugin)
	}
	return plugins, nil
}

func (p *PluginManager) IsExists(id string) bool {
	_, has := p.extenderDrivers.GetDriver(id)
	return has
}

func toConfig(v interface{}, t reflect.Type) (interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	obj := newConfig(t)
	err = json.Unmarshal(data, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func newConfig(t reflect.Type) interface{} {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}
