package static

import (
	"github.com/baker-yuan/go-gateway/drivers"
	"github.com/baker-yuan/go-gateway/eosc"
)

var name = "discovery_static"

// Register 注册静态服务发现的驱动工厂
func Register(register eosc.IExtenderDriverRegister) {
	register.RegisterExtenderDriver(name, NewFactory())
}
func NewFactory() eosc.IExtenderDriverFactory {
	return drivers.NewFactory[Config](Create)
}
