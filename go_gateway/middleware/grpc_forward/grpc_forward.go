package grpc_forward

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/baker-yuan/go-gateway/common"
	v1 "github.com/baker-yuan/go-gateway/go-gateway-admin/api/admin/v1"
	"github.com/baker-yuan/go-gateway/registry"
	"github.com/baker-yuan/go-gateway/third_party/httprouter"
	"github.com/fullstorydev/grpcurl"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

var (
	client = defaultH2Client()
)

func GrpcForwardModeMiddleware(ctx *fasthttp.RequestCtx, next func(error)) {
	// 从上下文中获取规则
	rule := ctx.UserValue(common.CtxRuleKey).(*httprouter.ServiceInfo)
	if rule.InterfaceType != v1.InterfaceType_G_RPC {
		next(nil)
		return
	}

	// grpc接口，只支持http为post的请求
	if string(ctx.Method()) != http.MethodPost {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	// 获取下游服务
	discovery := registry.NewDiscovery()
	service, err := discovery.GetService(ctx, rule.Application)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	fmt.Println(service)

	// 获取grpc配置
	cfg := v1.GrpcConfig{}
	if err := json.Unmarshal([]byte(rule.Config), &cfg); err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	// 转发 优先级：json -> 反射 -> proto文件
	if cfg.GetCallType() == v1.GrpcConfig_JSON {
		resp, err := callGrpcByHttpCodec(ctx, "127.0.0.1", "9100", rule.InterfaceUrl)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(resp)
		return
	}

	if cfg.GetCallType() == v1.GrpcConfig_REFLECTION {
		resp, err := callGrpcByReflection(ctx, "127.0.0.1", "9100", rule.InterfaceUrl)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(resp)
		return
	}

	if cfg.GetCallType() == v1.GrpcConfig_PROTO {
		resp, err := callGrpcByProto(ctx, "127.0.0.1", "9100", rule.InterfaceUrl, cfg.GetProto())
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(resp)
		return
	}

	ctx.SetStatusCode(http.StatusInternalServerError)
}

func callGrpcByProto(ctx *fasthttp.RequestCtx, ip string, port string, method string, proto string) ([]byte, error) {
	var (
		clientConn *grpc.ClientConn
		refSource  grpcurl.DescriptorSource
		err        error
	)
	if clientConn, err = grpc.DialContext(
		ctx,
		fmt.Sprintf("%s:%s", ip, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(),
	); err != nil {
		panic(err)
	}

	var (
		descSourceFiles = map[string]string{}
		fileNames       = make([]string, 0)
		p               *protoparse.Parser
	)

	descSourceFiles = map[string]string{
		method: proto,
	}
	fileNames = append(fileNames, method)
	p = &protoparse.Parser{Accessor: protoparse.FileContentsFromMap(descSourceFiles)}
	descSources, err := p.ParseFiles(fileNames...)

	if err != nil {
		return nil, err
	}
	refSource, err = grpcurl.DescriptorSourceFromFileDescriptors(descSources...)
	if err != nil {
		return nil, err
	}

	return doCallGrpcByReflection(ctx, clientConn, method, refSource)
}

func callGrpcByReflection(ctx *fasthttp.RequestCtx, ip, port, method string) ([]byte, error) {
	var (
		clientConn *grpc.ClientConn
		refSource  grpcurl.DescriptorSource
		err        error
	)
	if clientConn, err = grpc.DialContext(
		ctx,
		fmt.Sprintf("%s:%s", ip, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(),
	); err != nil {
		panic(err)
	}
	defer clientConn.Close()
	refClient := grpcreflect.NewClient(context.Background(), reflectpb.NewServerReflectionClient(clientConn))
	defer refClient.Reset()

	refSource = grpcurl.DescriptorSourceFromServer(context.Background(), refClient)

	return doCallGrpcByReflection(ctx, clientConn, method, refSource)
}

var (
	options = grpcurl.FormatOptions{
		AllowUnknownFields: true,
	}
)

func doCallGrpcByReflection(ctx *fasthttp.RequestCtx, clientConn *grpc.ClientConn, method string, descSource grpcurl.DescriptorSource) ([]byte, error) {
	in := strings.NewReader(string(ctx.PostBody()))

	rf, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, descSource, in, options)
	if err != nil {
		return nil, err
	}
	var md []string

	response := NewResponse()
	handler := &grpcurl.DefaultEventHandler{
		VerbosityLevel: 2,
		Out:            response,
		Formatter:      formatter,
	}
	err = grpcurl.InvokeRPC(ctx, descSource, clientConn, "helloworld.Greeter/SayHello", md, handler, rf.Next)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func callGrpcByHttpCodec(ctx *fasthttp.RequestCtx, ip, port, method string) ([]byte, error) {
	log.Infof("callGrpcByHttpCodec ip: %s, port: %s, method: %s, body: %s", ip, port, method, string(ctx.Request.Body()))
	b := ctx.PostBody() // Content-Type 和 Content-Length 一定要有，不然读取不到数据
	// b := []byte(`{"name":"baker"}`)
	// chatGPT：
	// 这段代码是将请求体的长度编码成4个字节的大端序，并将编码后的长度与请求体一起拼接成一个新的字节切片。
	// 这样做是为了在请求数据传输过程中，让服务端能够准确地知道请求体的长度，从而正确地解析请求数据。
	// 在gRPC的底层协议中，请求的长度信息是必须的，因此这个操作是非常常见的。
	//
	// 这样设计是为了流式接口
	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests
	bb := make([]byte, len(b)+5)
	binary.BigEndian.PutUint32(bb[1:], uint32(len(b)))
	copy(bb[5:], b)

	// 请求参数
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s:%s%s", ip, port, method), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/grpc+json")
	req.Header.Del("Content-Length")
	req.ContentLength = int64(len(bb))
	req.Body = io.NopCloser(bytes.NewReader(bb))

	// 请求后端服务
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析响应
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		return data, nil
	}

	return data[5:], err
}

var _dialTimeout = 200 * time.Millisecond

func defaultH2Client() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			// So http2.Transport doesn't complain the URL scheme isn't 'https'
			AllowHTTP:          true,
			DisableCompression: true,
			// Pretend we are dialing a TLS endpoint.
			// Note, we ignore the passed tls.Config
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.DialTimeout(network, addr, _dialTimeout)
			},
		},
	}
}
