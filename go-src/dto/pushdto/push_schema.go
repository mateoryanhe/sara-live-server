package pushdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dto/diamonddto"
	"xr-game-server/dto/golddto"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/dto/messagedto"
	"xr-game-server/dto/vipdto"
)

// EnterPushReq cmd=1 连接鉴权失败
type EnterPushReq struct {
	g.Meta `path:"/enter" method:"post" summary:"推送 cmd=1 Enter" description:"WebSocket 连接鉴权失败时返回" tags:"推送协议"`
}

type EnterPushResp struct {
	Cmd  int                  `json:"cmd" dc:"命令字 1"`
	Data *httpserver.AuthResp `json:"data"`
}

// HeartPushReq cmd=2 服务端心跳
type HeartPushReq struct {
	g.Meta `path:"/heart" method:"post" summary:"推送 cmd=2 Heart" description:"服务端心跳保活" tags:"推送协议"`
}

type HeartPushResp struct {
	Cmd  int   `json:"cmd" dc:"命令字 2"`
	Data int64 `json:"data" dc:"当前 Unix 时间戳(秒)"`
}

// KickPushReq cmd=3 踢下线
type KickPushReq struct {
	g.Meta `path:"/kick" method:"post" summary:"推送 cmd=3 Kick" description:"踢下线" tags:"推送协议"`
}

type KickPushResp struct {
	Cmd  int   `json:"cmd" dc:"命令字 3"`
	Data int64 `json:"data" dc:"踢出时间 Unix 毫秒时间戳"`
}

// CloseServerPushReq cmd=4 关服通知
type CloseServerPushReq struct {
	g.Meta `path:"/closeServer" method:"post" summary:"推送 cmd=4 CloseServer" description:"关服通知(仅 cmd，无 data)" tags:"推送协议"`
}

type CloseServerPushResp struct {
	Cmd int `json:"cmd" dc:"命令字 4(无 data)"`
}

// RepeatLoginPushReq cmd=5 重复登录
type RepeatLoginPushReq struct {
	g.Meta `path:"/repeatLogin" method:"post" summary:"推送 cmd=5 RepeatLogin" description:"重复登录，通知旧连接下线(仅 cmd，无 data)" tags:"推送协议"`
}

type RepeatLoginPushResp struct {
	Cmd int `json:"cmd" dc:"命令字 5(无 data)"`
}

// ErrorPushReq cmd=6 无参数错误
type ErrorPushReq struct {
	g.Meta `path:"/error" method:"post" summary:"推送 cmd=6 Error" description:"无参数错误" tags:"推送协议"`
}

type ErrorPushResp struct {
	Cmd  int `json:"cmd" dc:"命令字 6"`
	Data int `json:"data" dc:"错误码 XRCode"`
}

// ErrorParamPushReq cmd=7 带参数错误
type ErrorParamPushReq struct {
	g.Meta `path:"/errorParam" method:"post" summary:"推送 cmd=7 ErrorParam" description:"带参数错误" tags:"推送协议"`
}

type ErrorParamPushResp struct {
	Cmd  int            `json:"cmd" dc:"命令字 7"`
	Data *push.ErrorDto `json:"data"`
}

// LiveRoomGiftPushReq cmd=8 直播间送礼
type LiveRoomGiftPushReq struct {
	g.Meta `path:"/liveRoomGift" method:"post" summary:"推送 cmd=8 LiveRoomGift" description:"直播间送礼广播(房间内全体在线用户)" tags:"推送协议"`
}

type LiveRoomGiftPushResp struct {
	Cmd  int                       `json:"cmd" dc:"命令字 8"`
	Data *liveroomdto.GiftPushItem `json:"data"`
}

// LiveRoomChatPushReq cmd=9 直播间文字消息
type LiveRoomChatPushReq struct {
	g.Meta `path:"/liveRoomChat" method:"post" summary:"推送 cmd=9 LiveRoomChat" description:"直播间免费文字消息(房间内全体在线用户)" tags:"推送协议"`
}

type LiveRoomChatPushResp struct {
	Cmd  int                       `json:"cmd" dc:"命令字 9"`
	Data *liveroomdto.ChatPushItem `json:"data"`
}

// DiamondPushReq cmd=10 钻石余额
type DiamondPushReq struct {
	g.Meta `path:"/diamond" method:"post" summary:"推送 cmd=10 DiamondPush" description:"钻石余额变更(推送给指定用户)" tags:"推送协议"`
}

type DiamondPushResp struct {
	Cmd  int                         `json:"cmd" dc:"命令字 10"`
	Data *diamonddto.DiamondPushItem `json:"data"`
}

// GoldPushReq cmd=11 金币余额
type GoldPushReq struct {
	g.Meta `path:"/gold" method:"post" summary:"推送 cmd=11 GoldPush" description:"金币余额变更(推送给指定用户)" tags:"推送协议"`
}

type GoldPushResp struct {
	Cmd  int                   `json:"cmd" dc:"命令字 11"`
	Data *golddto.GoldPushItem `json:"data"`
}

// VipLevelPushReq cmd=12 VIP等级
type VipLevelPushReq struct {
	g.Meta `path:"/vipLevel" method:"post" summary:"推送 cmd=12 VipLevelPush" description:"VIP等级变更(推送给指定用户)" tags:"推送协议"`
}

type VipLevelPushResp struct {
	Cmd  int                      `json:"cmd" dc:"命令字 12"`
	Data *vipdto.VipLevelPushItem `json:"data"`
}

