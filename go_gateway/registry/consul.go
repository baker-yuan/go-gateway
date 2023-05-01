package registry

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

func NewDiscovery() registry.Discovery {
	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}
	return consul.New(consulClient)
}
