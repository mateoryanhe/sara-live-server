package livefollowdto

import "github.com/gogf/gf/v2/frame/g"

// BlockReq 拉黑用户
type BlockReq struct {
	g.Meta   `path:"/block" method:"post" summary:"拉黑用户" tags:"拉黑"`
	TargetId uint64 `json:"targetId" v:"required#目标用户ID不能为空" dc:"被拉黑的用户ID"`
}

type BlockRes struct {
	Success bool `json:"success" dc:"是否成功"`
	Blocked bool `json:"blocked" dc:"当前是否已拉黑"`
}

// UnblockReq 解除拉黑
type UnblockReq struct {
	g.Meta   `path:"/unblock" method:"post" summary:"解除拉黑" tags:"拉黑"`
	TargetId uint64 `json:"targetId" v:"required#目标用户ID不能为空" dc:"被解除拉黑的用户ID"`
}

type UnblockRes struct {
	Success bool `json:"success" dc:"是否成功"`
	Blocked bool `json:"blocked" dc:"当前是否仍被拉黑"`
}

// BlockListReq 拉黑列表
type BlockListReq struct {
	g.Meta   `path:"/blockList" method:"post" summary:"拉黑列表" tags:"拉黑"`
	Page     int `json:"page"     v:"min:1#页码从1开始" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100)"`
}

// BlockListItem 拉黑用户条目
type BlockListItem struct {
	TargetId  string `json:"targetId"  dc:"被拉黑用户ID"`
	Nickname  string `json:"nickname"  dc:"昵称"`
	Avatar    string `json:"avatar"    dc:"头像URL(已拼资源域名)"`
	VipLevel  uint32 `json:"vipLevel"  dc:"VIP等级"`
	Gender    uint8  `json:"gender"    dc:"性别(0未知,1男,2女)"`
	Age       int    `json:"age"       dc:"年龄(未设置出生日期时为0)"`
	BlockedAt int64  `json:"blockedAt" dc:"拉黑时间(秒)"`
}

type BlockListRes struct {
	Total    int              `json:"total"    dc:"总数"`
	Page     int              `json:"page"     dc:"当前页码"`
	PageSize int              `json:"pageSize" dc:"每页数量"`
	List     []*BlockListItem `json:"list"     dc:"拉黑列表"`
}
