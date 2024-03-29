package context

import "time"

// BalanceHandler 负载均衡
type BalanceHandler interface {
	IService                                           // 服务
	Select(ctx GatewayContext) (IInstance, int, error) // 选择一个接口
	Scheme() string                                    // 负载均衡类型
	TimeOut() time.Duration                            // 超时时间
}
