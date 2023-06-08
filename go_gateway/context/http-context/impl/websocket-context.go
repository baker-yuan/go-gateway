package http_context

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	eoscContext "github.com/baker-yuan/go-gateway/context"
	http_context "github.com/baker-yuan/go-gateway/context/http-context"
	http_service "github.com/baker-yuan/go-gateway/context/http-context"
	"github.com/baker-yuan/go-gateway/log"
	"github.com/baker-yuan/go-gateway/utils/config"
	"github.com/fasthttp/websocket"
)

var _ http_context.IWebsocketContext = (*WebsocketContext)(nil)

type WebsocketContext struct {
	*HttpContext
	upstreamConn net.Conn
}

var upgrader = websocket.FastHTTPUpgrader{}

func (w *WebsocketContext) Upgrade() error {
	err := upgrader.Upgrade(w.fastHttpRequestCtx, func(conn *websocket.Conn) {
		if w.upstreamConn == nil {
			// 上游连接失败，直接返回
			log.Error("fail to connect upstream")
			return
		}
		defer conn.Close()
		defer w.upstreamConn.Close()
		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func() {
			size, err := io.Copy(conn.UnderlyingConn(), w.upstreamConn)
			log.Infof("finish copy upstream: size is %d,err is %v", size, err)
			wg.Done()
		}()
		go func() {
			size, err := io.Copy(w.upstreamConn, conn.UnderlyingConn())
			log.Infof("finish copy upstream: size is %d,err is %v", size, err)
			wg.Done()
		}()
		wg.Wait()
	})

	return err
}

func (w *WebsocketContext) IsWebsocket() bool {
	return websocket.FastHTTPIsWebSocketUpgrade(w.fastHttpRequestCtx)
}

func NewWebsocketContext(ctx http_context.IHttpContext) (*WebsocketContext, error) {
	httpCtx, ok := ctx.(*HttpContext)
	if !ok {
		return nil, errors.New("unsupported context type")
	}
	return &WebsocketContext{HttpContext: httpCtx}, nil
}

func (w *WebsocketContext) SetUpstreamConn(conn net.Conn) {
	w.upstreamConn = conn
}

func (w *WebsocketContext) Assert(i interface{}) error {
	if v, ok := i.(*http_context.IWebsocketContext); ok {
		*v = w
		return nil
	}
	if v, ok := i.(*http_service.IHttpContext); ok {
		*v = w
		return nil
	}
	return fmt.Errorf("not suport:%s", config.TypeNameOf(i))
}

func (w *WebsocketContext) IsCloneable() bool {
	return false
}

func (w *WebsocketContext) Clone() (eoscContext.EoContext, error) {
	return nil, fmt.Errorf("%s %w", "WebsocketContext", eoscContext.ErrEoCtxUnCloneable)
}
