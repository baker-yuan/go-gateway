package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/baker-yuan/go-gateway/common/bean"
	eoscContext "github.com/baker-yuan/go-gateway/context"
	http_service "github.com/baker-yuan/go-gateway/context/http-context"
	http_context "github.com/baker-yuan/go-gateway/context/http-context/impl"
	plugin_manager "github.com/baker-yuan/go-gateway/drivers/plugin-manager"
	ip_restriction "github.com/baker-yuan/go-gateway/drivers/plugins/ip-restriction"
	"github.com/baker-yuan/go-gateway/drivers/service"
	"github.com/baker-yuan/go-gateway/eosc"
	"github.com/baker-yuan/go-gateway/plugin"
	"github.com/baker-yuan/go-gateway/process-worker/workers"
	"github.com/baker-yuan/go-gateway/professions"
	"github.com/baker-yuan/go-gateway/variable"
	"github.com/valyala/fasthttp"
	// _ "github.com/baker-yuan/go-gateway/drivers/plugin-manager"
)

// var professionJsonConfig = `
// [
//     {
//         "name": "router",
//         "label": "路由",
//         "desc": "路由",
//         "dependencies": [
//             "service",
//             "template"
//         ],
//         "appendLabels": [
//             "host",
//             "service",
//             "listen",
//             "disable"
//         ],
//         "drivers": [
//             {
//                 "id": "eolinker.com:apinto:http_router",
//                 "name": "http",
//                 "label": "http",
//                 "desc": "http路由"
//             },
//             {
//                 "id": "eolinker.com:apinto:grpc_router",
//                 "name": "grpc",
//                 "label": "grpc",
//                 "desc": "grpc路由"
//             },
//             {
//                 "id": "eolinker.com:apinto:dubbo2_router",
//                 "name": "dubbo2",
//                 "label": "dubbo2",
//                 "desc": "dubbo2路由"
//             }
//         ]
//     },
//     {
//         "name": "service",
//         "label": "服务",
//         "desc": "服务",
//         "dependencies": [
//             "discovery"
//         ],
//         "appendLabels": [
//             "discovery"
//         ],
//         "drivers": [
//             {
//                 "id": "eolinker.com:apinto:service_http",
//                 "name": "http",
//                 "label": "service",
//                 "desc": "服务"
//             }
//         ]
//     },
//     {
//         "name": "strategy",
//         "label": "策略",
//         "desc": "策略",
//         "drivers": [
//             {
//                 "id": "eolinker.com:apinto:strategy-limiting",
//                 "name": "limiting",
//                 "label": "限流策略",
//                 "desc": "限流策略"
//             },
//             {
//                 "id": "eolinker.com:apinto:strategy-cache",
//                 "name": "cache",
//                 "label": "缓存策略",
//                 "desc": "缓存策略"
//             },
//             {
//                 "id": "eolinker.com:apinto:strategy-grey",
//                 "name": "grey",
//                 "label": "灰度策略",
//                 "desc": "灰度策略"
//             },
//             {
//                 "id": "eolinker.com:apinto:strategy-visit",
//                 "name": "visit",
//                 "label": "访问策略",
//                 "desc": "访问策略"
//             },
//             {
//                 "id": "eolinker.com:apinto:strategy-fuse",
//                 "name": "fuse",
//                 "label": "熔断策略",
//                 "desc": "熔断策略"
//             }
//         ]
//     },
//     {
//         "name": "transcode",
//         "label": "编码器",
//         "desc": "编码器",
//         "drivers": [
//             {
//                 "id": "eolinker.com:apinto:protobuf_transcode",
//                 "name": "protobuf",
//                 "label": "protobuf编码器",
//                 "desc": "protobuf编码器"
//             }
//         ]
//     },
//     {
//         "name": "app",
//         "label": "应用",
//         "desc": "应用",
//         "appendLabels": [
//             "disable"
//         ],
//         "drivers": [
//             {
//                 "id": "eolinker.com:apinto:app",
//                 "name": "app",
//                 "label": "应用",
//                 "desc": "应用"
//             }
//         ]
//     },
//     {
//         "name": "certificate",
//         "label": "证书",
//         "desc": "证书",
//         "drivers": [
//             {
//                 "id": "eolinker.com:apinto:ssl-server",
//                 "name": "server",
//                 "label": "证书",
//                 "desc": "证书"
//             }
//         ]
//     },
//     {
//         "name": "discovery",
//         "label": "注册中心",
//         "desc": "注册中心",
//         "drivers": [
//             {
//                 "id": "eolinker.com:apinto:discovery_static",
//                 "name": "static",
//                 "label": "静态服务发现",
//                 "desc": "静态服务发现"
//             },
//             {
//                 "id": "eolinker.com:apinto:discovery_nacos",
//                 "name": "nacos",
//                 "label": "nacos服务发现",
//                 "desc": "nacos服务发现"
//             },
//             {
//                 "id": "eolinker.com:apinto:discovery_consul",
//                 "name": "consul",
//                 "label": "consul服务发现",
//                 "desc": "consul服务发现"
//             },
//             {
//                 "id": "eolinker.com:apinto:discovery_eureka",
//                 "name": "eureka",
//                 "label": "eureka服务发现",
//                 "desc": "consul服务发现"
//             }
//         ]
//     },
//     {
//         "name": "output",
//         "label": "输出",
//         "desc": "输出",
//         "drivers": [
//             {
//                 "id": "eolinker.com:apinto:file_output",
//                 "name": "file",
//                 "label": "文件输出",
//                 "desc": "文件输出"
//             },
//             {
//                 "id": "eolinker.com:apinto:nsqd",
//                 "name": "nsqd",
//                 "label": "NSQ输出",
//                 "desc": "NSQ输出"
//             },
//             {
//                 "id": "eolinker.com:apinto:http_output",
//                 "name": "http_output",
//                 "label": "http输出",
//                 "desc": "http输出"
//             },
//             {
//                 "id": "eolinker.com:apinto:syslog_output",
//                 "name": "syslog_output",
//                 "label": "syslog输出",
//                 "desc": "syslog输出"
//             },
//             {
//                 "id": "eolinker.com:apinto:kafka_output",
//                 "name": "kafka_output",
//                 "label": "kafka输出",
//                 "desc": "kafka输出"
//             },
//             {
//                 "id": "eolinker.com:apinto:redis",
//                 "name": "redis",
//                 "label": "redis 集群",
//                 "desc": "redis 集群"
//             },
//             {
//                 "id": "eolinker.com:apinto:influxdbv2",
//                 "name": "influxdbv2",
//                 "label": "influxdbv2输出",
//                 "desc": "influxdbv2输出"
//             },
//             {
//                 "id": "eolinker.com:apinto:prometheus_output",
//                 "name": "prometheus",
//                 "label": "prometheus输出",
//                 "desc": "prometheus输出"
//             }
//         ]
//     },
//     {
//         "name": "template",
//         "label": "模版",
//         "desc": "模版",
//         "drivers": [
//             {
//                 "id": "eolinker.com:apinto:plugin_template",
//                 "name": "plugin_template",
//                 "label": "插件模版",
//                 "desc": "插件模版"
//             }
//         ]
//     }
// ]
// `

