# 基于go实现分布式网关

# grpc泛化调用方式

## grpc服务直接支持json编码

> 1、文章：http://www.baker-yuan.cn/articles/317
>
> 2、案例：https://github.com/go-kratos/gateway/blob/main/middleware/transcoder/transcoder.go

## grpc服务开启反射

> 1、文章：http://www.baker-yuan.cn/articles/294
>
> 2、案例：https://github.com/eolinker/apinto/blob/v0.12.5/drivers/plugins/http-to-gRPC/complete.go


## 通过proto文件调用

> 1、文章：http://www.baker-yuan.cn/articles/294
>
> 2、案例：https://github.com/eolinker/apinto/blob/v0.12.5/drivers/plugins/http-to-gRPC/complete.go

# 三方工具

## grpcurl

> grpcurl是一个命令行工具，使用它可以在命令行中访问gRPC服务，就像使用curl访问http服务一样。
> https://github.com/fullstorydev/grpcurl
>
> 文章：
>
> 1、http://www.baker-yuan.cn/articles/294
>
> 2、http://www.baker-yuan.cn/articles/292

## fasthttp

fasthttp是一个高性能的HTTP服务器框架，它是Go语言中最快的HTTP服务器之一。它的设计目标是实现高吞吐量和低延迟，以满足高并发的需求。

fasthttp相较于标准库中的http包，有以下优势：

1. 更快的性能：fasthttp使用了更加高效的内存管理和I/O复用机制，能够实现更高的吞吐量和更低的延迟。
2. 更少的内存占用：fasthttp的内存管理非常精细，能够避免不必要的内存分配和释放，减少内存占用。
3. 更易于扩展：fasthttp提供了丰富的中间件和插件机制，可以方便地进行定制和扩展。
4. 更多的功能：fasthttp支持HTTP/1.1和HTTP/2协议，同时还支持WebSocket、TLS和gzip等功能。
   fasthttp的使用非常简单，只需要实现一个请求处理函数并指定监听地址即可，如下所示：

https://pkg.go.dev/github.com/valyala/fasthttp



## httprouter

> https://github.com/julienschmidt/httprouter
>
> chatGPT：
> httprouter是一个轻量级的Go语言HTTP请求路由器，它能够高效地处理HTTP请求，并支持RESTful API的设计。它的主要特点包括：
> 1. 高性能：httprouter使用了trie树的算法来实现路由匹配，比常规的正则表达式匹配更快。
> 2. 简单易用：httprouter的API非常简单，只需要调用一个函数并传入路由和处理函数即可。
> 3. 支持RESTful API：httprouter支持HTTP请求的GET、POST、PUT、DELETE等RESTful API，可以快速地设计和开发RESTful API。
> 4. 支持中间件：httprouter支持中间件，可以在路由处理函数之前或之后执行一些自定义的逻辑。
> 5. 可扩展性强：httprouter的代码非常简洁，易于扩展和修改，可以根据实际需求进行定制化开发。
>
> 总之，httprouter是一个简单、高效、易用、可扩展的Go语言HTTP请求路由器，非常适合用于构建高性能的Web应用程序和RESTful API。

> 文章：http://www.baker-yuan.cn/articles/322



## kratos

> Kratos 一套轻量级 Go 微服务框架，包含大量微服务相关功能及工具。
>
> https://github.com/go-kratos/kratos



## protoreflect

> 用于在Go语言中反射和操作Protocol Buffers（protobuf）消息的库。它允许您动态地读取、写入和修改protobuf消息，而无需生成代码或使用特定的结构体。该库还支持在运行时动态生成protobuf消息，并支持使用JSON和文本格式的消息序列化和反序列化。
>
> https://github.com/jhump/protoreflect

> 文章：http://www.baker-yuan.cn/articles/295



## Polaris

> 北极星是一个支持多语言和多框架的服务发现和治理平台，致力于解决分布式和微服务架构中的服务管理、流量管理、故障容错、配置管理和可观测性问题，针对不同的技术栈和环境提供服务治理的标准方案和最佳实践。
>
> https://github.com/polarismesh/polaris



# 开源网关推荐

## B站开源的简单网关

> A high-performance API Gateway with middlewares, supporting HTTP and gRPC protocols.
>
> https://github.com/go-kratos/gateway

这个项目处于起步阶段，不是很成熟，无法直接用于生产环境。

## Apache ShenYu

> 这是一个异步的，高性能的，跨语言的，响应式的 `API` 网关。
>
> https://github.com/apache/shenyu



## Apinto

> Apinto 是专门为微服务架构设计的开源 API 网关，完全由 Go 语言开发，拥有目前市面上最强的性能及稳定性表现，并且可以自由扩展几乎所有功能模块。 提供丰富的流量管理、数据处理、协议转换等功能，例如动态路由、负载均衡、服务发现、熔断降级、身份认证、监控与告警等。
>
> https://github.com/eolinker/apinto



> 悟空API网关升级版本
>
> - https://github.com/eolinker/goku_lite
> - https://community.apinto.com/d/34051-apinto-gateway
>
> - https://community.apinto.com/d/34044-gokuapinto

EOLINKER 旗下的微服务网关，都有对应的商业版本。开源版本感觉抠抠搜搜的。

# 项目截图

## 接口列表
![Image text](./doc/img/httr_rule_list.jpg)

## 新增gRPC接口



## 心中Http接口