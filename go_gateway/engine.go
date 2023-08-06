package gateway

import (
	"fmt"
	"net"

	"github.com/baker-yuan/go-gateway/pkg/config"
	"github.com/baker-yuan/go-gateway/plugin"
	http_router "github.com/baker-yuan/go-gateway/router/http"
	"github.com/baker-yuan/go-gateway/service"
	"github.com/valyala/fasthttp"
)

// Engine 网关引擎
type Engine struct {
	config         *config.GatewayConfig
	pluginManager  plugin.PluginManager
	serviceManager service.ServiceManager

	// 接收http请求
	httpRouteManager http_router.RouterManager
	httpServer       *fasthttp.Server
}

// New 创建网关实例
func New() (*Engine, error) {
	engine := &Engine{}

	// 初始化配置
	conf, err := config.Load()
	if err != nil {
		return nil, err
	}
	engine.config = conf

	// 初始化插件管理器
	engine.pluginManager = plugin.NewPluginManager()

	// 初始化服务管理器
	engine.serviceManager = service.NewServiceManager()

	// 初始化http路由管理器
	engine.httpRouteManager = http_router.NewRouterManager()

	// 初始化http服务器
	httpServe := &fasthttp.Server{
		StreamRequestBody:            true,
		DisablePreParseMultipartForm: true,
		MaxRequestBodySize:           100 * 1024 * 1024,
		ReadBufferSize:               16 * 1024,
		Handler: func(ctx *fasthttp.RequestCtx) {
			engine.httpRouteManager.FastHandler(ctx)
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
