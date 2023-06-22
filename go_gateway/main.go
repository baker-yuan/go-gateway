package main

import (
	"fmt"
	"net"

	eoscContext "github.com/baker-yuan/go-gateway/context"
	http_service "github.com/baker-yuan/go-gateway/context/http-context"
	http_context "github.com/baker-yuan/go-gateway/context/http-context/impl"
	ip_restriction "github.com/baker-yuan/go-gateway/drivers/plugins/ip-restriction"
	"github.com/baker-yuan/go-gateway/eosc"
	"github.com/valyala/fasthttp"
)

// https://github.com/eolinker/apinto/pull/113
// https://blog.csdn.net/weixin_41479678/article/details/111933900
func main() {

	register := eosc.NewExtenderRegister()
	ip_restriction.Register(register)

	driver, exist := register.GetDriver(ip_restriction.Name)
	if exist {
		fmt.Println(driver)
	}
	d, err := driver.Create("plugin@setting", "ip_restriction", "ip_restriction", "service", map[string]interface{}{})
	if err != nil {
		panic(err)
	}
	fmt.Println(d)

	// ip, err := d.Create("ip_restriction@plugin", "ip_restriction", cc.config, nil)

	port := 8848

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		panic(err)
	}
	notFound := &HttpNotFoundHandler{}
	server := fasthttp.Server{
		StreamRequestBody:            true,
		DisablePreParseMultipartForm: true,
		MaxRequestBodySize:           100 * 1024 * 1024,
		ReadBufferSize:               16 * 1024,
		Handler: func(ctx *fasthttp.RequestCtx) {
			fmt.Println("handle...start")
			httpContext := http_context.NewContext(ctx, port)
			// header := httpContext.Request().Header()
			// fmt.Println(header)

			_ = notFound.Complete(httpContext)

			// httpContext.SetCompleteHandler(notFound)
			// httpContext.SetFinish(notFound)
			_ = notFound.Finish(httpContext)
			fmt.Println("handle...end")
		}}
	_ = server.Serve(ln)
}

type HttpNotFoundHandler struct {
}

func (m *HttpNotFoundHandler) Complete(ctx eoscContext.EoContext) error {
	httpContext, err := http_service.Assert(ctx)
	if err != nil {
		return nil
	}
	httpContext.Response().SetHeader("Content-Type", "application/json")
	httpContext.Response().SetStatus(200, "200")
	httpContext.Response().SetBody([]byte(`{code:"200",message:"ok"}`))
	return nil
}

func (m *HttpNotFoundHandler) Finish(ctx eoscContext.EoContext) error {
	httpContext, err := http_service.Assert(ctx)
	if err != nil {
		return err
	}
	httpContext.FastFinish()
	return nil
}
