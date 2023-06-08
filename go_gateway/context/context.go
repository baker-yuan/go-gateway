package eocontext

import (
	"context"
	"errors"
	"net"
	"time"
)

var ErrEoCtxUnCloneable = errors.New("EoContext is UnCloneable. ")

// CompleteHandler 所有的插件执行完成后，需要执行的操作，1、转发到下游http服务 2、转发到下游grpc服务 3、返回路由失败异常
type CompleteHandler interface {
	Complete(ctx EoContext) error
}

// FinishHandler 请求执行完后，做清理的
type FinishHandler interface {
	Finish(ctx EoContext) error
}

// EoContext 上下文，不同的协议实现自己的上下文
type EoContext interface {
	RequestId() string                                  // 请求id唯一，每次请求随机生成
	AcceptTime() time.Time                              // 请求接收时间
	Context() context.Context                           // 原始context
	Value(key interface{}) interface{}                  // 从原始context中返回键对应的值
	WithValue(key, val interface{})                     // 往原始context添加键值对
	Scheme() string                                     // 协议 http、https、grpc、dubbo
	Assert(i interface{}) error                         // context类型断言
	SetLabel(name, value string)                        // 设置标签
	GetLabel(name string) string                        // 获取标签
	Labels() map[string]string                          // 返回所有标签
	GetComplete() CompleteHandler                       // 获取CompleteHandler
	SetCompleteHandler(handler CompleteHandler)         // 设置CompleteHandler
	GetFinish() FinishHandler                           // 获取FinishHandler
	SetFinish(handler FinishHandler)                    // 设置FinishHandler
	GetBalance() BalanceHandler                         //
	SetBalance(handler BalanceHandler)                  //
	GetUpstreamHostHandler() UpstreamHostHandler        //
	SetUpstreamHostHandler(handler UpstreamHostHandler) //
	RealIP() string                                     // 客户端ip
	LocalIP() net.IP                                    // 本机ip
	LocalAddr() net.Addr                                // 服务器监听的本地地址
	LocalPort() int                                     // 监听端口
	IsCloneable() bool                                  //
	Clone() (EoContext, error)                          // 克隆
}
