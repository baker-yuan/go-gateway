package eosc

import (
	"fmt"
)

// IExtenderDriverRegister 注册插件
type IExtenderDriverRegister interface {
	// RegisterExtenderDriver 具体的插件调用这个方法进行注册
	RegisterExtenderDriver(name string, factory IExtenderDriverFactory) error
}

// IExtenderDrivers 获取插件
type IExtenderDrivers interface {
	GetDriver(name string) (IExtenderDriverFactory, bool)
}

// ExtenderRegister 插件注册中心
// 1、实现了IExtenderDriverRegister 用于注册插件
// 2、实现了IExtenderDrivers 用于获取插件
type ExtenderRegister struct {
	data IRegister[IExtenderDriverFactory] // key=插件名称 value=插件工厂
}

// RegisterExtenderDriver 注册插件
func (p *ExtenderRegister) RegisterExtenderDriver(name string, factory IExtenderDriverFactory) error {
	err := p.data.Register(name, factory, false)
	if err != nil {
		return fmt.Errorf("register profession  driver %s:%w", name, err)
	}
	return nil
}

// GetDriver 获取插件
func (p *ExtenderRegister) GetDriver(name string) (IExtenderDriverFactory, bool) {
	if v, has := p.data.Get(name); has {
		return v, true
	}
	return nil, false
}

// NewExtenderRegister 创建插件工厂注册中心实例
func NewExtenderRegister() *ExtenderRegister {
	return &ExtenderRegister{
		data: NewRegister[IExtenderDriverFactory](),
	}
}
