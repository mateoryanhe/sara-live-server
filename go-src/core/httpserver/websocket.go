package httpserver

import (
	"encoding/json"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/constants/common"
	"xr-game-server/core/cfg"
	"xr-game-server/core/event"
	"xr-game-server/core/xrtoken"
	"xr-game-server/errercode"
)

const (
	//发送超时
	SendTimeOut = 100 * time.Millisecond
	//写入超时
	IdleTime = 5000 * time.Millisecond
	ChkTime  = 2 * time.Nanosecond
	MaxSize  = 10
	BuffSize = 100
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
	data, _ := json.Marshal(ret)
	return data
}

// 生成连接失败错误日志
func newHeart() []byte {
	ret := []*PushResp{
		&PushResp{
			Cmd:  cmd.Heart,
			Data: time.Now().Unix(),
		},
	}
	data, _ := json.Marshal(ret)
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
		userId := r.GetQuery(AuthId).String()
		if len(userId) == common.Zero {
			ws.WriteMessage(websocket.BinaryMessage, newError(errercode.EmptyUserId))
			return
		}
		if !cfg.GetAuthCfg().LoginOff {
			token := r.GetQuery(Token).String()
			if len(token) == common.Zero {
				ws.WriteMessage(websocket.BinaryMessage, newError(errercode.EmptyToken))
				return
			}
			flag := xrtoken.HasAppToken(gconv.Uint64(userId), token)
			if !flag {
				ws.WriteMessage(websocket.BinaryMessage, newError(errercode.Token))
				return
			}
		}
		client := newClient(gconv.Uint64(userId), ws)
		client.init()
		//发布客户端进入成功事件
		event.Pub(event.ClientEnter, client)
	})
}

type WebSocketClient struct {
	//用户唯一标识
	Id   uint64
	Conn *websocket.Conn
	//数据发送时间
	sendTime time.Time
	//上次数据发送时间
	lastTime   time.Time
	dataBuffer chan any
	sendData   []any
	Loop       bool
	Num        uint
}

func newClient(id uint64, conn *websocket.Conn) *WebSocketClient {
	buffSize := BuffSize
	if cfg.WebSocketBufferCfgModel.Period > common.Zero {
		buffSize = cfg.WebSocketBufferCfgModel.Period
	}
	return &WebSocketClient{
		Id:         id,
		Conn:       conn,
		sendTime:   time.Now(),
		lastTime:   time.Now(),
		dataBuffer: make(chan any, buffSize),
		sendData:   make([]any, common.Zero),
		Loop:       true,
	}
}
func (w *WebSocketClient) init() {
	//go w.ChkClose()
	//go w.consume()
}

func (w *WebSocketClient) ChkClose() {
	_, _, err := w.Conn.ReadMessage()
	if err != nil {
		w.exit()
		return // Connection closed or error occurred
	}
}

func (w *WebSocketClient) exit() {
	w.Loop = false
	event.Pub(event.ClientLeave, w)
	err := w.Conn.Close()
	if err != nil {
		return
	}
}

// 发送缓存区数据
func (w *WebSocketClient) flushDataBuffer() {
	if len(w.sendData) == common.Zero {
		return
	}
	w.sendTime = time.Now()
	w.lastTime = time.Now()
	w.Num = common.Zero
	ret, _ := json.Marshal(w.sendData)
	//开始发送数据
	err := w.Conn.WriteMessage(websocket.BinaryMessage, ret)
	if err != nil {
		w.exit()
		return
	}
	w.sendData = make([]any, common.Zero)
}

func (w *WebSocketClient) Consume() {
	//如果到点了,开始发送数据
	if w.Num > common.Zero && time.Now().After(w.sendTime) {
		w.flushDataBuffer()
	}
	//好久没有发数据了,发一下数据,防止客户端重连
	if w.lastTime.Add(IdleTime).Before(time.Now()) {
		err := w.Conn.WriteMessage(websocket.BinaryMessage, newHeart())
		if err != nil {
			w.exit()
			return
		}
		w.lastTime = time.Now()
	}
	w.consumeData()
}

func (w *WebSocketClient) consumeData() {
	for {
		select {
		case data, ok := <-w.dataBuffer:
			{
				if ok {
					//检查数据是否过多
					if w.Num >= MaxSize {
						//开始发送数据
						w.flushDataBuffer()
						return
					} else {
						w.sendData = append(w.sendData, data)
						//检查超时时间是否设置
						if w.Num == common.Zero {
							w.sendTime = time.Now().Add(SendTimeOut)
						}
						w.Num++
					}
				}
			}
			break
		case <-time.After(ChkTime):
			{

			}
			return
		}
	}

}

func (c *WebSocketClient) Send(data any) {
	if !c.Loop {
		return
	}
	c.dataBuffer <- data
}
