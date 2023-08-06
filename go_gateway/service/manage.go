package service

import (
	gcontext "github.com/baker-yuan/go-gateway/pkg/context"
)

type IService interface {
	gcontext.IService       // 后端服务
	gcontext.BalanceHandler // 负载均衡
}

// ServiceManager 服务管理器统一接口
type ServiceManager interface {
}

type ServiceManagerImpl struct {
}

func NewServiceManager() ServiceManager {
	return &ServiceManagerImpl{}
}
