package static

import (
	"github.com/baker-yuan/go-gateway/discovery"
	"github.com/baker-yuan/go-gateway/drivers"
	"github.com/baker-yuan/go-gateway/eosc"
)

// Create 创建静态服务发现驱动的实例
func Create(id, name string, cfg *Config, workers map[eosc.RequireId]eosc.IWorker) (eosc.IWorker, error) {

	s := &static{
		WorkerBase: drivers.Worker(id, name),
		services:   discovery.NewAppContainer(),
		cfg:        cfg,
	}
	return s, nil
}

func CreateAnonymous(conf *Config) discovery.IDiscovery {
	s := &static{
		cfg:      conf,
		services: discovery.NewAppContainer(),
	}
	s.Start()
	return s
}
