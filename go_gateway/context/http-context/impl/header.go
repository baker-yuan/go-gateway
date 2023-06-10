package http_context

import (
	"bytes"
	"net/http"
	"strings"

	http_service "github.com/baker-yuan/go-gateway/context/http-context"
	"github.com/valyala/fasthttp"
)

var _ http_service.IHeaderWriter = (*RequestHeader)(nil)

// RequestHeader
// go_gateway/context/http-context/context.go#IHeaderReader
// go_gateway/context/http-context/context.go#IHeaderWriter
type RequestHeader struct {
	header *fasthttp.RequestHeader //
	tmp    http.Header             // map[string][]string
}

// --------------------- go_gateway/context/http-context/context.go#IHeaderReader

func (h *RequestHeader) RawHeader() string {
	return h.header.String()
}

func (h *RequestHeader) GetHeader(name string) string {
	return h.Headers().Get(name)
}

func (h *RequestHeader) Headers() http.Header {
	h.initHeader()
	return h.tmp
}

func (h *RequestHeader) Host() string {
	return string(h.header.Host())
}

func (h *RequestHeader) GetCookie(key string) string {
	return string(h.header.Cookie(key))
}

// ---------------------- go_gateway/context/http-context/context.go#IHeaderWriter

func (h *RequestHeader) SetHeader(key, value string) {
	if h.tmp != nil {
		h.tmp.Set(key, value)
	}
	h.header.Set(key, value)
}

func (h *RequestHeader) AddHeader(key, value string) {
	if h.tmp != nil {
		h.tmp.Add(key, value)
	}
	h.header.Add(key, value)
}

func (h *RequestHeader) DelHeader(key string) {
	if h.tmp != nil {
		h.tmp.Del(key)
	}
	h.header.Del(key)
}

func (h *RequestHeader) SetHost(host string) {
	if h.tmp != nil {
		h.tmp.Set("Host", host)
	}
	h.header.SetHost(host)
}

// --------------------

func (h *RequestHeader) initHeader() {
	if h.tmp == nil {
		h.tmp = make(http.Header)
		h.header.VisitAll(func(key, value []byte) {
			bytes.SplitN(value, []byte(":"), 2)
			h.tmp[string(key)] = []string{string(value)}
		})
	}
}

func (h *RequestHeader) reset(header *fasthttp.RequestHeader) {
	h.header = header
	h.tmp = nil
}

// --------------------

type headerActionHandleFunc func(target *ResponseHeader, key string, value ...string)

var (
	headerActionAdd = func(target *ResponseHeader, key string, value ...string) {
		target.cache.Add(key, value[0])
		target.header.Add(key, value[0])
	}
	headerActionSet = func(target *ResponseHeader, key string, value ...string) {
		target.cache.Set(key, value[0])
		target.header.Set(key, value[0])
	}
	headerActionDel = func(target *ResponseHeader, key string, value ...string) {
		target.cache.Del(key)
		target.header.Del(key)
	}
)

type headerAction struct {
	Action headerActionHandleFunc
	Key    string
	Value  string
}

type ResponseHeader struct {
	header     *fasthttp.ResponseHeader //
	cache      http.Header              // map[string][]string
	actions    []*headerAction          //
	afterProxy bool                     //
}

func (r *ResponseHeader) reset(header *fasthttp.ResponseHeader) {
	r.header = header
	r.cache = http.Header{}
	r.actions = nil
	r.afterProxy = false
}

// refresh 刷新
func (r *ResponseHeader) refresh() {
	tmp := make(http.Header)
	hs := strings.Split(r.header.String(), "\r\n")
	for _, t := range hs {
		if strings.TrimSpace(t) == "" {
			continue
		}
		vs := strings.Split(t, ":")
		if len(vs) < 2 {
			if vs[0] == "" {
				continue
			}
			tmp[vs[0]] = []string{""}
			continue
		}
		tmp[vs[0]] = []string{strings.TrimSpace(vs[1])}
	}
	r.cache = tmp
	for _, ac := range r.actions {
		ac.Action(r, ac.Key, ac.Value)
	}
	r.afterProxy = true
	r.actions = nil
}

func (r *ResponseHeader) Finish() {
	r.header = nil
	r.cache = nil
	r.actions = nil
}

func (r *ResponseHeader) GetHeader(name string) string {
	return r.Headers().Get(name)
}

func (r *ResponseHeader) Headers() http.Header {
	return r.cache
}

func (r *ResponseHeader) SetHeader(key, value string) {
	// r.cache.Set(key, value)
	// if r.afterProxy {
	// 	r.header.Set(key, value)
	// } else {
	// 	r.actions = append(r.actions, &headerAction{
	// 		Key:    key,
	// 		Value:  value,
	// 		Action: headerActionSet,
	// 	})
	// }

	r.cache.Set(key, value)
	r.header.Set(key, value)
	if !r.afterProxy {
		r.actions = append(r.actions, &headerAction{
			Key:    key,
			Value:  value,
			Action: headerActionSet,
		})
	}

}

func (r *ResponseHeader) AddHeader(key, value string) {
	r.cache.Add(key, value)
	if r.afterProxy {
		r.header.Add(key, value)
	} else {
		r.actions = append(r.actions, &headerAction{
			Key:    key,
			Value:  value,
			Action: headerActionAdd,
		})
	}
}

func (r *ResponseHeader) DelHeader(key string) {
	r.cache.Del(key)
	if r.afterProxy {
		r.header.Del(key)
	} else {
		r.actions = append(r.actions, &headerAction{
			Key:    key,
			Action: headerActionDel,
		})
	}
}
