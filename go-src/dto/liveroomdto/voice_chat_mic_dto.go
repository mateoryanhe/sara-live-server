package liveroomdto

import "github.com/gogf/gf/v2/frame/g"

// VoiceChatMicUserBrief 麦位/申请列表用户信息
type VoiceChatMicUserBrief struct {
	UserId   string `json:"userId" dc:"用户ID"`
	Nickname string `json:"nickname" dc:"昵称"`
	Avatar   string `json:"avatar" dc:"头像URL"`
	VipLevel uint32 `json:"vipLevel" dc:"VIP等级"`
	Gender   uint8  `json:"gender" dc:"性别(0未知,1男,2女)"`
}

// VoiceChatMicSeatItem 单个麦位状态(POST /liveRoom/voiceChatMicState 与 WS cmd41 推送结构一致)
type VoiceChatMicSeatItem struct {
	SeatIndex int                    `json:"seatIndex" dc:"麦位序号(0=主播位,1~8=观众位)"`
	Status    uint8                  `json:"status" dc:"麦位状态(App按status选麦位底图/角标):1=空闲(空麦位,可上麦) 2=锁麦(空麦位,主播已锁,不可申请/上麦) 3=在麦(有人,可说话,展示user头像) 4=在麦禁音(有人,被禁音,展示user头像+禁音标识)"`
	Occupied  bool                   `json:"occupied" dc:"是否有人(true=展示头像;与status=3或4一致,false=status为1或2)"`
	User      *VoiceChatMicUserBrief `json:"user" dc:"麦上用户(status=3或4时有值;status=1或2时为null)"`
}

// VoiceChatMicApplyResultPushItem 上麦申请结果推送给申请人(WS cmd43)
type VoiceChatMicApplyResultPushItem struct {
	RoomId    string `json:"roomId" dc:"直播间ID"`
	Approved  bool   `json:"approved" dc:"是否同意(true=同意上麦,false=拒绝)"`
	SeatIndex int    `json:"seatIndex" dc:"同意时占用的麦位(1~8;拒绝时为0)"`
}

// VoiceChatMicKickedPushItem 被房主抱下麦推送给该观众(WS cmd44)
type VoiceChatMicKickedPushItem struct {
	RoomId    string `json:"roomId" dc:"直播间ID"`
	SeatIndex int    `json:"seatIndex" dc:"被抱下前所在麦位(1~8)"`
}

// VoiceChatMicSeatMutePushItem 麦位禁音/解禁推送给麦上用户
type VoiceChatMicSeatMutePushItem struct {
	RoomId    string `json:"roomId" dc:"直播间ID"`
	SeatIndex int    `json:"seatIndex" dc:"麦位序号"`
	Muted     bool   `json:"muted" dc:"是否禁音(true=该麦位已被主播禁音,对应status=4;客户端应关闭推流或展示禁音态)"`
}

// VoiceChatMicApplyItem 上麦申请条目
type VoiceChatMicApplyItem struct {
	User      *VoiceChatMicUserBrief `json:"user" dc:"申请用户"`
	AppliedAt int64                  `json:"appliedAt" dc:"申请时间(秒)"`
}

// VoiceChatMicStatePayload 麦位全量状态(接口响应与推送共用)
type VoiceChatMicStatePayload struct {
	RoomId           string                  `json:"roomId" dc:"直播间ID"`
	VoiceChatMicMode uint8                   `json:"voiceChatMicMode" dc:"上麦方式(1=自由上麦,2=申请上麦,3=房主单麦)"`
	UpdatedAt        int64                   `json:"updatedAt" dc:"状态更新时间(毫秒)"`
	Seats            []*VoiceChatMicSeatItem `json:"seats" dc:"麦位列表(固定9项,seatIndex 0~8;每项看status渲染UI)"`
}

// GetVoiceChatMicStateReq 查询语聊房麦位状态
type GetVoiceChatMicStateReq struct {
	g.Meta `path:"/voiceChatMicState" method:"post" summary:"查询语聊房麦位状态" tags:"语聊房"`
	RoomId uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
}

type GetVoiceChatMicStateRes = VoiceChatMicStatePayload

// ApplyVoiceChatMicReq 观众申请上麦
type ApplyVoiceChatMicReq struct {
	g.Meta    `path:"/applyVoiceChatMic" method:"post" summary:"申请上麦" tags:"语聊房"`
	RoomId    uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	SeatIndex int    `json:"seatIndex" dc:"自由上麦时可指定麦位(1~8,0或不传则自动分配;申请模式忽略)"`
}

