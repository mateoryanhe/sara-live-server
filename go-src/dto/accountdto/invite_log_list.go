package accountdto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSInviteLogListReq CMS分页查询邀请日志(全局)
type CMSInviteLogListReq struct {
	g.Meta `path:"/cmsInviteLogList" method:"post" summary:"CMS查询邀请日志列表" tags:"账号"`
	httpserver.CMSQueryReq
	Keyword   string `json:"keyword" dc:"关键字(被邀请人/邀请人用户ID或昵称,可选)"`
	StartTime int64  `json:"startTime" dc:"开始时间(Unix秒,按被邀请人资料创建时间,可选)"`
	EndTime   int64  `json:"endTime" dc:"结束时间(Unix秒,按被邀请人资料创建时间,可选)"`
}

// CMSInviteLogItem 邀请日志列表项
type CMSInviteLogItem struct {
	UserId           uint64     `json:"userId,string" dc:"被邀请人用户ID"`
	Nickname         string     `json:"nickname" dc:"被邀请人昵称"`
	Avatar           string     `json:"avatar" dc:"被邀请人头像URL"`
	InviterId        uint64     `json:"inviterId,string" dc:"邀请人用户ID"`
	InviterNickname  string     `json:"inviterNickname" dc:"邀请人昵称"`
	InviterAvatar    string     `json:"inviterAvatar" dc:"邀请人头像URL"`
	CreatedAt        *time.Time `json:"createdAt" dc:"被邀请人资料创建时间"`
}
