package service_manager

import (
	gcontext "github.com/baker-yuan/go-gateway/pkg/context"
)

type IService interface {
	gcontext.IService       // 后端服务
	gcontext.BalanceHandler // 负载均衡
}

// IServiceManager 服务管理器统一接口
type IServiceManager interface {
}

type ServiceManager struct {
}

func NewServiceManager() IServiceManager {
	return &ServiceManager{}
}
