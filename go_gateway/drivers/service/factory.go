package service

import (
	"github.com/baker-yuan/go-gateway/drivers"
	"github.com/baker-yuan/go-gateway/drivers/discovery/static"
	"github.com/baker-yuan/go-gateway/eosc"
	"github.com/baker-yuan/go-gateway/log"
	iphash "github.com/baker-yuan/go-gateway/upstream/ip-hash"
	roundrobin "github.com/baker-yuan/go-gateway/upstream/round-robin"
)

var DriverName = "service_http"
var (
	defaultHttpDiscovery = static.CreateAnonymous(&static.Config{
		Health:   nil,
		HealthOn: false,
	})
)

// Register 注册service_http驱动工厂
func Register(register eosc.IExtenderDriverRegister) {
	err := register.RegisterExtenderDriver(DriverName, NewFactory())
	if err != nil {
		log.Errorf("register %s %s", DriverName, err)
		return

	}
}

// NewFactory 创建service_http驱动工厂
func NewFactory() eosc.IExtenderDriverFactory {
	roundrobin.Register()
	iphash.Register()
	return drivers.NewFactory[Config](Create)
}
