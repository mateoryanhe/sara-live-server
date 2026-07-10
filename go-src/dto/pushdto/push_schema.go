package pushdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dto/calldto"
	"xr-game-server/dto/diamonddto"
	"xr-game-server/dto/golddto"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/dto/messagedto"
	"xr-game-server/dto/vipdto"
)

// Swagger 推送协议 tags(按业务模块划分)
const (
	PushTagConnection     = "推送-连接"
	PushTagError          = "推送-错误"
	PushTagWallet         = "推送-钱包"
	PushTagPrivateMessage = "推送-私信"
	PushTagMessage        = "推送-消息"
	PushTagLiveRoom       = "推送-直播间"
	PushTagCall           = "推送-通话"
	PushTagPrivateRoom    = "推送-私密房"
)

// --- 推送-通话 ---

// LiveRoomCallAnchorAcceptedAudiencePushReq cmd=33 主播开始接听视频通话(推送给直播间观众)
type LiveRoomCallAnchorAcceptedAudiencePushReq struct {
	g.Meta `path:"/liveRoomCallAnchorAcceptedAudience" method:"post" summary:"推送 cmd=33 主播开始接听视频通话(推送给直播间观众)" description:"主播开始接听视频通话(推送给直播间在线观众,不含主播与呼叫者)" tags:"推送-通话"`
}

type LiveRoomCallAnchorAcceptedAudiencePushResp struct {
	Cmd  int                                         `json:"cmd" dc:"命令字 33"`
	Data *calldto.CallAnchorAcceptedAudiencePushItem `json:"data"`
}

// LiveRoomCallTimeoutPushReq cmd=32 直播间通话呼叫超时
type LiveRoomCallTimeoutPushReq struct {
	g.Meta `path:"/liveRoomCallTimeout" method:"post" summary:"推送 cmd=32 直播间通话呼叫超时(推送给呼叫者与接听者)" description:"直播间通话呼叫超时(推送给呼叫者与接听者)" tags:"推送-通话"`
}

type LiveRoomCallTimeoutPushResp struct {
	Cmd  int                          `json:"cmd" dc:"命令字 32"`
	Data *calldto.CallTimeoutPushItem `json:"data"`
}

// LiveRoomCallStartedPushReq cmd=31 直播间通话开始
type LiveRoomCallStartedPushReq struct {
	g.Meta `path:"/liveRoomCallStarted" method:"post" summary:"推送 cmd=31 直播间通话开始(推送给呼叫者与接听者)" description:"直播间通话开始(推送给呼叫者与接听者)" tags:"推送-通话"`
}

type LiveRoomCallStartedPushResp struct {
	Cmd  int                          `json:"cmd" dc:"命令字 31"`
	Data *calldto.CallStartedPushItem `json:"data"`
}

// LiveRoomCallEndedPushReq cmd=28 直播间通话结束
type LiveRoomCallEndedPushReq struct {
	g.Meta `path:"/liveRoomCallEnded" method:"post" summary:"推送 cmd=28 直播间通话结束(推送给对方)" description:"直播间通话结束(推送给对方)" tags:"推送-通话"`
}

type LiveRoomCallEndedPushResp struct {
	Cmd  int                        `json:"cmd" dc:"命令字 28"`
	Data *calldto.CallEndedPushItem `json:"data"`
}

// LiveRoomCallAcceptedPushReq cmd=27 直播间通话被接听
type LiveRoomCallAcceptedPushReq struct {
	g.Meta `path:"/liveRoomCallAccepted" method:"post" summary:"推送 cmd=27 直播间通话被接听(推送给呼叫者)" description:"直播间通话被接听(推送给呼叫者)" tags:"推送-通话"`
}

type LiveRoomCallAcceptedPushResp struct {
	Cmd  int                           `json:"cmd" dc:"命令字 27"`
	Data *calldto.CallAcceptedPushItem `json:"data"`
}

// LiveRoomCallRejectedPushReq cmd=26 直播间通话被拒接
type LiveRoomCallRejectedPushReq struct {
	g.Meta `path:"/liveRoomCallRejected" method:"post" summary:"推送 cmd=26 直播间通话被拒接(推送给呼叫者)" description:"直播间通话被拒接(推送给呼叫者)" tags:"推送-通话"`
}

type LiveRoomCallRejectedPushResp struct {
	Cmd  int                           `json:"cmd" dc:"命令字 26"`
	Data *calldto.CallRejectedPushItem `json:"data"`
}

// LiveRoomCallRequestPushReq cmd=25 直播间通话请求
type LiveRoomCallRequestPushReq struct {
	g.Meta `path:"/liveRoomCallRequest" method:"post" summary:"推送 cmd=25 直播间通话请求(推送给主播)" description:"直播间通话请求(推送给主播)" tags:"推送-通话"`
}

