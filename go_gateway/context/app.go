package eocontext

// NodeStatus 节点状态类型
type NodeStatus int

const (
	// Running 节点运行中状态
	Running NodeStatus = 1
	// Down 节点不可用状态
	Down NodeStatus = 2
	// Leave 节点离开状态
	Leave NodeStatus = 3
)

// Attrs 属性集合
type Attrs map[string]string

// IAttributes 属性接口
type IAttributes interface {
	GetAttrs() Attrs
	GetAttrByName(name string) (string, bool)
}

// EoApp 网关
type EoApp interface {
	Nodes() []INode // 节点
}

// INode 节点接口
type INode interface {
	IAttributes         // 属性 key=string v=string
	ID() string         //
	IP() string         // IP
	Port() int          // 端口
	Addr() string       // 地址
	Status() NodeStatus // 节点状态
	Up()                //
	Down()              //
	Leave()             //
}
