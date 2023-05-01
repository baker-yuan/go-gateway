package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// HttpRule holds the schema definition for the HttpRule entity.
type HttpRule struct {
	ent.Schema
}

// Fields of the HttpRule.
func (HttpRule) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("id"),

		field.String("gw_url").MaxLen(128).Default("").NotEmpty().Comment("网关接口"),
		field.String("http_type").MaxLen(128).Default("").NotEmpty().Comment("接口类型"),
		field.Uint8("status").Default(0).Comment("接口状态 0-默认 1-上线 2-下线"),

		field.String("application").MaxLen(128).NotEmpty().Comment("应用名称"),
		field.Uint8("interface_type").Default(0).Comment("接口协议 0-未知 1-http 2-https 3-gRPC 4-Double"),
		field.String("interface_url").MaxLen(128).Default("").NotEmpty().Comment("接口方法"),
		field.String("config").MaxLen(2000).Default("").NotEmpty().Comment("接口特殊配置json格式"),

		field.Uint32("create_time").Default(0).Comment("创建时间"),
		field.Uint32("update_time").Default(0).Comment("修改时间"),
	}
}

// Edges of the HttpRule.
func (HttpRule) Edges() []ent.Edge {
	return nil
}