type ApplyVoiceChatMicRes struct {
	Success   bool `json:"success"`
	SeatIndex int  `json:"seatIndex" dc:"自由上麦成功时占用的麦位序号(1~8)"`
}

// CancelVoiceChatMicApplyReq 取消上麦申请
type CancelVoiceChatMicApplyReq struct {
	g.Meta `path:"/cancelVoiceChatMicApply" method:"post" summary:"取消上麦申请" tags:"语聊房"`
	RoomId uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
}

type CancelVoiceChatMicApplyRes struct {
	Success bool `json:"success"`
}

// VoiceChatMicApplyListReq 主播查询上麦申请列表
type VoiceChatMicApplyListReq struct {
	g.Meta `path:"/voiceChatMicApplyList" method:"post" summary:"上麦申请列表" tags:"语聊房"`
	RoomId uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
}

type VoiceChatMicApplyListRes struct {
	RoomId string                   `json:"roomId" dc:"直播间ID"`
	Total  int                      `json:"total" dc:"申请总数"`
	List   []*VoiceChatMicApplyItem `json:"list" dc:"申请列表(按申请时间先后)"`
}

// ApproveVoiceChatMicApplyReq 主播同意上麦申请
type ApproveVoiceChatMicApplyReq struct {
	g.Meta    `path:"/approveVoiceChatMicApply" method:"post" summary:"同意上麦申请" tags:"语聊房"`
	RoomId    uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	UserId    uint64 `json:"userId" v:"required|min:1#用户ID不能为空|用户ID无效" dc:"申请用户ID"`
	SeatIndex int    `json:"seatIndex" dc:"指定麦位(1~8,0或不传则自动分配空位)"`
}

type ApproveVoiceChatMicApplyRes struct {
	Success   bool `json:"success"`
	SeatIndex int  `json:"seatIndex" dc:"分配到的麦位序号"`
}

// RejectVoiceChatMicApplyReq 主播拒绝上麦申请
type RejectVoiceChatMicApplyReq struct {
	g.Meta `path:"/rejectVoiceChatMicApply" method:"post" summary:"拒绝上麦申请" tags:"语聊房"`
	RoomId uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	UserId uint64 `json:"userId" v:"required|min:1#用户ID不能为空|用户ID无效" dc:"申请用户ID"`
}

type RejectVoiceChatMicApplyRes struct {
	Success bool `json:"success"`
}

// LeaveVoiceChatMicReq 观众主动下麦
type LeaveVoiceChatMicReq struct {
	g.Meta `path:"/leaveVoiceChatMic" method:"post" summary:"下麦" tags:"语聊房"`
	RoomId uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
}

type LeaveVoiceChatMicRes struct {
	Success bool `json:"success"`
}

// KickVoiceChatMicReq 主播将观众抱下麦
type KickVoiceChatMicReq struct {
	g.Meta `path:"/kickVoiceChatMic" method:"post" summary:"抱下麦" tags:"语聊房"`
	RoomId uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	UserId uint64 `json:"userId" v:"required|min:1#用户ID不能为空|用户ID无效" dc:"被下麦用户ID"`
}

type KickVoiceChatMicRes struct {
	Success bool `json:"success"`
}

// LockVoiceChatMicSeatReq 主播锁麦(空观众位)
type LockVoiceChatMicSeatReq struct {
	g.Meta    `path:"/lockVoiceChatMicSeat" method:"post" summary:"锁麦" tags:"语聊房"`
	RoomId    uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	SeatIndex int    `json:"seatIndex" v:"required|min:1|max:8#麦位序号不能为空|麦位序号无效" dc:"观众麦位序号(1~8)"`
}

type LockVoiceChatMicSeatRes struct {
	Success bool `json:"success"`
}

// UnlockVoiceChatMicSeatReq 主播解锁麦位
type UnlockVoiceChatMicSeatReq struct {
	g.Meta    `path:"/unlockVoiceChatMicSeat" method:"post" summary:"解锁麦位" tags:"语聊房"`
	RoomId    uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	SeatIndex int    `json:"seatIndex" v:"required|min:1|max:8#麦位序号不能为空|麦位序号无效" dc:"观众麦位序号(1~8)"`
}

