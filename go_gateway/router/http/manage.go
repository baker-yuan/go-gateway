package router_manager

import (
	"net/http"
	"time"

	pb "github.com/baker-yuan/go-gateway/pb/router"
	gcontext "github.com/baker-yuan/go-gateway/pkg/context"
	http_context "github.com/baker-yuan/go-gateway/pkg/context/http-context"
	http_context_imp "github.com/baker-yuan/go-gateway/pkg/context/http-context/impl"
	pkg_model "github.com/baker-yuan/go-gateway/pkg/model"
	"github.com/baker-yuan/go-gateway/pkg/router"
	pkg_router "github.com/baker-yuan/go-gateway/pkg/router"
	http_router "github.com/baker-yuan/go-gateway/pkg/router/http-router"
	pkg_service "github.com/baker-yuan/go-gateway/pkg/service"
	http_complete "github.com/baker-yuan/go-gateway/router/http/http-complete"
	"github.com/valyala/fasthttp"
)

var (
	routerNotEnable = pkg_model.Result{
		Code:    http.StatusInternalServerError,
		Message: "路由被禁用",
	}
)

// IRouterManager 路由管理器统一接口
type IRouterManager interface {
	http_context.HandHttp
	Set(router *pb.HttpRouter) // 设置路由
	Delete(id uint32)          // 删除路由
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

	// 创建拦击器链
	chainCreate gcontext.IChainCreate
	serviceGet  pkg_service.IServiceGet
}

func NewRouterManager(chainCreate gcontext.IChainCreate, serviceGet pkg_service.IServiceGet) IRouterManager {
	return &RouterManager{
		origin:      make(map[uint32]*pb.HttpRouter, 0),
		router:      make(map[uint32]*Router, 0),
		chainCreate: chainCreate,
		serviceGet:  serviceGet,
	}
}

func (m *RouterManager) Set(cfg *pb.HttpRouter) {
	// 一个路由一个httpHandler，用于保存http请求执行需要用到的数据
	handler := &httpHandler{
		routerID:  cfg.GetId(),
		serviceID: cfg.GetServiceId(),
		finisher:  defaultFinisher,
		disable:   cfg.GetDisable(),
		retry:     cfg.GetRetry(),
		timeout:   time.Duration(cfg.GetTimeOut()) * time.Millisecond,
	}

	// 路由设置
	if cfg.GetDisable() {
		handler.completeHandler = http_complete.NewFailCompleteHandler(http.StatusInternalServerError, nil, routerNotEnable.ToString())
	} else {
		var plugins gcontext.IChainPro
		if cfg.PluginTemplateId != nil {
			// 插件
		} else {
			plugins = m.chainCreate.CreateChain(cfg.Plugins)
		}
		handler.filters = plugins

		// 服务设置
		handler.service = m.serviceGet.GetService(cfg.GetServiceId())
	}

	// 精细化匹配规则
	appendRule := make([]pkg_router.AppendRule, 0, len(cfg.GetRules()))
	for _, r := range cfg.Rules {
		appendRule = append(appendRule, pkg_router.AppendRule{
			Type:    r.GetType(),
			Name:    r.GetName(),
			Pattern: r.GetValue(),
		})
	}

	r := &Router{
		ID:          cfg.GetId(),
		Hosts:       cfg.Host,
		Method:      cfg.Method,
		Path:        cfg.GetLocation(),
		Appends:     appendRule,
		HttpHandler: handler,
	}

	m.origin[r.ID] = cfg
	m.router[r.ID] = r
	m.matcher, _ = parse(m.router)
}

func (m RouterManager) Delete(id uint32) {

}

func parse(router map[uint32]*Router) (router.IMatcher, error) {
	root := http_router.NewRoot()
	for _, v := range router {
		err := root.Add(v.ID, v.HttpHandler, defaultPort, v.Hosts, v.Method, v.Path, v.Appends)
		if err != nil {
			return nil, err
		}
	}
	return root.Build(), nil
}

func (m RouterManager) FastHandler(ctx *fasthttp.RequestCtx) {
	httpContext := http_context_imp.NewContext(ctx, defaultPort)
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