// PrivateMessagePushReq cmd=13 私信
type PrivateMessagePushReq struct {
	g.Meta `path:"/privateMessage" method:"post" summary:"推送 cmd=13 PrivateMessagePush" description:"私信消息(接收者与发送者均会收到)" tags:"推送协议"`
}

type PrivateMessagePushResp struct {
	Cmd  int                                `json:"cmd" dc:"命令字 13"`
	Data *messagedto.PrivateMessagePushItem `json:"data"`
}

// SystemMessagePushReq cmd=14 系统消息
type SystemMessagePushReq struct {
	g.Meta `path:"/systemMessage" method:"post" summary:"推送 cmd=14 SystemMessagePush" description:"系统消息(推送给接收者)" tags:"推送协议"`
}

type SystemMessagePushResp struct {
	Cmd  int                               `json:"cmd" dc:"命令字 14"`
	Data *messagedto.SystemMessagePushItem `json:"data"`
}

// LiveRoomAnchorBanPushReq cmd=15 主播封禁
type LiveRoomAnchorBanPushReq struct {
	g.Meta `path:"/liveRoomAnchorBan" method:"post" summary:"推送 cmd=15 LiveRoomAnchorBan" description:"主播封禁(推送给主播及直播间在线观众)" tags:"推送协议"`
}

type LiveRoomAnchorBanPushResp struct {
	Cmd  int                            `json:"cmd" dc:"命令字 15"`
	Data *liveroomdto.AnchorBanPushItem `json:"data"`
}

// LiveRoomAudienceMutePushReq cmd=17 观众禁言
type LiveRoomAudienceMutePushReq struct {
	g.Meta `path:"/liveRoomAudienceMute" method:"post" summary:"推送 cmd=17 LiveRoomAudienceMute" description:"观众禁言/解禁状态(推送给被操作用户)" tags:"推送协议"`
}

type LiveRoomAudienceMutePushResp struct {
	Cmd  int                               `json:"cmd" dc:"命令字 17"`
	Data *liveroomdto.AudienceMutePushItem `json:"data"`
}

// LiveRoomAudienceKickPushReq cmd=18 观众被踢
type LiveRoomAudienceKickPushReq struct {
	g.Meta `path:"/liveRoomAudienceKick" method:"post" summary:"推送 cmd=18 LiveRoomAudienceKick" description:"观众被踢出(推送给被踢用户)" tags:"推送协议"`
}

type LiveRoomAudienceKickPushResp struct {
	Cmd  int                               `json:"cmd" dc:"命令字 18"`
	Data *liveroomdto.AudienceKickPushItem `json:"data"`
}

// LiveRoomAudienceKickCancelPushReq cmd=19 取消进入限制
type LiveRoomAudienceKickCancelPushReq struct {
	g.Meta `path:"/liveRoomAudienceKickCancel" method:"post" summary:"推送 cmd=19 LiveRoomAudienceKickCancel" description:"取消观众进入限制(推送给被取消限制的用户)" tags:"推送协议"`
}

type LiveRoomAudienceKickCancelPushResp struct {
	Cmd  int                                     `json:"cmd" dc:"命令字 19"`
	Data *liveroomdto.AudienceKickCancelPushItem `json:"data"`
}

// LiveRoomAudienceJoinPushReq cmd=20 观众进房
type LiveRoomAudienceJoinPushReq struct {
	g.Meta `path:"/liveRoomAudienceJoin" method:"post" summary:"推送 cmd=20 LiveRoomAudienceJoin" description:"观众进入直播间(房间内全体在线用户)" tags:"推送协议"`
}

type LiveRoomAudienceJoinPushResp struct {
	Cmd  int                               `json:"cmd" dc:"命令字 20"`
	Data *liveroomdto.AudienceJoinPushItem `json:"data"`
}

// LiveRoomAudienceLeavePushReq cmd=21 观众离房
type LiveRoomAudienceLeavePushReq struct {
	g.Meta `path:"/liveRoomAudienceLeave" method:"post" summary:"推送 cmd=21 LiveRoomAudienceLeave" description:"观众离开直播间(房间内剩余在线用户)" tags:"推送协议"`
}

type LiveRoomAudienceLeavePushResp struct {
	Cmd  int                                `json:"cmd" dc:"命令字 21"`
	Data *liveroomdto.AudienceLeavePushItem `json:"data"`
}

// LiveRoomPaidDanmakuPushReq cmd=22 付费弹幕
type LiveRoomPaidDanmakuPushReq struct {
	g.Meta `path:"/liveRoomPaidDanmaku" method:"post" summary:"推送 cmd=22 LiveRoomPaidDanmaku" description:"付费弹幕(房间内全体在线用户)" tags:"推送协议"`
}

type LiveRoomPaidDanmakuPushResp struct {
	Cmd  int                              `json:"cmd" dc:"命令字 22"`
	Data *liveroomdto.PaidDanmakuPushItem `json:"data"`
}

// LiveRoomStopLivePushReq cmd=23 主播下播
type LiveRoomStopLivePushReq struct {
	g.Meta `path:"/liveRoomStopLive" method:"post" summary:"推送 cmd=23 LiveRoomStopLive" description:"主播下播(推送给直播间在线观众,不含主播)" tags:"推送协议"`
}

type LiveRoomStopLivePushResp struct {
	Cmd  int                                 `json:"cmd" dc:"命令字 23"`
	Data *liveroomdto.AnchorStopLivePushItem `json:"data"`
}