type LiveRoomCallRequestPushResp struct {
	Cmd  int                          `json:"cmd" dc:"命令字 25"`
	Data *calldto.CallRequestPushItem `json:"data"`
}

// --- 推送-私密房 ---

// LiveRoomPrivateGiftPushReq cmd=30 给指定主播送礼
type LiveRoomPrivateGiftPushReq struct {
	g.Meta `path:"/liveRoomPrivateGift" method:"post" summary:"推送 cmd=30 给指定主播送礼(推送给发送者与主播)" description:"给指定主播送礼(推送给发送者与主播)" tags:"推送-私密房"`
}

type LiveRoomPrivateGiftPushResp struct {
	Cmd  int                              `json:"cmd" dc:"命令字 30"`
	Data *liveroomdto.PrivateGiftPushItem `json:"data"`
}

// LiveRoomPrivateChatPushReq cmd=29 私密房文字消息
type LiveRoomPrivateChatPushReq struct {
	g.Meta `path:"/liveRoomPrivateChat" method:"post" summary:"推送 cmd=29 私密房文字消息(推送给发送者与目标用户)" description:"私密房文字消息(推送给发送者与目标用户)" tags:"推送-私密房"`
}

type LiveRoomPrivateChatPushResp struct {
	Cmd  int                                  `json:"cmd" dc:"命令字 29"`
	Data *liveroomdto.PrivateRoomChatPushItem `json:"data"`
}

// --- 推送-直播间 ---

// LiveRoomStartLivePushReq cmd=34 主播开播
type LiveRoomStartLivePushReq struct {
	g.Meta `path:"/liveRoomStartLive" method:"post" summary:"推送 cmd=34 主播开播(全服广播)" description:"主播开播(全服广播)" tags:"推送-直播间"`
}

type LiveRoomStartLivePushResp struct {
	Cmd  int                                  `json:"cmd" dc:"命令字 34"`
	Data *liveroomdto.AnchorStartLivePushItem `json:"data"`
}

// LiveRoomAudienceListRefreshPushReq cmd=24 观众列表刷新
type LiveRoomAudienceListRefreshPushReq struct {
	g.Meta `path:"/liveRoomAudienceListRefresh" method:"post" summary:"推送 cmd=24 观众列表刷新(房间内全体在线用户,含主播)" description:"观众列表刷新(房间内全体在线用户,含主播)" tags:"推送-直播间"`
}

type LiveRoomAudienceListRefreshPushResp struct {
	Cmd  int                                      `json:"cmd" dc:"命令字 24"`
	Data *liveroomdto.AudienceListRefreshPushItem `json:"data"`
}

// LiveRoomStopLivePushReq cmd=23 主播下播
type LiveRoomStopLivePushReq struct {
	g.Meta `path:"/liveRoomStopLive" method:"post" summary:"推送 cmd=23 主播下播(全服广播)" description:"主播下播(全服广播)" tags:"推送-直播间"`
}

type LiveRoomStopLivePushResp struct {
	Cmd  int                                 `json:"cmd" dc:"命令字 23"`
	Data *liveroomdto.AnchorStopLivePushItem `json:"data"`
}

// LiveRoomPaidDanmakuPushReq cmd=22 付费弹幕
type LiveRoomPaidDanmakuPushReq struct {
	g.Meta `path:"/liveRoomPaidDanmaku" method:"post" summary:"推送 cmd=22 付费弹幕(房间内全体在线用户)" description:"付费弹幕(房间内全体在线用户)" tags:"推送-直播间"`
}

type LiveRoomPaidDanmakuPushResp struct {
	Cmd  int                              `json:"cmd" dc:"命令字 22"`
	Data *liveroomdto.PaidDanmakuPushItem `json:"data"`
}

// LiveRoomAudienceLeavePushReq cmd=21 观众离房
type LiveRoomAudienceLeavePushReq struct {
	g.Meta `path:"/liveRoomAudienceLeave" method:"post" summary:"推送 cmd=21 观众离开直播间(房间内剩余在线用户)" description:"观众离开直播间(房间内剩余在线用户)" tags:"推送-直播间"`
}

type LiveRoomAudienceLeavePushResp struct {
	Cmd  int                                `json:"cmd" dc:"命令字 21"`
	Data *liveroomdto.AudienceLeavePushItem `json:"data"`
}

// LiveRoomAudienceJoinPushReq cmd=20 观众进房
type LiveRoomAudienceJoinPushReq struct {
	g.Meta `path:"/liveRoomAudienceJoin" method:"post" summary:"推送 cmd=20 观众进入直播间(房间内全体在线用户)" description:"观众进入直播间(房间内全体在线用户)" tags:"推送-直播间"`
}

