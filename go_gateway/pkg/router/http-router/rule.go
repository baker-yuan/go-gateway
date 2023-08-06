package http_router

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/baker-yuan/go-gateway/pkg/checker"
	"github.com/baker-yuan/go-gateway/pkg/router"
)

var ErrorDuplicate = errors.New("duplicate")

// 路由匹配规则
// 路由可配置五种指标，分别是host、method、location、header、query。
//
// 指标匹配顺序
// 优先级递减 host -> method -> location -> header（key根据单词字母升序进行排序） -> query（key根据单词字母升序进行排序）

type Root struct {
	ports map[int]*Ports // key=端口 value=数据
}

type Ports struct {
	hosts map[string]*Hosts
}

type Hosts struct {
	methods map[string]*Methods
}

type Methods struct {
	paths map[string]*Paths
}

type Paths struct {
	handlers map[string]*Handler
	checker  checker.Checker
}

type IBuilder interface {
	Build() router.IMatcher
}

func (r *Root) Build() router.IMatcher {
	portsHandlers := make(map[string]router.IMatcher)
	for p, c := range r.ports {
		name := strconv.Itoa(p)
		if p == 0 {
			name = router.All
		}
		portsHandlers[name] = c.Build()
	}
	return newPortMatcher(portsHandlers)
}

func (p *Ports) Build() router.IMatcher {
	hostMatchers := make(map[string]router.IMatcher)
	for h, c := range p.hosts {
		hostMatchers[h] = c.Build()
	}
	return newHostMatcher(hostMatchers)
}

func (h *Hosts) Build() router.IMatcher {
	methodMatchers := make(map[string]router.IMatcher)
	for m, c := range h.methods {
		methodMatchers[m] = c.Build()
	}
	return newMethodMatcher(methodMatchers, nil)
}

func (m *Methods) Build() router.IMatcher {
	checkers := make([]*CheckerHandler, 0, len(m.paths))
	equals := make(map[string]router.IMatcher, len(m.paths))
	var all router.IMatcher
	for _, next := range m.paths {
		matcher := next.Build()
		if next.checker.CheckType() == checker.CheckTypeEqual {
			equals[next.checker.Value()] = matcher
			continue
		}
		if next.checker.CheckType() == checker.CheckTypeAll {
			all = next.Build()
			continue
		}

		checkers = append(checkers, &CheckerHandler{
			checker: next.checker,
			next:    matcher,
		})
	}
	return NewPathMatcher(equals, checkers, all)
}

func (p *Paths) Build() router.IMatcher {
	if len(p.handlers) == 0 {
		return &EmptyMatcher{handler: nil, has: false}
	}
	if all, has := p.handlers[router.All]; has {
		if len(p.handlers) == 1 {
			return &EmptyMatcher{handler: all.handler, has: true}
		}
	}

	nexts := make(AppendMatchers, 0, len(p.handlers))
	for _, h := range p.handlers {
		nexts = append(nexts, &AppendMatcher{
			handler:  h.handler,
			checkers: Parse(h.rules),
		})
	}
	sort.Sort(nexts)
	return nexts
}

type Handler struct {
	id      uint32
	handler router.IRouterHandler
	rules   []router.AppendRule
}

func (h *Handler) Build() router.IMatcher {
	return &AppendMatcher{
		handler:  h.handler,
		checkers: Parse(h.rules),
	}
}

func NewRoot() *Root {
	return &Root{
		ports: map[int]*Ports{},
	}
}

func NewPorts() *Ports {
	return &Ports{
		hosts: map[string]*Hosts{},
	}
}
func NewHosts() *Hosts {
	return &Hosts{
		methods: map[string]*Methods{},
	}
}
func NewMethods() *Methods {
	return &Methods{paths: map[string]*Paths{}}
}
func NewPaths(checker checker.Checker) *Paths {
	return &Paths{
		checker:  checker,
		handlers: map[string]*Handler{},
	}
}

func NewHandler(id uint32, handler router.IRouterHandler, appends []router.AppendRule) *Handler {
	return &Handler{id: id, handler: handler, rules: appends}
}

func (r *Root) Add(id uint32, handler router.IRouterHandler, port int, hosts []string, methods []string, path string, append []router.AppendRule) error {
	// 指标匹配顺序
	// 优先级递减 host -> method -> location -> header（key根据单词字母升序进行排序） -> query（key根据单词字母升序进行排序）
	if r.ports == nil {
		r.ports = make(map[int]*Ports)
	}
	pN, has := r.ports[port]
	if !has {
		pN = NewPorts()
		r.ports[port] = pN
	}
	// 调用Ports#Add
	err := pN.Add(id, handler, hosts, methods, path, append)
	if err != nil {
		return fmt.Errorf("port=%d %w", port, err)
	}
	return nil
}

func (p *Ports) Add(id uint32, handler router.IRouterHandler, hosts []string, methods []string, path string, append []router.AppendRule) error {
	// host为空
	if len(hosts) == 0 {
		return p.add(id, handler, router.All, methods, path, append)
	}
	// host不为空
	for _, host := range hosts {
		err := p.add(id, handler, host, methods, path, append)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Ports) add(id uint32, handler router.IRouterHandler, host string, methods []string, path string, append []router.AppendRule) error {
	// 没有就创建
	hN, has := p.hosts[host]
	if !has {
		hN = NewHosts()
		p.hosts[host] = hN
	}
	// Hosts#Add
	err := hN.Add(id, handler, methods, path, append)
	if err != nil {
		return fmt.Errorf("host=%s %w", host, err)
	}
	return nil
}

func (h *Hosts) add(id uint32, handler router.IRouterHandler, method string, path string, append []router.AppendRule) error {
	// 没有就创建
	methods, has := h.methods[method]
	if !has {
		methods = NewMethods()
		h.methods[method] = methods
	}
	// Methods#Add
	err := methods.Add(id, handler, path, append)
	if err != nil {
		return fmt.Errorf("method=%s %w", method, err)
	}
	return nil
}
func (h *Hosts) Add(id uint32, handler router.IRouterHandler, methods []string, path string, append []router.AppendRule) error {
	// 没有就创建
	if len(methods) == 0 {
		return h.add(id, handler, router.All, path, append)
	}
	// 添加
	for _, m := range methods {
		err := h.add(id, handler, m, path, append)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Methods) Add(id uint32, handler router.IRouterHandler, path string, append []router.AppendRule) error {
	// 没有就创建
	ck, err := checker.Parse(path)
	if err != nil {
		return fmt.Errorf("path=%s %w", path, err)
	}
	key := ck.Key()
	p, has := m.paths[key]
	if !has {
		p = NewPaths(ck)
		m.paths[key] = p
	}
	// Paths#Add
	err = p.Add(id, handler, append)
	if err != nil {
		return fmt.Errorf("path=%s %w", key, err)
	}
	return nil
}

func (p *Paths) Add(id uint32, handler router.IRouterHandler, append []router.AppendRule) error {
	key := router.Key(append)
	h, has := p.handlers[key]
	if has && h.id != id {
		return fmt.Errorf(" append{%s}:%w for (%s %s) ", key, ErrorDuplicate, h.id, id)
	}
	p.handlers[key] = NewHandler(id, handler, append)
	return nil
}