type UnlockVoiceChatMicSeatRes struct {
	Success bool `json:"success"`
}

// MuteVoiceChatMicSeatReq 主播对指定麦位禁音(在麦用户不可说话)
type MuteVoiceChatMicSeatReq struct {
	g.Meta    `path:"/muteVoiceChatMicSeat" method:"post" summary:"麦位禁音" tags:"语聊房"`
	RoomId    uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	SeatIndex int    `json:"seatIndex" v:"required|min:0|max:8#麦位序号不能为空|麦位序号无效" dc:"麦位序号(0~8)"`
}

type MuteVoiceChatMicSeatRes struct {
	Success bool `json:"success"`
}

// UnmuteVoiceChatMicSeatReq 主播解除麦位禁音
type UnmuteVoiceChatMicSeatReq struct {
	g.Meta    `path:"/unmuteVoiceChatMicSeat" method:"post" summary:"解除麦位禁音" tags:"语聊房"`
	RoomId    uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	SeatIndex int    `json:"seatIndex" v:"required|min:0|max:8#麦位序号不能为空|麦位序号无效" dc:"麦位序号(0~8)"`
}

type UnmuteVoiceChatMicSeatRes struct {
	Success bool `json:"success"`
}

// InviteVoiceChatMicReq 主播抱观众上麦(无需申请队列)
type InviteVoiceChatMicReq struct {
	g.Meta    `path:"/inviteVoiceChatMic" method:"post" summary:"抱上麦" tags:"语聊房"`
	RoomId    uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	UserId    uint64 `json:"userId" v:"required|min:1#用户ID不能为空|用户ID无效" dc:"被抱上麦的观众ID"`
	SeatIndex int    `json:"seatIndex" dc:"指定麦位(1~8,0或不传则自动分配空位)"`
}

type InviteVoiceChatMicRes struct {
	Success   bool `json:"success"`
	SeatIndex int  `json:"seatIndex" dc:"分配到的麦位序号(1~8)"`
}

// SwitchVoiceChatMicSeatReq 在麦观众换麦位
type SwitchVoiceChatMicSeatReq struct {
	g.Meta    `path:"/switchVoiceChatMicSeat" method:"post" summary:"换麦" tags:"语聊房"`
	RoomId    uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	SeatIndex int    `json:"seatIndex" v:"required|min:1|max:8#目标麦位不能为空|目标麦位无效" dc:"目标观众麦位(1~8)"`
}

type SwitchVoiceChatMicSeatRes struct {
	Success       bool `json:"success"`
	FromSeatIndex int  `json:"fromSeatIndex" dc:"原麦位序号"`
	SeatIndex     int  `json:"seatIndex" dc:"当前麦位序号"`
}

// BatchLockVoiceChatMicSeatsReq 主播批量锁麦(空观众位)
type BatchLockVoiceChatMicSeatsReq struct {
	g.Meta      `path:"/batchLockVoiceChatMicSeats" method:"post" summary:"批量锁麦" tags:"语聊房"`
	RoomId      uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	SeatIndexes []int  `json:"seatIndexes" v:"required#麦位列表不能为空" dc:"观众麦位序号列表(1~8,可多个)"`
}

type BatchLockVoiceChatMicSeatsRes struct {
	Success bool `json:"success"`
	Count   int  `json:"count" dc:"成功锁定的麦位数"`
}

// BatchUnlockVoiceChatMicSeatsReq 主播批量解锁麦位
type BatchUnlockVoiceChatMicSeatsReq struct {
	g.Meta      `path:"/batchUnlockVoiceChatMicSeats" method:"post" summary:"批量解锁麦位" tags:"语聊房"`
	RoomId      uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
	SeatIndexes []int  `json:"seatIndexes" v:"required#麦位列表不能为空" dc:"观众麦位序号列表(1~8,可多个)"`
}

type BatchUnlockVoiceChatMicSeatsRes struct {
	Success bool `json:"success"`
	Count   int  `json:"count" dc:"成功解锁的麦位数"`
}

// ClearVoiceChatMicAppliesReq 主播清空上麦申请列表
type ClearVoiceChatMicAppliesReq struct {
	g.Meta `path:"/clearVoiceChatMicApplies" method:"post" summary:"清空上麦申请" tags:"语聊房"`
	RoomId uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
}

type ClearVoiceChatMicAppliesRes struct {
	Success bool `json:"success"`
}
