package eocontext

import (
	"context"
	"errors"
	"net"
	"time"
)

var ErrEoCtxUnCloneable = errors.New("EoContext is UnCloneable. ")

// CompleteHandler 所有的插件执行完成后，转发到下游服务/返回失败的逻辑
type CompleteHandler interface {
	Complete(ctx EoContext) error
}

// FinishHandler 请求执行完后，做清理的
type FinishHandler interface {
	Finish(ctx EoContext) error
}

// EoContext 上下文，不通的协议实现自己的上下文
type EoContext interface {
	RequestId() string                 // 请求id唯一
	AcceptTime() time.Time             // 请求接收时间
	Context() context.Context          // 原始context
	Value(key interface{}) interface{} //
	WithValue(key, val interface{})    // put k v

	Scheme() string             // 协议
	Assert(i interface{}) error // 类型断言

	SetLabel(name, value string) // 设置标签
	GetLabel(name string) string // 获取标签
	Labels() map[string]string   // 返回所有标签

	GetComplete() CompleteHandler                       //
	SetCompleteHandler(handler CompleteHandler)         // 设置CompleteHandler
	GetFinish() FinishHandler                           //
	SetFinish(handler FinishHandler)                    //
	GetBalance() BalanceHandler                         //
	SetBalance(handler BalanceHandler)                  //
	GetUpstreamHostHandler() UpstreamHostHandler        //
	SetUpstreamHostHandler(handler UpstreamHostHandler) //

	RealIP() string      //
	LocalIP() net.IP     //
	LocalAddr() net.Addr //
	LocalPort() int      //

	IsCloneable() bool         //
	Clone() (EoContext, error) //
}
