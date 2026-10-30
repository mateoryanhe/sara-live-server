package userinfodto

import "github.com/gogf/gf/v2/frame/g"

// InviteeListReq 名下邀请人列表
type InviteeListReq struct {
	g.Meta   `path:"/inviteeList" method:"post" summary:"名下邀请人列表" tags:"用户信息"`
	Page     int `json:"page"     v:"min:1#页码从1开始" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100;前2页走缓存)"`
}

// InviteeListItem 被邀请用户条目
type InviteeListItem struct {
	UserId            string  `json:"userId" dc:"被邀请人用户ID"`
	Nickname          string  `json:"nickname" dc:"昵称"`
	Avatar            string  `json:"avatar" dc:"头像URL(已拼资源域名)"`
	InviterRewardGold float64 `json:"inviterRewardGold" dc:"邀请人从该用户充值获得的返现金币累计"`
}

type InviteeListRes struct {
	Page     int                `json:"page" dc:"当前页码"`
	PageSize int                `json:"pageSize" dc:"每页数量"`
	List     []*InviteeListItem `json:"list" dc:"邀请人列表"`
}
