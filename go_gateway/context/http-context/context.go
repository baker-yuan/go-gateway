package http_context

import (
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"time"

	gcontext "github.com/baker-yuan/go-gateway/context"
)

type keyCloneCtx struct{}
type keyHttpRetry struct{}
type keyHttpTimeout struct{}

var (
	KeyHttpRetry   = keyHttpRetry{}
	KeyHttpTimeout = keyHttpTimeout{}
	KeyCloneCtx    = keyCloneCtx{}
)

type FileHeader struct {
	FileName string
	Header   textproto.MIMEHeader
	Data     []byte
}

// IWebsocketContext websocket 继承了IHttpContext接口
type IWebsocketContext interface {
	IHttpContext
	Upgrade() error
	SetUpstreamConn(conn net.Conn)
	IsWebsocket() bool
}

// IHttpContext http 扩展EoContext接口
type IHttpContext interface {
	gcontext.EoContext                                                      // 组合EoContext
	Request() IRequestReader                                                // 读取原始请求
	Proxy() IRequest                                                        // 组装转发的request
	Response() IResponse                                                    // 处理返回结果，可读可写
	SendTo(scheme string, node gcontext.INode, timeout time.Duration) error // 如果下游是http服务，通过这个方法转发到下游
	Proxies() []IProxy                                                      //
	FastFinish()                                                            // 结束请求，释放资源
}

// IHeaderReader 请求头读取
type IHeaderReader interface {
	RawHeader() string            // 返回原始的http请求头数据
	GetHeader(name string) string // 根据key获取请求头
	Headers() http.Header         // 返回所有的请求头
	Host() string                 // 获取HTTP请求头中的 "Host" 字段。"Host" 字段是一个HTTP/1.1请求中必须存在的字段，它指定了请求的目标主机名和（可选的）端口号。
	GetCookie(key string) string  // 获取HTTP请求头中的特定cookie值
}

// IHeaderWriter 设置响应头
type IHeaderWriter interface {
	IHeaderReader                // 组合了IHeaderReader
	SetHeader(key, value string) //
	AddHeader(key, value string) //
	DelHeader(key string)        //
	SetHost(host string)         //
}

// IResponseHeader 设置响应头
type IResponseHeader interface {
	GetHeader(name string) string
	Headers() http.Header
	HeadersString() string
	SetHeader(key, value string)
	AddHeader(key, value string)
	DelHeader(key string)
}

// IBodyGet 请求体获取
type IBodyGet interface {
	GetBody() []byte
	BodyLen() int
}

// IBodySet 设置请求体
type IBodySet interface {
	SetBody([]byte)
}

type IBodyDataReader interface {
	// ContentType "Content-Type" 字段用于表示发送给接收者的实体主体的媒体类型。
	ContentType() string
	// BodyForm
	// content-Type = application/x-www-form-urlencoded 或 multipart/form-data，与原生request.Form不同，这里不包括 query 参数
	BodyForm() (url.Values, error)
	// Files
	// content-Type = multipart/form-data 时有效
	Files() (map[string][]*multipart.FileHeader, error)

	GetForm(key string) string

	GetFile(key string) (file []*multipart.FileHeader, has bool)

	RawBody() ([]byte, error)
}

type IBodyDataWriter interface {
	IBodyDataReader

	// SetForm
	// 设置form数据并将content-type设置 为 application/x-www-form-urlencoded 或 multipart/form-data
	SetForm(values url.Values) error
	SetToForm(key, value string) error
	AddForm(key, value string) error

	// AddFile
	// 会替换掉对应掉file信息，并且将content-type 设置为 multipart/form-data
	AddFile(key string, file *multipart.FileHeader) error

	// SetRaw
	// 设置 multipartForm 数据并将content-type设置 为 multipart/form-data
	// 重置body，会清除掉未处理掉 form和file
	SetRaw(contentType string, body []byte)
}

// IStatusGet 获取http状态吗
type IStatusGet interface {
	StatusCode() int      // 获取响应状态吗
	Status() string       // 获取字符串格式的响应状态吗
	ProxyStatusCode() int //
	ProxyStatus() string  //
}

// IStatusSet 设置http状态吗
type IStatusSet interface {
	SetStatus(code int, status string)      // 设置http状态吗
	SetProxyStatus(code int, status string) //
}

// IQueryReader 获取url查询参数
type IQueryReader interface {
	GetQuery(key string) string // 根据key获取查询参数
	RawQuery() string           // 获取所有的查询参数，url ? 后面的
}

// IQueryWriter 设置url查询参数
type IQueryWriter interface {
	IQueryReader                // 组合IQueryReader
	SetQuery(key, value string) //
	AddQuery(key, value string) //
	DelQuery(key string)        //
	SetRawQuery(raw string)     //
}

type IURIReader interface {
	IQueryReader        //
	RequestURI() string //
	Scheme() string     //
	RawURL() string     //
	Host() string       //
	Path() string       //
}

type IURIWriter interface {
	IURIReader               //
	IQueryWriter             //
	SetPath(string)          //
	SetScheme(scheme string) //
	SetHost(host string)     //
}

// IRequestReader 原始请求数据的读取接口
type IRequestReader interface {
	URI() IURIReader       // url读取
	Header() IHeaderReader // 请求头读取
	Body() IBodyDataReader // 请求体读取
	RemoteAddr() string    // 客户端地址
	RemotePort() string    //
	RealIp() string        // 客户端ip
	ForwardIP() string     //
	Method() string        // 请求方法
	String() string        //
	ContentLength() int    //
	ContentType() string   // 请求类型
}

// IRequest 用于组装转发的request
type IRequest interface {
	Header() IHeaderWriter   //
	Body() IBodyDataWriter   //
	URI() IURIWriter         //
	Method() string          //
	ContentLength() int      //
	ContentType() string     //
	SetMethod(method string) //
}

// IProxy 记录转发相关信息
type IProxy interface {
	IRequest              // 组装转发的request
	StatusCode() int      //
	Status() string       //
	ProxyTime() time.Time //
	ResponseLength() int  //
	ResponseTime() int64  //
}

// IResponse 返回给client端的
type IResponse interface {
	IStatusGet                              //
	IResponseHeader                         //
	IStatusSet                              // 设置返回状态
	IBodySet                                // 设置返回内容
	IBodyGet                                //
	ResponseError() error                   //
	ClearError()                            //
	String() string                         //
	SetResponseTime(duration time.Duration) //
	ResponseTime() time.Duration            //
	ContentLength() int                     //
	ContentType() string                    //
}
