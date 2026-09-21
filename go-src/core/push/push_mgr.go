package push

import (
	"sort"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/cmd"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrjson"
	"xr-game-server/core/xrlog"

	"xr-game-server/errercode"
)

var clientMap = gmap.New(true)

func Init() {
	event.Sub(event.ClientEnter, addClient)
	event.Sub(event.ClientLeave, rmClient)
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
	jStr, _ := xrjson.Marshal(msg)
	xrlog.DetailLog.Infof(gctx.New(), "authId=%v,发送数据cmd=%v,data=%s", clientId, cmd, jStr)
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
	xrlog.DetailLog.Infof(gctx.New(), "authId=%v,发送数据cmd=%v", clientId, cmd)
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
	leaveClient, ok := data.(*httpserver.WebSocketClient)
	if !ok || leaveClient == nil {
		return
	}
	xrlog.DetailLog.Infof(gctx.New(), "authId=%v,ip=%v,离线了", leaveClient.Id, clientRemoteAddr(leaveClient))
	event.Pub(event.Offline, event.NewOfflineData(leaveClient.Id))
	mapData := clientMap.Get(leaveClient.Id)
	if mapData == nil {
		return
	}
	mapClient, ok := mapData.(*httpserver.WebSocketClient)
	if !ok || mapClient == nil {
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
	client, ok := data.(*httpserver.WebSocketClient)
	if !ok || client == nil {
		return
	}
	xrlog.DetailLog.Infof(gctx.New(), "authId=%v,ip=%v,上线了", client.Id, clientRemoteAddr(client))
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

func clientRemoteAddr(client *httpserver.WebSocketClient) string {
	if client == nil || client.Conn == nil {
		return ""
	}
	return client.Conn.RemoteAddr().String()
}

// OnlineCount 当前在线连接数
func OnlineCount() int {
	if clientMap == nil {
		return 0
	}
	return clientMap.Size()
}

type onlineUserSnapshot struct {
	userId   uint64
	onlineAt int64
}

// OnlineUserIds 返回当前进程已建立 WebSocket 连接的用户ID快照，按上线时间倒序排列。
func OnlineUserIds() []uint64 {
	if clientMap == nil {
		return []uint64{}
	}
	keys := clientMap.Keys()
	snapshots := make([]onlineUserSnapshot, 0, len(keys))
	for _, key := range keys {
		userId, ok := key.(uint64)
		if !ok || userId == 0 {
			continue
		}
		data := clientMap.Get(userId)
		client, ok := data.(*httpserver.WebSocketClient)
		if !ok || client == nil {
			continue
		}
		snapshots = append(snapshots, onlineUserSnapshot{
			userId:   userId,
			onlineAt: client.OnlineAt,
		})
	}
	sort.Slice(snapshots, func(i, j int) bool {
		if snapshots[i].onlineAt == snapshots[j].onlineAt {
			return snapshots[i].userId > snapshots[j].userId
		}
		return snapshots[i].onlineAt > snapshots[j].onlineAt
	})
	userIds := make([]uint64, 0, len(snapshots))
	for _, snapshot := range snapshots {
		userIds = append(userIds, snapshot.userId)
	}
	return userIds
}

// IsOnline 根据当前 WebSocket 连接判断用户是否在线。
func IsOnline(clientId uint64) bool {
	if clientId == 0 || clientMap == nil {
		return false
	}
	return clientMap.Contains(clientId)
}
