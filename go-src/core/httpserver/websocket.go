package httpserver

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/cfg"
	"xr-game-server/core/event"
	"xr-game-server/core/xrjson"
	"xr-game-server/core/xrtoken"
	"xr-game-server/errercode"
)

const (
	// SendTimeOut 批量发送等待窗口,窗口内消息合并一次发送
	SendTimeOut = 10 * time.Millisecond
	// IdleTime 空闲超过该时间才发送心跳
	IdleTime        = 5000 * time.Millisecond
	defaultMaxBatch = 10
	wsPushSlowLogMs = int64(5)
)

type PushResp struct {
	Cmd  int `json:"cmd" dc:"命令"`
	Data any `json:"data" dc:"数据"`
}

type AuthResp struct {
	Code int `json:"code" dc:"错误码"`
}

// 生成连接失败错误日志
func newError(code errercode.XRCode) []byte {
	ret := []*PushResp{
		&PushResp{
			Cmd:  cmd.Enter,
			Data: &AuthResp{Code: int(code)},
		},
	}
	data := xrjson.MustMarshal(ret)
	return data
}

// 生成连接失败错误日志
func newHeart() []byte {
	ret := []*PushResp{
		&PushResp{
			Cmd:  cmd.Heart,
			Data: time.Now().UnixMilli(),
		},
	}
	data := xrjson.MustMarshal(ret)
	return data
}

func InitWebsocket() {
	var (
		wsUpGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// In production, you should implement proper origin checking
				return true
			},
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				// Error callback function.
			},
		}
	)

	httpServer.BindHandler(Ws, func(r *ghttp.Request) {
		ws, err := wsUpGrader.Upgrade(r.Response.Writer, r.Request, nil)
		if err != nil {
			return
		}
		//开始检查token
		authStr := r.GetQuery("Authorization", "").String()
		if authStr == "" || len(strings.Split(authStr, ".")) != 2 {
			ws.WriteMessage(websocket.BinaryMessage, newError(errercode.Token))
			return
		}
		userId := strings.Split(authStr, ".")[0]
		token := strings.Split(authStr, ".")[1]
		flag := xrtoken.HasAppToken(gconv.Uint64(userId), token)
		if !flag {
			ws.WriteMessage(websocket.BinaryMessage, newError(errercode.Token))
			return
		}
		client := newClient(gconv.Uint64(userId), ws)
		client.init()
		//发布客户端进入成功事件
		event.Pub(event.ClientEnter, client)
	})
}

type WebSocketClient struct {
	//用户唯一标识
	Id         uint64
	Conn       *websocket.Conn
	dataBuffer *gqueue.TQueue[any]
	Loop       bool
	cancel     context.CancelFunc
}

func newClient(id uint64, conn *websocket.Conn) *WebSocketClient {
	return &WebSocketClient{
		Id:         id,
		Conn:       conn,
		dataBuffer: gqueue.NewTQueue[any](),
		Loop:       true,
	}
}

func (w *WebSocketClient) init() {
	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel
	go w.readPump()
	go w.consumeData(ctx)
}

// readPump 独立读协程,连接断开时 ReadMessage 返回错误并触发 exit
func (w *WebSocketClient) readPump() {
	defer w.exit()
	for w.Loop && w.Conn != nil {
		if _, _, err := w.Conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (w *WebSocketClient) exit() {
	if !w.Loop {
		return
	}
	w.Loop = false
	w.clearPendingData()
	if w.cancel != nil {
		w.cancel()
	}

	event.Pub(event.ClientLeave, w)
	if w.Conn != nil {
		_ = w.Conn.Close()
		w.Conn = nil
	}
}

func (w *WebSocketClient) clearPendingData() {
	for {
		select {
		case <-w.dataBuffer.C:
		default:
			return
		}
	}
}

func getMaxBatchSize() int {
	if cfg.WebSocketBufferCfgModel.Size > 0 {
		return cfg.WebSocketBufferCfgModel.Size
	}
	return defaultMaxBatch
}

func (w *WebSocketClient) flushBatch(batch []any, idleTimer *time.Timer) {
	if !w.Loop || len(batch) == 0 || w.Conn == nil {
		return
	}
	err := w.Conn.WriteMessage(websocket.BinaryMessage, xrjson.MustMarshal(batch))
	if err != nil {
		w.exit()
		return
	}
	resetIdleTimer(idleTimer)
}

func (w *WebSocketClient) sendHeart(idleTimer *time.Timer) {
	if !w.Loop || w.Conn == nil {
		return
	}
	if err := w.Conn.WriteMessage(websocket.BinaryMessage, newHeart()); err != nil {
		w.exit()
		return
	}
	resetIdleTimer(idleTimer)
}

func resetIdleTimer(idleTimer *time.Timer) {
	stopTimer(idleTimer)
	idleTimer.Reset(IdleTime)
}

func stopTimer(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
}

func (w *WebSocketClient) consumeData(ctx context.Context) {
	maxBatch := getMaxBatchSize()
	batch := make([]any, 0, maxBatch)
	flushTimer := time.NewTimer(time.Hour)
	stopTimer(flushTimer)
	defer flushTimer.Stop()

	idleTimer := time.NewTimer(IdleTime)
	defer idleTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			w.flushBatch(batch, idleTimer)
			return
		case <-idleTimer.C:
			if w.Loop {
				w.sendHeart(idleTimer)
			}
		case data, ok := <-w.dataBuffer.C:
			if !ok {
				w.flushBatch(batch, idleTimer)
				return
			}
			if !w.Loop {
				w.flushBatch(batch, idleTimer)
				return
			}
			batch = append(batch, data)
			if len(batch) == 1 {
				flushTimer.Reset(SendTimeOut)
			}
			if len(batch) >= maxBatch {
				stopTimer(flushTimer)
				w.flushBatch(batch, idleTimer)
				batch = batch[:0]
			}
		case <-flushTimer.C:
			w.flushBatch(batch, idleTimer)
			batch = batch[:0]
		}
	}
}

func (c *WebSocketClient) Send(data any) {
	if !c.Loop {
		return
	}
	start := time.Now()
	c.dataBuffer.Push(data)
	if costMs := time.Since(start).Milliseconds(); costMs >= wsPushSlowLogMs {
		g.Log().Warningf(gctx.New(), "websocket Push慢,userId=%v,耗时=%vms,data=%v", c.Id, costMs, data)
	}
}
