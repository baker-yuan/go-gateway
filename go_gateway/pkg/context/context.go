package context

import (
	"context"
	"errors"
	"net"
	"time"
)

var ErrEoCtxUnCloneable = errors.New("EoContext is UnCloneable. ")

// CompleteHandler 完成转发请求操作
// 主要作用：
// 1、转发到下游http、grpc、double服务
// 2、路由失败，返回路由失败异常
type CompleteHandler interface {
	Complete(ctx GatewayContext) error
}

// FinishHandler 结束请求操作，请求执行完后，做资源清理用
type FinishHandler interface {
	Finish(ctx GatewayContext) error
}

// GatewayContext 上下文，不同的协议实现自己的上下文
type GatewayContext interface {
	RequestId() string                          // 请求id唯一，每次请求随机生成
	AcceptTime() time.Time                      // 请求接收时间
	Context() context.Context                   // 原始context
	Value(key interface{}) interface{}          // 从原始context中返回键对应的值
	WithValue(key, val interface{})             // 往原始context添加键值对
	Scheme() string                             // 协议 http、https、grpc、dubbo
	Assert(i interface{}) error                 // context类型断言
	SetLabel(name, value string)                // 设置标签
	GetLabel(name string) string                // 获取标签
	Labels() map[string]string                  // 返回所有标签
	GetComplete() CompleteHandler               // 获取CompleteHandler
	SetCompleteHandler(handler CompleteHandler) // 设置CompleteHandler
	GetFinish() FinishHandler                   // 获取FinishHandler
	SetFinish(handler FinishHandler)            // 设置FinishHandler
	GetBalance() BalanceHandler                 // 获取负载均衡器
	SetBalance(handler BalanceHandler)          // 设置负载均衡器
	RealIP() string                             // 客户端IP
	LocalIP() net.IP                            // 本机IP
	LocalAddr() net.Addr                        // 服务器监听的本地地址
	LocalPort() int                             // 监听端口
	IsCloneable() bool                          // 是否可克隆
	Clone() (GatewayContext, error)             // 克隆
}
