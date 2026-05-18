package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ApplyJoinGuildReq 申请加入工会(userId 由鉴权中间件提供)
type ApplyJoinGuildReq struct {
	g.Meta  `path:"/applyJoin" method:"post" summary:"申请加入工会" tags:"直播工会"`
	GuildId uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
}

type ApplyJoinGuildRes struct {
	MemberId string `json:"memberId" dc:"工会成员记录ID"`
}

// ApproveGuildMemberReq 审批通过工会申请(由会长操作)
type ApproveGuildMemberReq struct {
	g.Meta   `path:"/approveMember" method:"post" summary:"审批通过工会申请" tags:"直播工会"`
	MemberId uint64 `json:"memberId" v:"required#成员记录ID不能为空" dc:"成员记录ID"`
}

type ApproveGuildMemberRes struct {
	Success bool `json:"success"`
}

// RejectGuildMemberReq 拒绝工会申请(由会长操作)
type RejectGuildMemberReq struct {
	g.Meta   `path:"/rejectMember" method:"post" summary:"拒绝工会申请" tags:"直播工会"`
	MemberId uint64 `json:"memberId" v:"required#成员记录ID不能为空" dc:"成员记录ID"`
}

type RejectGuildMemberRes struct {
	Success bool `json:"success"`
}

// KickGuildMemberReq 剔除工会成员(由会长操作)
type KickGuildMemberReq struct {
	g.Meta   `path:"/kickMember" method:"post" summary:"剔除工会成员" tags:"直播工会"`
	MemberId uint64 `json:"memberId" v:"required#成员记录ID不能为空" dc:"成员记录ID"`
}

type KickGuildMemberRes struct {
	Success bool `json:"success"`
}

// GuildMemberListReq 查询工会成员列表(仅返回已成功加入的成员)
type GuildMemberListReq struct {
	g.Meta  `path:"/memberList" method:"post" summary:"获取工会成员列表" tags:"直播工会"`
	GuildId uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
}

// PendingMemberListReq 查询工会待审批成员列表(仅会长可调用)
type PendingMemberListReq struct {
	g.Meta  `path:"/pendingMemberList" method:"post" summary:"获取工会待审批成员列表" tags:"直播工会"`
	GuildId uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
}

type PendingMemberListRes struct {
	List []*GuildMemberItem `json:"list"`
}

type GuildMemberItem struct {
	MemberId  string `json:"memberId"`
	GuildId   string `json:"guildId"`
	UserId    string `json:"userId"`
	Status    uint8  `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type GuildMemberListRes struct {
	List []*GuildMemberItem `json:"list"`
}