type LiveRoomAudienceJoinPushResp struct {
	Cmd  int                               `json:"cmd" dc:"命令字 20"`
	Data *liveroomdto.AudienceJoinPushItem `json:"data"`
}

// LiveRoomAudienceKickCancelPushReq cmd=19 取消进入限制
type LiveRoomAudienceKickCancelPushReq struct {
	g.Meta `path:"/liveRoomAudienceKickCancel" method:"post" summary:"推送 cmd=19 取消观众进入限制(推送给被取消限制的用户)" description:"取消观众进入限制(推送给被取消限制的用户)" tags:"推送-直播间"`
}

type LiveRoomAudienceKickCancelPushResp struct {
	Cmd  int                                     `json:"cmd" dc:"命令字 19"`
	Data *liveroomdto.AudienceKickCancelPushItem `json:"data"`
}

// LiveRoomAudienceKickPushReq cmd=18 观众被踢
type LiveRoomAudienceKickPushReq struct {
	g.Meta `path:"/liveRoomAudienceKick" method:"post" summary:"推送 cmd=18 观众被踢出(推送给被踢用户)" description:"观众被踢出(推送给被踢用户)" tags:"推送-直播间"`
}

type LiveRoomAudienceKickPushResp struct {
	Cmd  int                               `json:"cmd" dc:"命令字 18"`
	Data *liveroomdto.AudienceKickPushItem `json:"data"`
}

// LiveRoomAudienceMutePushReq cmd=17 观众禁言
type LiveRoomAudienceMutePushReq struct {
	g.Meta `path:"/liveRoomAudienceMute" method:"post" summary:"推送 cmd=17 观众禁言/解禁状态(推送给被操作用户)" description:"观众禁言/解禁状态(推送给被操作用户)" tags:"推送-直播间"`
}

type LiveRoomAudienceMutePushResp struct {
	Cmd  int                               `json:"cmd" dc:"命令字 17"`
	Data *liveroomdto.AudienceMutePushItem `json:"data"`
}

// LiveRoomAnchorBanPushReq cmd=15 主播封禁
type LiveRoomAnchorBanPushReq struct {
	g.Meta `path:"/liveRoomAnchorBan" method:"post" summary:"推送 cmd=15 主播封禁(推送给主播及直播间在线观众)" description:"主播封禁(推送给主播及直播间在线观众)" tags:"推送-直播间"`
}

type LiveRoomAnchorBanPushResp struct {
	Cmd  int                            `json:"cmd" dc:"命令字 15"`
	Data *liveroomdto.AnchorBanPushItem `json:"data"`
}

// LiveRoomChatPushReq cmd=9 直播间文字消息
type LiveRoomChatPushReq struct {
	g.Meta `path:"/liveRoomChat" method:"post" summary:"推送 cmd=9 直播间免费文字消息(房间内全体在线用户)" description:"直播间免费文字消息(房间内全体在线用户)" tags:"推送-直播间"`
}

type LiveRoomChatPushResp struct {
	Cmd  int                       `json:"cmd" dc:"命令字 9"`
	Data *liveroomdto.ChatPushItem `json:"data"`
}

// LiveRoomGiftPushReq cmd=8 直播间送礼
type LiveRoomGiftPushReq struct {
	g.Meta `path:"/liveRoomGift" method:"post" summary:"推送 cmd=8 直播间送礼广播(房间内全体在线用户)" description:"直播间送礼广播(房间内全体在线用户)" tags:"推送-直播间"`
}

type LiveRoomGiftPushResp struct {
	Cmd  int                       `json:"cmd" dc:"命令字 8"`
	Data *liveroomdto.GiftPushItem `json:"data"`
}

// --- 推送-消息 ---

// SystemMessagePushReq cmd=14 系统消息
type SystemMessagePushReq struct {
	g.Meta `path:"/systemMessage" method:"post" summary:"推送 cmd=14 系统消息(推送给接收者)" description:"系统消息(推送给接收者)" tags:"推送-消息"`
}

type SystemMessagePushResp struct {
	Cmd  int                               `json:"cmd" dc:"命令字 14"`
	Data *messagedto.SystemMessagePushItem `json:"data"`
}

// --- 推送-私信 ---

// PrivateMessagePushReq cmd=13 私信
type PrivateMessagePushReq struct {
	g.Meta `path:"/privateMessage" method:"post" summary:"推送 cmd=13 私信消息(接收者与发送者均会收到)" description:"私信消息(接收者与发送者均会收到)" tags:"推送-私信"`
}

