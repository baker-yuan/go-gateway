package http_context

import (
	"strconv"
	"strings"
	"time"

	http_service "github.com/baker-yuan/go-gateway/context/http-context"
	"github.com/valyala/fasthttp"
)

var _ http_service.IResponse = (*Response)(nil)

// Response
// go_gateway/context/http-context/context.go#IBodyGet
// go_gateway/context/http-context/context.go#IResponse
// go_gateway/context/http-context/context.go#IResponseHeader
// go_gateway/context/http-context/context.go#IBodySet
type Response struct {
	ResponseHeader
	*fasthttp.Response
	length          int
	responseTime    time.Duration
	proxyStatusCode int
	responseError   error
}

// --------------- go_gateway/context/http-context/context.go#IBodyGet

func (r *Response) GetBody() []byte {
	if strings.Contains(r.GetHeader("Content-Encoding"), "gzip") {
		body, _ := r.BodyGunzip()
		r.DelHeader("Content-Encoding")
		r.SetHeader("Content-Length", strconv.Itoa(len(body)))
		r.Response.SetBody(body)
	}
	return r.Response.Body()
}

func (r *Response) BodyLen() int {
	return r.header.ContentLength()
}

// ---------------  go_gateway/context/http-context/context.go#IResponse

func (r *Response) ResponseError() error {
	return r.responseError
}

func (r *Response) ClearError() {
	r.responseError = nil
}

func (r *Response) SetResponseTime(t time.Duration) {
	r.responseTime = t
}

func (r *Response) ResponseTime() time.Duration {
	return r.responseTime
}

func (r *Response) ContentLength() int {
	if r.length == 0 {
		return r.Response.Header.ContentLength()
	}
	return r.length
}

func (r *Response) ContentType() string {
	return string(r.Response.Header.ContentType())
}

// --------------- go_gateway/context/http-context/context.go#IResponseHeader

func (r *Response) HeadersString() string {
	return r.header.String()
}

// --------------- go_gateway/context/http-context/context.go#IBodySet

func (r *Response) SetBody(bytes []byte) {
	if strings.Contains(r.GetHeader("Content-Encoding"), "gzip") {
		r.DelHeader("Content-Encoding")
	}
	r.Response.SetBody(bytes)
	r.length = len(bytes)
	r.SetHeader("Content-Length", strconv.Itoa(r.length))
	r.responseError = nil
}

// --------------- go_gateway/context/http-context/context.go#IStatusGet

func (r *Response) StatusCode() int {
	if r.responseError != nil {
		return 504
	}
	return r.Response.StatusCode()
}

// ProxyStatusCode 原始的响应状态码
func (r *Response) ProxyStatusCode() int {
	return r.proxyStatusCode
}

func (r *Response) ProxyStatus() string {
	return strconv.Itoa(r.proxyStatusCode)
}

func (r *Response) Status() string {
	return strconv.Itoa(r.StatusCode())
}

// --------------- go_gateway/context/http-context/context.go#IStatusSet

func (r *Response) SetStatus(code int, status string) {
	r.Response.SetStatusCode(code)
	r.responseError = nil
}

func (r *Response) SetProxyStatus(code int, status string) {
	r.proxyStatusCode = code
}

// ---------------

func (r *Response) Finish() error {
	r.ResponseHeader.Finish()
	r.Response = nil
	r.responseError = nil
	r.proxyStatusCode = 0
	return nil
}
func (r *Response) reset(resp *fasthttp.Response) {
	r.Response = resp
	r.ResponseHeader.reset(&resp.Header)
	r.responseError = nil
	r.proxyStatusCode = 0
}
