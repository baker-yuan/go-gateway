package service

import (
	"context"

	v1 "github.com/baker-yuan/go-gateway/go-gateway-admin/api/admin/v1"
	"github.com/baker-yuan/go-gateway/go-gateway-admin/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/proto"
)

// HttpRuleService is a httpRule service.
type HttpRuleService struct {
	v1.UnimplementedGatewayAdminServer
	rule *biz.HttpRuleBiz
}

// NewHttpRuleService new a httpRule service.
func NewHttpRuleService(rule *biz.HttpRuleBiz, logger log.Logger) *HttpRuleService {
	return &HttpRuleService{rule: rule}
}

func (s *HttpRuleService) ListHttpRule(ctx context.Context, req *v1.ListHttpRuleReq) (*v1.HttpRulesRsp, error) {
	rules, total, err := s.rule.ListHttpRule(ctx, req.GetPage(), req.GetPageSize(), req)
	if err != nil {
		return nil, err
	}
	res := &v1.HttpRulesRsp{
		Total: proto.Uint32(total),
	}
	for _, v := range rules {
		res.Rules = append(res.Rules, httpRule2pb(v))
	}
	return res, nil
}

func (s *HttpRuleService) AddHttpRule(ctx context.Context, req *v1.HttpRule) (*v1.HttpRuleRsp, error) {
	err := s.rule.AddHttpRule(ctx, pbtHttpRule(req))
	return &v1.HttpRuleRsp{}, err
}

func (s *HttpRuleService) DeleteHttpRule(ctx context.Context, req *v1.DeleteHttpRuleReq) (*v1.HttpRuleRsp, error) {
	err := s.rule.DeleteHttpRule(ctx, req.GetId())
	return &v1.HttpRuleRsp{}, err
}

func (s *HttpRuleService) UpdateHttpRule(ctx context.Context, req *v1.HttpRule) (*v1.HttpRuleRsp, error) {
	err := s.rule.UpdateHttpRule(ctx, pbtHttpRule(req))
	return &v1.HttpRuleRsp{}, err
}

func pbtHttpRule(rule *v1.HttpRule) *biz.HttpRule {
	return &biz.HttpRule{
		ID:            rule.GetId(),
		CreateTime:    rule.GetCreateTime(),
		UpdateTime:    rule.GetUpdateTime(),
		GwURL:         rule.GetGwUrl(),
		HTTPType:      rule.GetHttpType(),
		Status:        rule.GetStatus(),
		Application:   rule.GetApplication(),
		InterfaceType: rule.GetInterfaceType(),
		Config:        rule.GetConfig(),
		InterfaceURL:  rule.GetInterfaceUrl(),
	}
}

func httpRule2pb(rule *biz.HttpRule) *v1.HttpRule {
	return &v1.HttpRule{
		Id:            proto.Uint32(rule.ID),
		GwUrl:         proto.String(rule.GwURL),
		HttpType:      proto.String(rule.HTTPType),
		Application:   proto.String(rule.Application),
		InterfaceType: rule.InterfaceType.Enum(),
		Config:        proto.String(rule.Config),
		InterfaceUrl:  proto.String(rule.InterfaceURL),
		Status:        rule.Status.Enum(),
		CreateTime:    proto.Uint32(rule.CreateTime),
		UpdateTime:    proto.Uint32(rule.UpdateTime),
	}
}
