package ip_restriction

import (
	"net"
	"testing"

	http_service "github.com/baker-yuan/go-gateway/context/http-context"
	http_context "github.com/baker-yuan/go-gateway/context/http-context/impl"
	"github.com/valyala/fasthttp"
)

// 127.0.0.1:8080
var ctx http_service.IHttpContext

// 初始化context，设置发送请求的远程地址为address
func getContext(address string) (http_service.IHttpContext, error) {
	if ctx == nil {
		return initTestContext(address)
	}
	if address == ctx.Request().RemoteAddr() {
		return ctx, nil
	}
	return initTestContext(address)
}

func initTestContext(address string) (http_service.IHttpContext, error) {
	fast := &fasthttp.RequestCtx{}
	freq := fasthttp.AcquireRequest()
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	fast.Init(freq, addr, nil)
	return http_context.NewContext(fast, 0), nil
}

func TestDoRestriction(t *testing.T) {
	// 创建context
	http_ctx, err := getContext("127.0.0.1:8080")
	if err != nil {
		t.Fatal(err)
	}

	// 创建工厂
	factory := NewFactory()
	driver, err := factory.Create("plugin@setting", "ip_restriction", "ip_restriction", "service", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name   string
		config *Config
		want   string
	}{
		{
			// 全加入黑名单
			name: "limit_black_all",
			config: &Config{
				IPListType: "black",
				IPBlackList: []string{
					"*",
				},
			},
			want: "403",
		},
		{
			// 127.0.0.1 加入黑名单
			// 拒绝
			name: "limit_black",
			config: &Config{
				IPListType: "black",
				IPBlackList: []string{
					"127.0.0.1",
				},
			},
			want: "403",
		},
		{
			// 127.0.0.2 加入黑名单
			// 127.0.0.1通过
			name: "pass_black",
			config: &Config{
				IPListType: "black",
				IPBlackList: []string{
					"127.0.0.2",
				},
			},
			want: "200",
		},
		{
			// 127.0.0.1加入ip白名单列表
			name: "limit_white",
			config: &Config{
				IPListType: "black",
				IPWhiteList: []string{
					"127.0.0.1",
				},
			},
			want: "200",
		},
		{
			// 127.0.0.2加入ip白名单列表
			name: "pass_white",
			config: &Config{
				IPListType: "white",
				IPWhiteList: []string{
					"127.0.0.2",
				},
			},
			want: "403",
		},
		{
			// 全都加入白名单放开所有
			name: "pass_white_all",
			config: &Config{
				IPListType: "white",
				IPWhiteList: []string{
					"*",
				},
			},
			want: "200",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			// 重置响应200
			http_ctx.Response().SetStatus(200, "200")

			// 创建具体的插件
			ip, err := driver.Create("ip_restriction@plugin", "ip_restriction", cc.config, nil)
			if err != nil {
				t.Errorf("create handler error : %v", err)
			}

			// 类型转换
			filter, ok := ip.(http_service.HttpFilter)
			if !ok {
				t.Errorf("parse filter error")
				return
			}
			// 执行拦击器逻辑
			_ = filter.DoHttpFilter(http_ctx, nil)

			// 判断
			if http_ctx.Response().Status() != cc.want {
				t.Errorf("do restriction error; want %s, got %s", cc.want, http_ctx.Response().Status())
			}
		})
	}

}