var professionJsonConfig = `
[
    {
        "name": "router",
        "label": "路由",
        "desc": "路由",
        "dependencies": [
            "service",
            "template"
        ],
        "appendLabels": [
            "host",
            "service",
            "listen",
            "disable"
        ],
        "drivers": [
            {
                "id": "eolinker.com:apinto:http_router",
                "name": "http",
                "label": "http",
                "desc": "http路由"
            },
            {
                "id": "eolinker.com:apinto:grpc_router",
                "name": "grpc",
                "label": "grpc",
                "desc": "grpc路由"
            },
            {
                "id": "eolinker.com:apinto:dubbo2_router",
                "name": "dubbo2",
                "label": "dubbo2",
                "desc": "dubbo2路由"
            }
        ]
    },
    {
        "name": "service",
        "label": "服务",
        "desc": "服务",
        "dependencies": [
            "discovery"
        ],
        "appendLabels": [
            "discovery"
        ],
        "drivers": [
            {
                "id": "service_http",
                "name": "http",
                "label": "service",
                "desc": "服务"
            }
        ]
    },
    {
        "name": "strategy",
        "label": "策略",
        "desc": "策略",
        "drivers": [
            {
                "id": "eolinker.com:apinto:strategy-limiting",
                "name": "limiting",
                "label": "限流策略",
                "desc": "限流策略"
            },
            {
                "id": "eolinker.com:apinto:strategy-cache",
                "name": "cache",
                "label": "缓存策略",
                "desc": "缓存策略"
            },
            {
                "id": "eolinker.com:apinto:strategy-grey",
                "name": "grey",
                "label": "灰度策略",
                "desc": "灰度策略"
            },
            {
                "id": "eolinker.com:apinto:strategy-visit",
                "name": "visit",
                "label": "访问策略",
                "desc": "访问策略"
            },
            {
                "id": "eolinker.com:apinto:strategy-fuse",
                "name": "fuse",
                "label": "熔断策略",
                "desc": "熔断策略"
            }
        ]
    },
    {
        "name": "transcode",
        "label": "编码器",
        "desc": "编码器",
        "drivers": [
            {
                "id": "eolinker.com:apinto:protobuf_transcode",
                "name": "protobuf",
                "label": "protobuf编码器",
                "desc": "protobuf编码器"
            }
        ]
    },
    {
        "name": "app",
        "label": "应用",
        "desc": "应用",
        "appendLabels": [
            "disable"
        ],
        "drivers": [
            {
                "id": "eolinker.com:apinto:app",
                "name": "app",
                "label": "应用",
                "desc": "应用"
            }
        ]
    },
    {
        "name": "certificate",
        "label": "证书",
        "desc": "证书",
        "drivers": [
            {
                "id": "eolinker.com:apinto:ssl-server",
                "name": "server",
                "label": "证书",
                "desc": "证书"
            }
        ]
    },
    {
        "name": "discovery",
        "label": "注册中心",
        "desc": "注册中心",
        "drivers": [
            {
                "id": "eolinker.com:apinto:discovery_static",
                "name": "static",
                "label": "静态服务发现",
                "desc": "静态服务发现"
            },
            {
                "id": "eolinker.com:apinto:discovery_nacos",
                "name": "nacos",
                "label": "nacos服务发现",
                "desc": "nacos服务发现"
            },
            {
                "id": "eolinker.com:apinto:discovery_consul",
                "name": "consul",
                "label": "consul服务发现",
                "desc": "consul服务发现"
            },
            {
                "id": "eolinker.com:apinto:discovery_eureka",
                "name": "eureka",
                "label": "eureka服务发现",
                "desc": "consul服务发现"
            }
        ]
    },
    {
        "name": "output",
        "label": "输出",
        "desc": "输出",
        "drivers": [
            {
                "id": "eolinker.com:apinto:file_output",
                "name": "file",
                "label": "文件输出",
                "desc": "文件输出"
            },
            {
                "id": "eolinker.com:apinto:nsqd",
                "name": "nsqd",
                "label": "NSQ输出",
                "desc": "NSQ输出"
            },
            {
                "id": "eolinker.com:apinto:http_output",
                "name": "http_output",
                "label": "http输出",
                "desc": "http输出"
            },
            {
                "id": "eolinker.com:apinto:syslog_output",
                "name": "syslog_output",
                "label": "syslog输出",
                "desc": "syslog输出"
            },
            {
                "id": "eolinker.com:apinto:kafka_output",
                "name": "kafka_output",
                "label": "kafka输出",
                "desc": "kafka输出"
            },
            {
                "id": "eolinker.com:apinto:redis",
                "name": "redis",
                "label": "redis 集群",
                "desc": "redis 集群"
            },
            {
                "id": "eolinker.com:apinto:influxdbv2",
                "name": "influxdbv2",
                "label": "influxdbv2输出",
                "desc": "influxdbv2输出"
            },
            {
                "id": "eolinker.com:apinto:prometheus_output",
                "name": "prometheus",
                "label": "prometheus输出",
                "desc": "prometheus输出"
            }
        ]
    },
    {
        "name": "template",
        "label": "模版",
        "desc": "模版",
        "drivers": [
            {
                "id": "eolinker.com:apinto:plugin_template",
                "name": "plugin_template",
                "label": "插件模版",
                "desc": "插件模版"
            }
        ]
    }
]
`

