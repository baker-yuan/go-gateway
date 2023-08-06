package gateway

import (
	"context"
	"time"

	pb "github.com/baker-yuan/go-gateway/pb/router"
	"github.com/baker-yuan/go-gateway/pkg/util"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	gatewayHttpRouter = "/gateway/httpRouter/" // 网关插件
)

func (e *Engine) Sync() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		watcher clientv3.Watcher
		err     error
	)

	cfg := e.config.Sync
	// 初始化配置
	config = clientv3.Config{
		Endpoints:   cfg.Etcd.Endpoints,
		DialTimeout: time.Duration(cfg.Etcd.DialTimeout) * time.Millisecond,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}
	watcher = clientv3.NewWatcher(client)
	kv = clientv3.NewKV(client)

	go func() {
		e.syncHttpRouter(kv, watcher)
	}()

}

func (e *Engine) syncHttpRouter(kv clientv3.KV, watcher clientv3.Watcher) {
	ctx := context.Background()
	// 获取全量数据
	value, err := kv.Get(ctx, gatewayHttpRouter)
	if err != nil {
		return
	}
	for _, kvpair := range value.Kvs {
		router, err := util.Unmarshal[pb.HttpRouter](kvpair.Value)
		if err != nil {
			continue
		}
		e.httpRouteManager.Set(&router, e.serviceManager)
	}

	// 监听增量数据
	watchChan := watcher.Watch(ctx, gatewayHttpRouter)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT: // 数据保存修改
				router, err := util.Unmarshal[pb.HttpRouter](event.Kv.Value)
				if err != nil {
					continue
				}
				e.httpRouteManager.Set(&router, e.serviceManager)
			case mvccpb.DELETE: // 数据删除
				e.httpRouteManager.Delete(util.StrToUint32Def0(string(event.Kv.Value)))
			}
		}
	}

}
