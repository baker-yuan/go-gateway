package data

import (
	"context"
	"time"

	v1 "github.com/baker-yuan/go-gateway/go-gateway-admin/api/admin/v1"
	"github.com/baker-yuan/go-gateway/go-gateway-admin/internal/biz"
	"github.com/baker-yuan/go-gateway/go-gateway-admin/internal/data/ent"
	"github.com/baker-yuan/go-gateway/go-gateway-admin/internal/data/ent/httprule"
	"github.com/baker-yuan/go-gateway/go-gateway-admin/internal/data/ent/predicate"
	"github.com/go-kratos/kratos/v2/log"
)

type httpRuleRepo struct {
	data *Data
	log  *log.Helper
}

func NewHttpRuleRepo(data *Data, logger log.Logger) biz.HttpRuleRepo {
	return &httpRuleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r httpRuleRepo) ListHttpRule(ctx context.Context, page uint32, pageSize uint32, search *v1.ListHttpRuleReq) ([]*biz.HttpRule, uint32, error) {
	// 定义查询条件
	rules := make([]predicate.HttpRule, 0)
	// if search.Status != nil {
	// 	rules = append(rules, httprule.ApplicationEQ(search.GetStatus()))
	// }
	if search != nil && search.Application != nil {
		rules = append(rules, httprule.ApplicationEQ(search.GetApplication()))
	}
	where := httprule.And(rules...)

	// 查询总数
	total, err := r.data.db.HttpRule.Query().Where(where).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 计算页数和偏移量
	offset := (page - 1) * pageSize

	// 查询一页的数据
	httpRules, err := r.data.db.HttpRule.Query().
		Where(where).
		Offset(int(offset)).
		Limit(int(pageSize)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return data2Ent(httpRules), uint32(total), nil
}

func (r httpRuleRepo) AddHttpRule(ctx context.Context, rule *biz.HttpRule) error {
	_, err := r.data.db.HttpRule.Create().
		SetGwURL(rule.GwURL).
		SetHTTPType(rule.HTTPType).
		SetStatus(uint8(rule.Status)).
		SetApplication(rule.Application).
		SetInterfaceType(uint8(rule.InterfaceType)).
		SetConfig(rule.Config).
		SetInterfaceURL(rule.InterfaceURL).
		SetUpdateTime(uint32(time.Now().Unix())).
		SetCreateTime(uint32(time.Now().Unix())).
		Save(ctx)
	return err
}

func (r httpRuleRepo) DeleteHttpRule(ctx context.Context, id uint32) error {
	return r.data.db.HttpRule.
		DeleteOneID(id).
		Exec(ctx)
}

func (r httpRuleRepo) UpdateHttpRule(ctx context.Context, rule *biz.HttpRule) error {
	_, err := r.data.db.HttpRule.
		UpdateOneID(rule.ID).
		SetGwURL(rule.GwURL).
		SetHTTPType(rule.HTTPType).
		SetStatus(uint8(rule.Status)).
		SetApplication(rule.Application).
		SetInterfaceType(uint8(rule.InterfaceType)).
		SetConfig(rule.Config).
		SetInterfaceURL(rule.InterfaceURL).
		SetUpdateTime(uint32(time.Now().Unix())).
		Save(ctx)
	return err
}

func data2Ent(rules []*ent.HttpRule) []*biz.HttpRule {
	res := make([]*biz.HttpRule, 0)
	for _, v := range rules {
		r := &biz.HttpRule{
			ID:            v.ID,
			CreateTime:    v.CreateTime,
			UpdateTime:    v.UpdateTime,
			GwURL:         v.GwURL,
			HTTPType:      v.HTTPType,
			Status:        v1.Status(v.Status),
			Application:   v.Application,
			InterfaceType: v1.InterfaceType(v.InterfaceType),
			Config:        v.Config,
			InterfaceURL:  v.InterfaceURL,
		}
		res = append(res, r)
	}
	return res
}
