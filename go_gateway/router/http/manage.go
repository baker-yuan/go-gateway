package router_manager

import (
	"time"

	pb "github.com/baker-yuan/go-gateway/pb/router"
	gocontext "github.com/baker-yuan/go-gateway/pkg/context"
	http_context "github.com/baker-yuan/go-gateway/pkg/context/http-context/impl"
	"github.com/baker-yuan/go-gateway/pkg/router"
	pkg_router "github.com/baker-yuan/go-gateway/pkg/router"
	http_router "github.com/baker-yuan/go-gateway/pkg/router/http-router"
	plugin_manager "github.com/baker-yuan/go-gateway/plugin"
	service_manager "github.com/baker-yuan/go-gateway/service"
	"github.com/valyala/fasthttp"
)

// IRouterManager 路由管理器统一接口
type IRouterManager interface {
	Set(router *pb.HttpRouter, serviceManager service_manager.IServiceManager) // 设置路由
	Delete(id uint32)                                                          // 删除路由
	FastHandler(ctx *fasthttp.RequestCtx)                                      // 处理http请求
}

var (
	defaultPort = 0
)

// Router 路由数据
type Router struct {
	ID          uint32                  // 路由ID
	Hosts       []string                // 请求host
	Method      []string                // 请求方式
	Path        string                  // 请求路径
	Appends     []pkg_router.AppendRule // 规则
	HttpHandler router.IRouterHandler   // 处理请求
}

type RouterManager struct {
	origin  map[uint32]*pb.HttpRouter
	router  map[uint32]*Router
	matcher router.IMatcher

	pluginManager  plugin_manager.IPluginManager
	serviceManager service_manager.IServiceManager
}

func NewRouterManager(pluginManager plugin_manager.IPluginManager, serviceManager service_manager.IServiceManager) IRouterManager {
	return &RouterManager{
		origin:         make(map[uint32]*pb.HttpRouter, 0),
		router:         make(map[uint32]*Router, 0),
		pluginManager:  pluginManager,
		serviceManager: serviceManager,
	}
}

func (m *RouterManager) Set(router *pb.HttpRouter, serviceManager service_manager.IServiceManager) {
	// 一个路由一个httpHandler，用于串联http请求执行
	handler := &httpHandler{
		routerID:  router.GetId(),
		serviceID: router.GetServiceId(),
		finisher:  defaultFinisher,
		disable:   router.GetDisable(),
		retry:     router.GetRetry(),
		timeout:   time.Duration(router.GetTimeOut()) * time.Millisecond,
	}

	// 路由设置
	var plugins gocontext.IChainPro
	if !router.GetDisable() {
		if router.PluginTemplateId != nil {
		} else {
			plugins = m.pluginManager.CreateRequest(router.Plugins)
		}
	}
	handler.filters = plugins

	// 服务设置

	// 精细化匹配规则
	appendRule := make([]pkg_router.AppendRule, 0, len(router.GetRules()))
	for _, r := range router.Rules {
		appendRule = append(appendRule, pkg_router.AppendRule{
			Type:    r.GetType(),
			Name:    r.GetName(),
			Pattern: r.GetValue(),
		})
	}

	r := &Router{
		ID:          router.GetId(),
		Hosts:       router.Host,
		Method:      router.Method,
		Path:        router.GetLocation(),
		Appends:     appendRule,
		HttpHandler: handler,
	}

	m.router[r.ID] = r

	m.matcher, _ = m.Parse()
}

func (m RouterManager) Delete(id uint32) {

}

func (m *RouterManager) Parse() (router.IMatcher, error) {
	root := http_router.NewRoot()
	for _, v := range m.router {
		err := root.Add(v.ID, v.HttpHandler, defaultPort, v.Hosts, v.Method, v.Path, v.Appends)
		if err != nil {
			return nil, err
		}
	}
	return root.Build(), nil
}

func (m RouterManager) FastHandler(ctx *fasthttp.RequestCtx) {
	httpContext := http_context.NewContext(ctx, defaultPort)
	// if m.matcher == nil {
	// 	httpContext.SetFinish(notFound)
	// 	httpContext.SetCompleteHandler(notFound)
	// 	globalFilters := m.globalFilters.Load()
	// 	if globalFilters != nil {
	// 		(*globalFilters).Chain(httpContext, completeCaller)
	// 	}
	// 	return
	// }

	// 匹配路由
	routerHandler, has := m.matcher.Match(defaultPort, httpContext.Request())
	if !has {
		// 匹配失败，返回404
		// httpContext.SetFinish(notFound)
		// httpContext.SetCompleteHandler(notFound)
		// globalFilters := m.globalFilters.Load()
		// if globalFilters != nil {
		// 	(*globalFilters).Chain(httpContext, completeCaller)
		// }
	} else {
		// 匹配成功
		routerHandler.ServeHTTP(httpContext)
	}

	finishHandler := httpContext.GetFinish()
	if finishHandler != nil {
		finishHandler.Finish(httpContext)
	}

}
