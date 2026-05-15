package push

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gutil"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/constants/common"
	"xr-game-server/core/cfg"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrtimer"
	"xr-game-server/errercode"
)

var clientMap = gmap.New(true)

func Init() {
	event.Sub(event.ClientEnter, addClient)
	event.Sub(event.ClientLeave, rmClient)
	period := 10
	if cfg.WebSocketBufferCfgModel.Period > common.Zero {
		period = cfg.WebSocketBufferCfgModel.Period
	}
	xrtimer.ModuleLoop(time.Duration(period)*time.Millisecond, consumeBuff)
	xrtimer.ModuleLoop(time.Duration(period)*time.Millisecond, chkLife)
}

func consumeBuff() {
	clients := clientMap.Values()
	if len(clients) > common.Zero {
		for _, tempClient := range clients {
			client := tempClient.(*httpserver.WebSocketClient)
			client.Consume()
		}
	}
}

func chkLife() {
	clients := clientMap.Values()
	if len(clients) > common.Zero {
		for _, tempClient := range clients {
			gutil.TryCatch(gctx.New(), func(try context.Context) {
				client := tempClient.(*httpserver.WebSocketClient)
				client.ChkClose()
			}, func(catch context.Context, err error) {

			})
		}
	}
}

// Data 推送消息给指定玩家
func Data(clientId uint64, cmd int, msg any) {
	tempClient := clientMap.Get(clientId)
	if tempClient == nil {
		return
	}
	client := tempClient.(*httpserver.WebSocketClient)
	if client == nil {
		return
	}
	jStr, _ := json.MarshalIndent(msg, "", " ")
	g.Log().Infof(gctx.New(), "userid=%v,发送数据cmd=%v,data=%s", clientId, cmd, jStr)
	client.Send(&httpserver.PushResp{
		Data: msg,
		Cmd:  cmd,
	})
}

// OutData 推送消息
func OutData(clientId uint64, cmd int) {
	tempClient := clientMap.Get(clientId)
	if tempClient == nil {
		return
	}
	client := tempClient.(*httpserver.WebSocketClient)
	if client == nil {
		return
	}
	g.Log().Infof(gctx.New(), "userid=%v,发送数据cmd=%v", clientId, cmd)
	client.Send(&httpserver.PushResp{
		Cmd: cmd,
	})
}

// WithIds 推送消息
func WithIds(clientIds []uint64, cmd int, msg any) {
	for _, clientId := range clientIds {
		Data(clientId, cmd, msg)
	}
}

// Broadcast 广播
func Broadcast(cmd int, msg any) {
	for _, val := range clientMap.Keys() {
		Data(val.(uint64), cmd, msg)
	}
}

// BroadcastNonData 只广播命令字
func BroadcastNonData(cmd int) {
	for _, val := range clientMap.Keys() {
		OutData(val.(uint64), cmd)
	}
}

// WithIdsNonData 推送消息
func WithIdsNonData(clientIds []uint64, cmd int) {
	for _, clientId := range clientIds {
		OutData(clientId, cmd)
	}
}

func Kick(clientId uint64) {
	data := clientMap.Get(clientId)
	if data == nil {
		return
	}
	client := data.(*httpserver.WebSocketClient)
	client.Send(&httpserver.PushResp{
		Data: time.Now().UnixMilli(),
		Cmd:  cmd.Kick,
	})
}

func rmClient(data any) {
	if data == nil {
		return
	}
	leaveClient := data.(*httpserver.WebSocketClient)
	if leaveClient == nil {
		return
	}
	g.Log().Infof(gctx.New(), "ip=%v,玩家=%v,离线了", leaveClient.Conn.RemoteAddr(), leaveClient.Id)
	//
	event.Pub(event.Offline, event.NewOfflineData(leaveClient.Id))
	mapClient := (clientMap.Get(leaveClient.Id)).(*httpserver.WebSocketClient)
	if mapClient == nil {
		return
	}
	if leaveClient.Conn != mapClient.Conn {
		return
	}
	clientMap.Remove(leaveClient.Id)
}

func addClient(data any) {
	if data == nil {
		return
	}
	client := data.(*httpserver.WebSocketClient)
	g.Log().Infof(gctx.New(), "ip=%v,玩家=%v,上线了", client.Conn.RemoteAddr(), client.Id)
	//尝试下线当前客户端,由客户端主动登出
	OutData(client.Id, cmd.RepeatLogin)
	clientMap.Set(client.Id, client)
	//发布上线事件
	event.Pub(event.Online, event.NewOnlineData(client.Id))
}

func Error(clientId uint64, err errercode.XRCode) {
	data := clientMap.Get(clientId)
	if data == nil {
		return
	}
	client := data.(*httpserver.WebSocketClient)
	client.Send(&httpserver.PushResp{
		Data: int(err),
		Cmd:  cmd.Error,
	})
}

func ErrorWithParam(clientId uint64, err errercode.XRCode, param any) {
	data := clientMap.Get(clientId)
	if data == nil {
		return
	}
	client := data.(*httpserver.WebSocketClient)
	client.Send(&httpserver.PushResp{
		Data: NewErrorDto(err, param),
		Cmd:  cmd.ErrorParam,
	})
}
