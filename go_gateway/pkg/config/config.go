package config

import (
	"flag"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// -conf=/Users/yuanyu/code/go-study/gateway/baker-gateway/go_gateway/config/gateway.yaml
var (
	configPath = flag.String("conf", "./conf/gateway.yaml", "gateway config path")
)

type GatewayConfig struct {
	Proxy ProxyConfig `json:"proxy" yaml:"proxy"`
	Sync  SyncConfig  `json:"sync" yaml:"sync"`
}

// ProxyConfig 转发服务配置
type ProxyConfig struct {
	Http struct {
		Port uint32 `json:"port" yaml:"port"`
	} `json:"http" yaml:"http"`
}

// SyncConfig 数据同步配置
type SyncConfig struct {
	Etcd struct {
		Endpoints   []string `json:"endpoints" yaml:"endpoints"`       // 集群地址
		DialTimeout uint32   `json:"dial_timeout" yaml:"dial_timeout"` // 连接超时
	} `json:"etcd" yaml:"etcd"`
}

// Load 加载配置文件
func Load() (*GatewayConfig, error) {
	// 参数解析
	flag.Parse()

	// 读取配置
	buf, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, err
	}
	cfg := defaultConfig()
	if err := yaml.Unmarshal(buf, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func defaultConfig() *GatewayConfig {
	return &GatewayConfig{}
}
