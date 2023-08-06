package http_context

import (
	"strconv"
	"time"

	http_service "github.com/baker-yuan/go-gateway/pkg/context/http-context"
)

var _ http_service.IProxy = (*requestAgent)(nil)

// requestAgent
// go_gateway/context/http-context/context.go#IProxy
// go_gateway/context/http-context/context.go#IRequest
type requestAgent struct {
	http_service.IRequest           //
	host                  string    //
	scheme                string    // http协议
	statusCode            int       //
	status                string    //
	responseLength        int       //
	beginTime             time.Time // 请求执行时间
	endTime               time.Time // 请求结束时间
	hostAgent             *UrlAgent //
}

func newRequestAgent(IRequest http_service.IRequest, host string, scheme string, beginTime, endTime time.Time) *requestAgent {
	return &requestAgent{
		IRequest:  IRequest,
		host:      host,
		scheme:    scheme,
		beginTime: beginTime,
		endTime:   endTime,
	}
}

// ------------------- go_gateway/context/http-context/context.go#IProxy

func (a *requestAgent) StatusCode() int {
	return a.statusCode
}

func (a *requestAgent) Status() string {
	return a.status
}

func (a *requestAgent) ProxyTime() time.Time {
	return a.beginTime
}

func (a *requestAgent) ResponseLength() int {
	return a.responseLength
}

func (a *requestAgent) ResponseTime() int64 {
	return a.endTime.Sub(a.beginTime).Milliseconds()
}

// ------------------- go_gateway/context/http-context/context.go#IRequest

func (a *requestAgent) URI() http_service.IURIWriter {
	if a.hostAgent == nil {
		a.hostAgent = NewUrlAgent(a.IRequest.URI(), a.host, a.scheme)
	}
	return a.hostAgent
}

func (a *requestAgent) setStatusCode(code int) {
	a.statusCode = code
	a.status = strconv.Itoa(code)
}

func (a *requestAgent) setResponseLength(length int) {
	if length > 0 {
		a.responseLength = length
	}
}

// -------------------

type UrlAgent struct {
	http_service.IURIWriter
	host   string
	scheme string
}

func NewUrlAgent(IURIWriter http_service.IURIWriter, host string, scheme string) *UrlAgent {
	return &UrlAgent{IURIWriter: IURIWriter, host: host, scheme: scheme}
}

func (u *UrlAgent) SetScheme(scheme string) {
	u.scheme = scheme
}
func (u *UrlAgent) Scheme() string {
	return u.scheme
}

func (u *UrlAgent) Host() string {
	return u.host
}

func (u *UrlAgent) SetHost(host string) {
	u.host = host
}
