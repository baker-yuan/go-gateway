package biz

import (
	"context"
	"net/http"

	v1 "github.com/baker-yuan/go-gateway/go-gateway-admin/api/admin/v1"
	"github.com/go-kratos/kratos/v2/log"
)

// HttpRule 网关接口映射信息
type HttpRule struct {
	ID         uint32 // ID
	CreateTime uint32 // 创建时间
	UpdateTime uint32 // 修改时间
	// 网关
	GwURL    string    // 网关接口路径
	HTTPType string    // 接口类型 net/http/method.go
	Status   v1.Status // 接口状态
	// 接口
	Application   string           // 应用名称
	InterfaceType v1.InterfaceType // 接口协议
	Config        string           // 指定协议的配置
	InterfaceURL  string           // 接口路径
}

type HttpRuleRepo interface {
	// db
	ListHttpRule(ctx context.Context, page uint32, pageSize uint32, search *v1.ListHttpRuleReq) ([]*HttpRule, uint32, error)
	AddHttpRule(ctx context.Context, rule *HttpRule) error
	DeleteHttpRule(ctx context.Context, id uint32) error
	UpdateHttpRule(ctx context.Context, rule *HttpRule) error
}

type HttpRuleBiz struct {
	repo   HttpRuleRepo
	logger log.Logger
}

func NewHttpRuleBiz(repo HttpRuleRepo, logger log.Logger) *HttpRuleBiz {
	return &HttpRuleBiz{repo: repo}
}

func (r *HttpRuleBiz) ListHttpRule(ctx context.Context, page uint32, size uint32, search *v1.ListHttpRuleReq) ([]*HttpRule, uint32, error) {
	rule, total, err := r.repo.ListHttpRule(ctx, page, size, search)
	if err != nil {
		return nil, 0, err
	}
	return rule, total, err
}

func (r *HttpRuleBiz) AddHttpRule(ctx context.Context, rule *HttpRule) error {
	if rule.InterfaceType == v1.InterfaceType_G_RPC && rule.HTTPType != http.MethodPost {
		return nil // todo
	}
	return r.repo.AddHttpRule(ctx, rule)
}

func (r *HttpRuleBiz) DeleteHttpRule(ctx context.Context, id uint32) error {
	return r.repo.DeleteHttpRule(ctx, id)
}

func (r *HttpRuleBiz) UpdateHttpRule(ctx context.Context, rule *HttpRule) error {
	return r.repo.UpdateHttpRule(ctx, rule)
}
