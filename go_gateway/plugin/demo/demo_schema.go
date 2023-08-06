package demo

import (
	pkg_plugin "github.com/baker-yuan/go-gateway/pkg/plugin"
)

const (
	Priority   = 90100
	PluginName = "demo"
)

type DemoPluginSchema struct {
}

func PluginSchema() *pkg_plugin.PluginSchema {
	return &pkg_plugin.PluginSchema{
		Name:       PluginName,
		JsonSchema: `[{"label":"key","valueName":"key","dataType":"String","elementType":"input","defaultValue":"","required":false,"disabled":false,"check":false,"maxlength":"","innerType":"","placeholder":"","key":"1594104154244","level":1}]`,
		Creator:    NewPlugin,
		Priority:   Priority,
	}
}

func NewPlugin(cv []byte) (pkg_plugin.IPluginInstance, error) {
	return nil, nil
}
