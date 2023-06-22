package plugin

import (
	"reflect"

	eocontext "github.com/baker-yuan/go-gateway/context"
)

// Config 普通插件配置，在router、service、upstream的插件格式
type Config struct {
	Disable bool        `json:"disable"` // 是否禁用
	Config  interface{} `json:"config"`  // 配置
}

// IPluginManager 插件管理器
type IPluginManager interface {
	CreateRequest(id string, conf map[string]*Config) eocontext.IChainPro // 获取拦击器链
	Global() eocontext.IChainPro                                          // 获取全局插件
	GetConfigType(name string) (reflect.Type, bool)                       // 配置
}

func MergeConfig(high, low map[string]*Config) map[string]*Config {
	if high == nil && low == nil {
		return make(map[string]*Config)
	}
	if high == nil {
		return clone(low)
	}
	if low == nil {
		return clone(high)
	}

	mv := clone(low)

	for k, hv := range high {
		lv, has := mv[k]
		if has {
			*lv = *hv
		} else {
			c := new(Config)
			*c = *hv
			mv[k] = c
		}
	}
	return mv
}

func clone(v map[string]*Config) map[string]*Config {
	cv := make(map[string]*Config)
	if v == nil {
		return cv
	}
	for k, v := range v {
		c := new(Config)
		*c = *v
		cv[k] = c
	}
	return cv
}
