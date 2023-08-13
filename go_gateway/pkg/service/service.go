package service

import gcontext "github.com/baker-yuan/go-gateway/pkg/context"

type IService interface {
	gcontext.IService       // 后端服务
	gcontext.BalanceHandler // 负载均衡
}

// IServiceGet 获取服务
type IServiceGet interface {
	// GetService 获取服务
	GetService(id uint32) IService
}
