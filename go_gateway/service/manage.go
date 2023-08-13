package service_manager

import (
	"time"

	pb_service "github.com/baker-yuan/go-gateway/pb/service"
	gocontext "github.com/baker-yuan/go-gateway/pkg/context"
	pkg_service "github.com/baker-yuan/go-gateway/pkg/service"
)

// IServiceManager 服务管理器统一接口
type IServiceManager interface {
	pkg_service.IServiceGet
	Set(service *pb_service.Service)
}

type ServiceManager struct {
	origin   map[uint32]pb_service.Service   // key=服务ID value=服务
	services map[uint32]pkg_service.IService // key=服务ID value=服务
}

var (
	_ gocontext.BalanceHandler = (*Service)(nil)
	_ gocontext.IService       = (*Service)(nil)
)

type Service struct {
}

func (s Service) Instances() []gocontext.IInstance {
	// TODO implement me
	panic("implement me")
}

func (s Service) Select(ctx gocontext.GatewayContext) (gocontext.IInstance, int, error) {
	// TODO implement me
	panic("implement me")
}

func (s Service) Scheme() string {
	// TODO implement me
	panic("implement me")
}

func (s Service) TimeOut() time.Duration {
	// TODO implement me
	panic("implement me")
}

func (s *ServiceManager) GetService(id uint32) pkg_service.IService {
	return &Service{}
}

func (s *ServiceManager) Set(service *pb_service.Service) {

}

func NewServiceManager() IServiceManager {
	return &ServiceManager{}
}
