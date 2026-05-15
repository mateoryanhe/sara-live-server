package cfg

import (
	"encoding/json"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	WebSocketBufferStr = "bufferSize.websocket"
)

type WebSocketBufferCfg struct {
	Size   int
	Period int
}

var WebSocketBufferCfgModel = &WebSocketBufferCfg{}

func initWebSocketBufferCfg() {
	data, _ := g.Cfg().GetWithCmd(gctx.New(), WebSocketBufferStr)
	err := data.Scan(&WebSocketBufferCfgModel)
	if err != nil {
		g.Log().Error(gctx.New(), "无法加载到websocket缓冲池配置数据")
	} else {
		cfgJson, _ := json.MarshalIndent(WebSocketBufferCfgModel, "", " ")
		g.Log().Warningf(gctx.New(), "成功加载到websocket缓冲池配置数据:%s", cfgJson)
	}
}
