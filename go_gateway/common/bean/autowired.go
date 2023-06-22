package bean

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/baker-yuan/go-gateway/log"
	// "github.com/baker-yuan/go-gateway/log"
)

// Container bean 容器接口
type Container interface {
	Autowired(p interface{})                             // 声明
	Injection(i interface{})                             // 注入
	InjectionDefault(i interface{})                      // 注入默认值， 这里注入的只会对没有被其他注入对接口生效
	Check() error                                        // 检查是否实现相关dao类
	AddInitializingBean(handler InitializingBeanHandler) // 注入完成后的回调
	AddInitializingBeanFunc(handler func())              //
}

// InitializingBeanHandler 注入完成后的回调
type InitializingBeanHandler interface {
	AfterPropertiesSet()
}

type initializingBeanFunc func()

func (fun initializingBeanFunc) AfterPropertiesSet() {
	fun()
}

type container struct {
	autowiredInterfaces map[string][]reflect.Value // key=接口体唯一标识 value=依赖key结构体的结构体
	cache               map[string]reflect.Value   // key=接口体唯一标识 value=结构体
	defaultInterface    map[string]reflect.Value   // key=接口体唯一标识 value=默认值
	initializingBean    []InitializingBeanHandler  //
	lock                sync.Mutex                 //
	once                sync.Once                  //
}

func (m *container) AddInitializingBean(handler InitializingBeanHandler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if need := m.check(); len(need) == 0 {
		handler.AfterPropertiesSet()
		return
	} else {
		log.Debug("AddInitializingBean need:", strings.Join(need, ","))
	}
	m.initializingBean = append(m.initializingBean, handler)

}

func (m *container) AddInitializingBeanFunc(handler func()) {
	m.AddInitializingBean(initializingBeanFunc(handler))
}

// NewContainer 创建新的 bean 容器
func NewContainer() Container {
	return &container{
		autowiredInterfaces: make(map[string][]reflect.Value),
		cache:               make(map[string]reflect.Value),
		defaultInterface:    make(map[string]reflect.Value),
		lock:                sync.Mutex{},
		once:                sync.Once{},
	}
}

// 我依赖别人，我要去缓存找依赖，注入到我的变量中
func (m *container) add(key string, v reflect.Value) {

	// 缓存里面找（Injection放入的）
	if e, has := m.cache[key]; has {
		log.DebugF("autowired set:%s,%v", key, e)
		v.Set(e)
		return
	}

	// 有默认值设置默认值
	if ed, has := m.defaultInterface[key]; has {
		log.DebugF("autowired set default:%s,%v", key, ed)
		v.Set(ed)
	}

	log.DebugF("autowired cache :%s,%v", key, v)

	// 缓存和默认都没找到，我暂停依赖注入，我需要依赖别人依赖注入完成，在初始化我
	m.autowiredInterfaces[key] = append(m.autowiredInterfaces[key], v)
}

func (m *container) set(key string, v reflect.Value) {
	// 如果有人依赖我完成初始化，这里把依赖我的结构体进行依赖注入
	values := m.autowiredInterfaces[key]
	delete(m.autowiredInterfaces, key) // 依赖注入完了，直接删除

	for _, e := range values {
		e.Set(v)
	}
}
func (m *container) check() []string {

	r := make([]string, 0, len(m.autowiredInterfaces))
	for pkg := range m.autowiredInterfaces {
		if _, has := m.defaultInterface[pkg]; !has {
			r = append(r, pkg)
		}
	}

	return r
}

// Autowired 声明
func (m *container) Autowired(p interface{}) {
	pkg, e := pkg(p)
	m.lock.Lock()
	log.Debug("autowired: ", pkg, " point: ", p)

	m.add(pkg, e)
	m.lock.Unlock()
}

// Injection 注入
func (m *container) Injection(i interface{}) {
	// 1、获取结构体唯一标识
	// 2、把结构体加入缓存
	// 3、有人依赖我完成注入，把我注入给依赖我的这些结构体。没人依赖我，啥都不做（加了缓存）。
	pkg, v := pkg(i)

	m.lock.Lock()
	log.Debug("injection: ", pkg, " point: ", i)

	m.cache[pkg] = v
	m.set(pkg, v)
	m.lock.Unlock()
}

// InjectionDefault 注入默认值， 这里注入的只会对没有被其他注入对接口生效
func (m *container) InjectionDefault(i interface{}) {

	pkg, v := pkg(i)
	m.lock.Lock()
	m.defaultInterface[pkg] = v
	m.lock.Unlock()

}

func (m *container) injectionAll() {

	cache := m.cache
	for k, v := range cache {
		m.set(k, v)
	}

	defaults := m.defaultInterface
	for k, v := range defaults {
		m.set(k, v)
	}

}

// 获取一个接口体的唯一标识
func pkg(i interface{}) (string, reflect.Value) {
	v := reflect.ValueOf(i)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	pkg := key(v.Type())
	if pkg == "" {
		panic("invalid interface")
	}
	return pkg, v
}
func key(t reflect.Type) string {
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.String())
}

// Check 检查是否实现相关dao类
func (m *container) Check() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	var err error = nil
	m.once.Do(func() {

		m.injectionAll()
		rs := m.check()
		if len(rs) > 0 {
			err = fmt.Errorf("need:%v", rs)
			return
		}

		m.dispatchAfterPropertiesSet()
	})

	return err
}

func (m *container) dispatchAfterPropertiesSet() {

	beanHandlers := m.initializingBean
	for _, h := range beanHandlers {
		h.AfterPropertiesSet()
	}
}
