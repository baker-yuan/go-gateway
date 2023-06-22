package drivers

import (
	"fmt"
	"reflect"

	"github.com/baker-yuan/go-gateway/eosc"
	"github.com/baker-yuan/go-gateway/utils/config"
	"github.com/baker-yuan/go-gateway/utils/schema"
)

// Factory 插件工厂IExtenderDriverFactory实现
type Factory[T any] struct {
	// 插件配置类型
	configType reflect.Type
	// schema.Generate
	render interface{}
	// 创建插件的函数
	createFunc func(id, name string, v *T, workers map[eosc.RequireId]eosc.IWorker) (eosc.IWorker, error)
	// 检查插件配置
	configCheckFunc func(v *T, workers map[eosc.RequireId]eosc.IWorker) error
}

// NewFactory 创建插件工厂
//
// @createFunc 创建插件函数
func NewFactory[T any](createFunc func(id, name string, v *T, workers map[eosc.RequireId]eosc.IWorker) (eosc.IWorker, error), configCheckFunc ...func(v *T, workers map[eosc.RequireId]eosc.IWorker) error) *Factory[T] {
	// 在这段代码中，configType 是一个反射 Type 类型的变量，它表示泛型参数 T 所对应的类型。
	// reflect.TypeOf 函数可以接受任意类型的值作为参数，并返回该值的 reflect.Type 类型的变量。在这里，我们使用 (*T)(nil) 来获取泛型类型 T 的反射 Type 类型的变量。由于 *T 是一个指针类型，因此我们首先使用 nil 来创建一个空指针，然后使用 (*T) 将其转换为 *T 类型的指针变量，最后使用 reflect.TypeOf 函数获取该指针变量的反射 Type 类型的变量。
	// 这样做的目的是为了获取泛型参数 T 所对应的类型，并将其存储在 configType 中。通过这个变量，我们可以在后续的代码中使用反射来访问和操作 T 类型的值。
	configType := reflect.TypeOf((*T)(nil))

	render, _ := schema.Generate(configType, nil)
	f := &Factory[T]{
		createFunc: createFunc,
		configType: configType,
		render:     render,
	}

	if len(configCheckFunc) == 1 {
		f.configCheckFunc = configCheckFunc[0]
	}
	return f
}

func (f *Factory[T]) Render() interface{} {
	return f.render
}

func (f *Factory[T]) Create(profession string, name string, label string, desc string, params map[string]interface{}) (eosc.IExtenderDriver, error) {
	if f.configCheckFunc == nil {
		return &Driver[T]{
			profession: profession,
			driver:     name,
			label:      label,
			desc:       desc,
			configType: f.configType,
			createFunc: f.createFunc,
		}, nil
	} else {
		return &DriverConfigChecker[T]{
			Driver: Driver[T]{
				profession: profession,
				driver:     name,
				label:      label,
				desc:       desc,
				configType: f.configType,
				createFunc: f.createFunc,
			},
			configCheckFunc: f.configCheckFunc,
		}, nil
	}
}

// Driver 插件
// https://help.apinto.com/docs/apinto/router/http.html#http-%E5%8D%8F%E8%AE%AE%E8%B7%AF%E7%94%B1
type Driver[T any] struct {
	profession string                                                                                     // 模块名
	driver     string                                                                                     // 驱动名
	label      string                                                                                     //
	desc       string                                                                                     //
	configType reflect.Type                                                                               //
	createFunc func(id, name string, v *T, workers map[eosc.RequireId]eosc.IWorker) (eosc.IWorker, error) // 创建插件函数
}

func (d *Driver[T]) Create(id, name string, v interface{}, workers map[eosc.RequireId]eosc.IWorker) (eosc.IWorker, error) {
	cfg, err := d.Assert(v)
	if err != nil {
		return nil, err
	}
	return d.createFunc(id, name, cfg, workers)
}

func (d *Driver[T]) ConfigType() reflect.Type {
	return d.configType
}

func (d *Driver[T]) Assert(v interface{}) (*T, error) {
	return Assert[T](v)
}

// DriverConfigChecker 插件配置检查
type DriverConfigChecker[T any] struct {
	Driver[T]                                                                 // 聚合了插件
	configCheckFunc func(v *T, workers map[eosc.RequireId]eosc.IWorker) error // 插件配置检查
}

func (d *DriverConfigChecker[T]) Check(v interface{}, workers map[eosc.RequireId]eosc.IWorker) error {
	cfg, err := d.Assert(v)
	if err != nil {
		return err
	}
	return d.configCheckFunc(cfg, workers)
}

// Assert 断言输入的参数v是否为类型T
func Assert[T any](v interface{}) (*T, error) {
	cfg, ok := v.(*T)
	if !ok {
		return nil, fmt.Errorf("%w:need %s,now %s", eosc.ErrorConfigType, config.TypeNameOf((*T)(nil)), config.TypeNameOf(v))
	}
	return cfg, nil
}

type WorkerBase struct {
	id   string
	name string
}

func Worker(id string, name string) WorkerBase {
	return WorkerBase{id: id, name: name}
}

func (w *WorkerBase) Id() string {
	return w.id
}
func (w *WorkerBase) Name() string {
	return w.name
}
