package plugin

// PluginSchema 插件模型
type PluginSchema struct {
	Name       string                                     `json:"name"`        // 插件名称
	JsonSchema string                                     `json:"json_schema"` // 插件配置，json schema 前端通过配置渲染页面
	Priority   int                                        `json:"priority"`    // 优先级，优先级大的先执行
	Creator    func(conf []byte) (IPluginInstance, error) // 创建插件
}

// IPluginInstance 插件接口
type IPluginInstance interface {
}
