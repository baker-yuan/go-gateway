package main

import (
	gateway "github.com/baker-yuan/go-gateway"
)

func main() {

	// 创建引擎
	engine, err := gateway.New()
	if err != nil {
		panic(err)
	}

	// 网关数据同步
	engine.Sync()

	// 启动网关
	if err = engine.Start(); err != nil {
		panic(err)
	}
}
