package gateway

import (
	"fmt"
	"net"

	"github.com/baker-yuan/go-gateway/pkg/config"
	plugin_manager "github.com/baker-yuan/go-gateway/plugin"
	http_router "github.com/baker-yuan/go-gateway/router/http"
	service_manager "github.com/baker-yuan/go-gateway/service"
	"github.com/valyala/fasthttp"
)

// Engine 网关引擎
type Engine struct {
	config         *config.GatewayConfig
	pluginManager  plugin_manager.IPluginManager
	serviceManager service_manager.IServiceManager

	// 接收http请求
	httpRouteManager http_router.IRouterManager
	httpServer       *fasthttp.Server
}

// New 创建网关实例
func New() (*Engine, error) {
	engine := &Engine{}

	// 远程拉取管理端数据

	// 初始化配置
	conf, err := config.Load()
	if err != nil {
		return nil, err
	}
	engine.config = conf

	// 初始化插件管理器
	pluginManager := plugin_manager.NewPluginManager()
	engine.pluginManager = pluginManager

	// 初始化服务管理器
	serviceManager := service_manager.NewServiceManager()
	engine.serviceManager = serviceManager

	// 初始化http路由管理器
	httpRouteManager := http_router.NewRouterManager(pluginManager, serviceManager)
	engine.httpRouteManager = httpRouteManager

	// 初始化http服务器
	httpServe := &fasthttp.Server{
		StreamRequestBody:            true,
		DisablePreParseMultipartForm: true,
		MaxRequestBodySize:           100 * 1024 * 1024,
		ReadBufferSize:               16 * 1024,
		Handler: func(ctx *fasthttp.RequestCtx) {
			httpRouteManager.FastHandler(ctx)
		}}
	engine.httpServer = httpServe

	return engine, nil
}

// Start 启动网关
func (e *Engine) Start() error {

	if err := e.startHttpServer(); err != nil {
		return err
	}

	return nil
}

func (e *Engine) startHttpServer() error {
	httpCfg := e.config.Proxy.Http
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", httpCfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %+v", httpCfg.Port, err)
	}
	return e.httpServer.Serve(ln)
}
