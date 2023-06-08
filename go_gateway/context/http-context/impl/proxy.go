package http_context

import (
	"bytes"
	"fmt"

	http_service "github.com/baker-yuan/go-gateway/context/http-context"
	"github.com/baker-yuan/go-gateway/log"
	"github.com/valyala/fasthttp"
)

var _ http_service.IRequest = (*ProxyRequest)(nil)

type ProxyRequest struct {
	RequestReader
}

func (r *ProxyRequest) Finish() error {
	err := r.RequestReader.Finish()
	if err != nil {
		log.Warn(err)
	}
	return nil
}

func (r *ProxyRequest) Header() http_service.IHeaderWriter {
	return &r.headers
}

func (r *ProxyRequest) Body() http_service.IBodyDataWriter {
	return &r.body
}

func (r *ProxyRequest) URI() http_service.IURIWriter {
	return &r.uri
}

var (
	xforwardedforKey = []byte("x-forwarded-for")
)

func (r *ProxyRequest) reset(request *fasthttp.Request, remoteAddr string) {
	r.req = request
	forwardedFor := r.req.Header.PeekBytes(xforwardedforKey)
	if len(forwardedFor) > 0 {
		if i := bytes.IndexByte(forwardedFor, ','); i > 0 {
			r.realIP = string(forwardedFor[:i])
		} else {
			r.realIP = string(forwardedFor)
		}
		r.req.Header.Set("x-forwarded-for", fmt.Sprint(string(forwardedFor), ",", r.remoteAddr))
	} else {
		r.req.Header.Set("x-forwarded-for", r.remoteAddr)
		r.realIP = r.remoteAddr
	}
	r.RequestReader.reset(r.req, remoteAddr)
}

func (r *ProxyRequest) SetMethod(s string) {
	r.Request().Header.SetMethod(s)
}