// https://github.com/baker-yuan/go-gateway/pull/113
// https://blog.csdn.net/weixin_41479678/article/details/111933900
func main() {
	// 创建ExtenderRegister 插件工厂
	register := eosc.NewExtenderRegister()

	// 加载内置插件 插件放入工厂
	ip_restriction.Register(register)
	service.Register(register)

	// 往容器注入eosc.IExtenderDrivers 通过插件名称获取对应的插件工厂
	var extenderDrivers eosc.IExtenderDrivers = register
	bean.Injection(&extenderDrivers)
	fmt.Println("[info] 注入插件工厂 IExtenderDrivers main.go")

	// IWorkers
	profession := professions.NewProfessions(register)

	var configs []*eosc.ProfessionConfig
	_ = json.Unmarshal([]byte(professionJsonConfig), &configs)
	profession.Reset(configs)
	var workers *workers.Workers = workers.NewWorkerManager(profession)
	var iw eosc.IWorkers = workers
	bean.Injection(&iw) // 插件获取
	fmt.Println("[info] 注入Workers main.go")

	plugin_manager.Init()
	plugin_manager.Register(register)

	var pluginManager plugin.IPluginManager
	bean.Autowired(&pluginManager)
	pg := pluginManager.(*plugin_manager.PluginManager)

	// 设置启用插件
	cfg := &plugin_manager.PluginWorkerConfig{
		Plugins: []*plugin_manager.PluginConfig{
			{
				Name:   "my_ip_restriction",
				ID:     "ip_restriction",
				Status: "enable",
			},
		},
	}
	err := pg.Set(cfg)
	if err != nil {
		panic(err)
	}

	//
	variable := variable.NewVariables(nil)
	body := []byte(`
	{
		"name": "rate_limiting_service",
		"driver": "http",
		"timeout": 3000,
		"retry": 3,
		"description": "使用黑白ip插件",
		"scheme": "http",
		"nodes": [
			"demo.apinto.com:8280"
		],
		"balance": "round-robin",
		"plugins": {
			"my_ip_restriction": {
				"disable": false,
				"config": {
					"ip_list_type": "black",
					"ip_black_list": [
						"127.0.0.1"
					]
				}
			}
		}
	}
	`)
	_ = workers.Set("rate_limiting_service@service", "service", "rate_limiting_service", "http", body, variable)

	conf := map[string]*plugin.Config{
		"my_ip_restriction": {
			Disable: false,
			Config: &ip_restriction.Config{
				IPListType: "white",
				IPWhiteList: []string{
					"127.0.0.1",
				},
			},
		},

		// "my_ip_restriction": {
		// 	Disable: false,
		// 	Config: &ip_restriction.Config{
		// 		IPListType: "black",
		// 		IPBlackList: []string{
		// 			"*",
		// 		},
		// 	},
		// },
	}
	conf = conf
	var filters eoscContext.IChainPro = pluginManager.CreateRequest("ip_restriction_router@router", conf)
	fmt.Println(filters)

	port := 8848

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		panic(err)
	}
	notFound := &HttpNotFoundHandler{}
	server := fasthttp.Server{
		StreamRequestBody:            true,
		DisablePreParseMultipartForm: true,
		MaxRequestBodySize:           100 * 1024 * 1024,
		ReadBufferSize:               16 * 1024,
		Handler: func(ctx *fasthttp.RequestCtx) {
			fmt.Println("handle...start")
			var httpContext eoscContext.EoContext = http_context.NewContext(ctx, port)

			err := filters.Chain(httpContext)

			if err != nil {
				_ = notFound.Finish(httpContext)
				return
			}

			_ = notFound.Complete(httpContext)

			// httpContext.SetCompleteHandler(notFound)
			// httpContext.SetFinish(notFound)

			fmt.Println("handle...end")
		}}
	_ = server.Serve(ln)
}

type HttpNotFoundHandler struct {
}

func (m *HttpNotFoundHandler) Complete(ctx eoscContext.EoContext) error {
	httpContext, err := http_service.Assert(ctx)
	if err != nil {
		return nil
	}
	httpContext.Response().SetHeader("Content-Type", "application/json")
	httpContext.Response().SetStatus(200, "200")
	httpContext.Response().SetBody([]byte(`{code:"200",message:"ok"}`))
	return nil
}

func (m *HttpNotFoundHandler) Finish(ctx eoscContext.EoContext) error {
	httpContext, err := http_service.Assert(ctx)
	if err != nil {
		return err
	}
	httpContext.FastFinish()
	return nil
}
