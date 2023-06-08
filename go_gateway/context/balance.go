package eocontext

import "time"

// BalanceHandler 负载均衡
type BalanceHandler interface {
	Select(ctx EoContext) (INode, int, error) // 选择一个接口
	Scheme() string                           // 负载均衡类型
	TimeOut() time.Duration                   // 超时时间
	EoApp                                     // 网关app
}
