package http_context

import (
	http_service "github.com/baker-yuan/go-gateway/context/http-context"
	"github.com/valyala/fasthttp"
)

var _ http_service.IURIWriter = (*URIRequest)(nil)

// URIRequest
// go_gateway/context/http-context/context.go#IQueryReader
// go_gateway/context/http-context/context.go#IQueryWriter
// go_gateway/context/http-context/context.go#IURIReader
// go_gateway/context/http-context/context.go#IURIWriter
type URIRequest struct {
	uri *fasthttp.URI
}

func NewURIRequest(uri *fasthttp.URI) *URIRequest {
	return &URIRequest{uri: uri}
}

// ----------------- go_gateway/context/http-context/context.go#IQueryReader

func (ur *URIRequest) GetQuery(key string) string {
	return string(ur.uri.QueryArgs().Peek(key))
}

func (ur *URIRequest) RawQuery() string {
	return string(ur.uri.QueryString())
}

// ----------------- go_gateway/context/http-context/context.go#IQueryWriter

func (ur *URIRequest) SetQuery(key, value string) {
	ur.uri.QueryArgs().Set(key, value)
}

func (ur *URIRequest) AddQuery(key, value string) {
	ur.uri.QueryArgs().Add(key, value)
}

func (ur *URIRequest) DelQuery(key string) {
	queryArgs := ur.uri.QueryArgs()
	queryArgs.Del(key)
	if queryArgs.Len() == 0 {
		ur.uri.SetQueryStringBytes(nil)
	}
}

func (ur *URIRequest) SetRawQuery(raw string) {
	ur.uri.SetQueryString(raw)
}

// ----------------- go_gateway/context/http-context/context.go#IURIReader

func (ur *URIRequest) RequestURI() string {
	return string(ur.uri.RequestURI())
}

func (ur *URIRequest) Scheme() string {
	return string(ur.uri.Scheme())
}

func (ur *URIRequest) RawURL() string {
	return string(ur.uri.FullURI())
}

func (ur *URIRequest) Host() string {
	return string(ur.uri.Host())
}

func (ur *URIRequest) Path() string {
	return string(ur.uri.Path())
}

// ----------------- go_gateway/context/http-context/context.go#IURIWriter

func (ur *URIRequest) SetPath(s string) {
	ur.uri.SetPath(s)

}

func (ur *URIRequest) SetScheme(scheme string) {
	ur.uri.SetScheme(scheme)
}

func (ur *URIRequest) SetHost(host string) {
	ur.uri.SetHost(host)
}

// -----------------

func (ur *URIRequest) reset(uri *fasthttp.URI) {
	ur.uri = uri
}
