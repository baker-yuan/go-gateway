package http_rule

import (
	"context"
	"log"

	v1 "github.com/baker-yuan/go-gateway/go-gateway-admin/api/admin/v1"
	mregistry "github.com/baker-yuan/go-gateway/registry"
	"github.com/baker-yuan/go-gateway/third_party/httprouter"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	rules      []*v1.HttpRule
	routerRepo = httprouter.New()
)

func LoadHttpRule() error {
	discovery := mregistry.NewDiscovery()

	// 请求admin加载所有接口
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///go-gateway-admin"),
		grpc.WithDiscovery(discovery),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	gClient := v1.NewHttpRuleSrvClient(conn)

	req := &v1.HttpRuleReq{
		Page:     proto.Uint32(1),
		PageSize: proto.Uint32(100),
	}
	rsp, err := gClient.ListHttpRule(context.Background(), req)
	if err != nil {
		panic(err)
	}
	allRules := rsp.Rules
	for _, v := range allRules {
		routerRepo.Handle(v.GetHttpType(), v.GetGwUrl(), convert(v))
	}
	return nil
}

func GetRule(method, path string) (*httprouter.ServiceInfo, httprouter.Params) {
	lookup, params, _ := routerRepo.Lookup(method, path)
	return lookup, params
}

func convert(rule *v1.HttpRule) *httprouter.ServiceInfo {
	return &httprouter.ServiceInfo{
		GwUrl:         rule.GetGwUrl(),
		HttpType:      rule.GetHttpType(),
		Application:   rule.GetApplication(),
		Config:        rule.GetConfig(),
		InterfaceUrl:  rule.GetInterfaceUrl(),
		Status:        rule.GetStatus(),
		InterfaceType: rule.GetInterfaceType(),
	}
}
