package eosc

import (
	"time"
)

type RequireId string

// IWorker 插件
type IWorker interface {
	Id() string                                                  // id唯一
	Start() error                                                // 启动
	Reset(conf interface{}, workers map[RequireId]IWorker) error // 重置
	Stop() error                                                 // 停止
	CheckSkill(skill string) bool                                //
}

// IWorkerDestroy 插件销毁
type IWorkerDestroy interface {
	Destroy() error // 销毁
}

// IWorkers 插件获取
type IWorkers interface {
	Get(id string) (IWorker, bool) // 获取插件
}

type TWorker struct {
	Id         string      `json:"id,omitempty" yaml:"id"`
	Name       string      `json:"name,omitempty" yaml:"name"`
	Driver     string      `json:"driver,omitempty" yaml:"driver"`
	Profession string      `json:"profession,omitempty" yaml:"profession"`
	Create     time.Time   `json:"create" yaml:"create"`
	Update     time.Time   `json:"update" yaml:"update"`
	Data       interface{} `json:"data,omitempty" yaml:"data"`
}
