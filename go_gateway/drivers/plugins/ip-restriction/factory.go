package ip_restriction

import (
	"github.com/baker-yuan/go-gateway/drivers"
	"github.com/baker-yuan/go-gateway/eosc"
)

const (
	Name = "ip_restriction" // IP黑白名单
)

func Register(register eosc.IExtenderDriverRegister) {
	_ = register.RegisterExtenderDriver(Name, NewFactory())
}

func NewFactory() eosc.IExtenderDriverFactory {
	// 创建的是Factory对象
	var factory *drivers.Factory[Config] = drivers.NewFactory[Config](Create, Check)
	// 返回
	return factory
}