type PrivateMessagePushResp struct {
	Cmd  int                                `json:"cmd" dc:"命令字 13"`
	Data *messagedto.PrivateMessagePushItem `json:"data"`
}

// --- 推送-钱包 ---

// VipLevelPushReq cmd=12 VIP等级
type VipLevelPushReq struct {
	g.Meta `path:"/vipLevel" method:"post" summary:"推送 cmd=12 VIP等级变更(推送给指定用户)" description:"VIP等级变更(推送给指定用户)" tags:"推送-钱包"`
}

type VipLevelPushResp struct {
	Cmd  int                      `json:"cmd" dc:"命令字 12"`
	Data *vipdto.VipLevelPushItem `json:"data"`
}

// GoldPushReq cmd=11 金币余额
type GoldPushReq struct {
	g.Meta `path:"/gold" method:"post" summary:"推送 cmd=11 金币余额变更(推送给指定用户)" description:"金币余额变更(推送给指定用户)" tags:"推送-钱包"`
}

type GoldPushResp struct {
	Cmd  int                   `json:"cmd" dc:"命令字 11"`
	Data *golddto.GoldPushItem `json:"data"`
}

// DiamondPushReq cmd=10 钻石余额
type DiamondPushReq struct {
	g.Meta `path:"/diamond" method:"post" summary:"推送 cmd=10 钻石余额变更(推送给指定用户)" description:"钻石余额变更(推送给指定用户)" tags:"推送-钱包"`
}

type DiamondPushResp struct {
	Cmd  int                         `json:"cmd" dc:"命令字 10"`
	Data *diamonddto.DiamondPushItem `json:"data"`
}

// --- 推送-错误 ---

// ErrorParamPushReq cmd=7 带参数错误
type ErrorParamPushReq struct {
	g.Meta `path:"/errorParam" method:"post" summary:"推送 cmd=7 带参数错误" description:"带参数错误" tags:"推送-错误"`
}

type ErrorParamPushResp struct {
	Cmd  int            `json:"cmd" dc:"命令字 7"`
	Data *push.ErrorDto `json:"data"`
}

// ErrorPushReq cmd=6 无参数错误
type ErrorPushReq struct {
	g.Meta `path:"/error" method:"post" summary:"推送 cmd=6 无参数错误" description:"无参数错误" tags:"推送-错误"`
}

type ErrorPushResp struct {
	Cmd  int `json:"cmd" dc:"命令字 6"`
	Data int `json:"data" dc:"错误码 XRCode"`
}

// --- 推送-连接 ---

// RepeatLoginPushReq cmd=5 重复登录
type RepeatLoginPushReq struct {
	g.Meta `path:"/repeatLogin" method:"post" summary:"推送 cmd=5 重复登录，通知旧连接下线(仅 cmd，无 data)" description:"重复登录，通知旧连接下线(仅 cmd，无 data)" tags:"推送-连接"`
}

type RepeatLoginPushResp struct {
	Cmd int `json:"cmd" dc:"命令字 5(无 data)"`
}

// CloseServerPushReq cmd=4 关服通知
type CloseServerPushReq struct {
	g.Meta `path:"/closeServer" method:"post" summary:"推送 cmd=4 关服通知(仅 cmd，无 data)" description:"关服通知(仅 cmd，无 data)" tags:"推送-连接"`
}

type CloseServerPushResp struct {
	Cmd int `json:"cmd" dc:"命令字 4(无 data)"`
}

// KickPushReq cmd=3 踢下线
type KickPushReq struct {
	g.Meta `path:"/kick" method:"post" summary:"推送 cmd=3 踢下线" description:"踢下线" tags:"推送-连接"`
}

type KickPushResp struct {
	Cmd  int   `json:"cmd" dc:"命令字 3"`
	Data int64 `json:"data" dc:"踢出时间 Unix 毫秒时间戳"`
}

// HeartPushReq cmd=2 服务端心跳
type HeartPushReq struct {
	g.Meta `path:"/heart" method:"post" summary:"推送 cmd=2 服务端心跳保活" description:"服务端心跳保活" tags:"推送-连接"`
}

type HeartPushResp struct {
	Cmd  int   `json:"cmd" dc:"命令字 2"`
	Data int64 `json:"data" dc:"当前 Unix 时间戳(秒)"`
}

// EnterPushReq cmd=1 连接鉴权失败
type EnterPushReq struct {
	g.Meta `path:"/enter" method:"post" summary:"推送 cmd=1 WebSocket 连接鉴权失败时返回" description:"WebSocket 连接鉴权失败时返回" tags:"推送-连接"`
}

type EnterPushResp struct {
	Cmd  int                  `json:"cmd" dc:"命令字 1"`
	Data *httpserver.AuthResp `json:"data"`
}
