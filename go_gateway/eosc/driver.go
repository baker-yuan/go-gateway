package eosc

import (
	"reflect"
)

// IExtenderDriverFactory 插件工厂
// 1、实现类 Factory
type IExtenderDriverFactory interface {
	Render() interface{}
	// Create 创建IExtenderDriver
	Create(profession string, name string, label string, desc string, params map[string]interface{}) (IExtenderDriver, error)
}

// IExtenderConfigChecker 插件配置检查
type IExtenderConfigChecker interface {
	// Check 插件配置检查
	Check(v interface{}, workers map[RequireId]IWorker) error
}

// IExtenderDriver 插件
// 1、Driver 插件
// 2、DriverConfigChecker 插件配置检查
type IExtenderDriver interface {
	// ConfigType 插件配置类型
	ConfigType() reflect.Type
	// Create 创建具体的插件
	Create(id, name string, v interface{}, workers map[RequireId]IWorker) (IWorker, error)
}

type SettingMode int

type ISetting interface {
	ConfigType() reflect.Type
	Set(conf interface{}) (err error)
	Get() interface{}
	Mode() SettingMode
	Check(cfg interface{}) (profession, name, driver, desc string, err error)
	AllWorkers() []string
}

type ISettings interface {
	GetDriver(name string) (ISetting, bool)
	SettingWorker(name string, config []byte, variable IVariable) error
	Update(name string, variable IVariable) (err error)
	CheckVariable(name string, variable IVariable) (err error)
	GetConfig(name string) interface{}
}
