package accountdto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSInviteeListReq CMS分页查询某用户邀请的下级用户
type CMSInviteeListReq struct {
	g.Meta `path:"/cmsInviteeList" method:"post" summary:"CMS查询用户邀请列表" tags:"账号"`
	httpserver.CMSQueryReq
	InviterId uint64 `json:"inviterId,string" v:"required#邀请人ID不能为空" dc:"邀请人用户ID"`
	StartTime int64  `json:"startTime" dc:"开始时间(Unix秒,按被邀请人资料创建时间过滤,可选)"`
	EndTime   int64  `json:"endTime" dc:"结束时间(Unix秒,按被邀请人资料创建时间过滤,可选)"`
}

// CMSInviteeItem 被邀请用户列表项
type CMSInviteeItem struct {
	UserId    uint64     `json:"userId,string" dc:"被邀请人用户ID"`
	Nickname  string     `json:"nickname" dc:"昵称"`
	Avatar    string     `json:"avatar" dc:"头像URL"`
	CreatedAt *time.Time `json:"createdAt" dc:"资料创建时间"`
}
