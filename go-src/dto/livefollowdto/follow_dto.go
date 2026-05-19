package livefollowdto

import "github.com/gogf/gf/v2/frame/g"

// FollowReq 关注主播
type FollowReq struct {
	g.Meta   `path:"/follow" method:"post" summary:"关注主播" tags:"关注"`
	AnchorId uint64 `json:"anchorId" v:"required#主播ID不能为空" dc:"被关注的主播用户ID"`
}

type FollowRes struct {
	Following     bool `json:"following"     dc:"当前是否已关注"`
	FollowerCount int  `json:"followerCount" dc:"主播当前粉丝数"`
}

// UnfollowReq 取消关注
type UnfollowReq struct {
	g.Meta   `path:"/unfollow" method:"post" summary:"取消关注主播" tags:"关注"`
	AnchorId uint64 `json:"anchorId" v:"required#主播ID不能为空" dc:"主播用户ID"`
}

type UnfollowRes struct {
	Following     bool `json:"following"     dc:"当前是否已关注"`
	FollowerCount int  `json:"followerCount" dc:"主播当前粉丝数"`
}

// IsFollowingReq 查询是否关注
type IsFollowingReq struct {
	g.Meta   `path:"/isFollowing" method:"post" summary:"查询是否关注" tags:"关注"`
	AnchorId uint64 `json:"anchorId" v:"required#主播ID不能为空" dc:"主播用户ID"`
}

type IsFollowingRes struct {
	Following bool `json:"following" dc:"是否已关注"`
}

// FollowingListReq 当前用户关注的主播分页列表
type FollowingListReq struct {
	g.Meta   `path:"/followingList" method:"post" summary:"我关注的主播列表" tags:"关注"`
	Page     int `json:"page"     v:"min:1#页码从1开始" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100)"`
}

// AnchorItem 关注主播条目(附带主播基础信息)
type AnchorItem struct {
	AnchorId   string `json:"anchorId"   dc:"主播用户ID"`
	Nickname   string `json:"nickname"   dc:"主播昵称"`
	Avatar     string `json:"avatar"     dc:"主播头像URL(已拼资源域名)"`
	GuildId    string `json:"guildId"    dc:"主播所属工会ID"`
	RoomStatus uint8  `json:"roomStatus" dc:"主播直播间状态(0未开播,1直播中,2无直播间)"`
	FollowedAt int64  `json:"followedAt" dc:"关注时间(秒)"`
}

type FollowingListRes struct {
	Total    int           `json:"total"    dc:"总数"`
	Page     int           `json:"page"     dc:"当前页码"`
	PageSize int           `json:"pageSize" dc:"每页数量"`
	List     []*AnchorItem `json:"list"     dc:"主播列表"`
}

// FollowerListReq 主播粉丝列表(查询自己的)
type FollowerListReq struct {
	g.Meta   `path:"/followerList" method:"post" summary:"我的粉丝列表" tags:"关注"`
	Page     int `json:"page"     v:"min:1#页码从1开始" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100)"`
}

// FollowerItem 粉丝条目
type FollowerItem struct {
	UserId     string `json:"userId"     dc:"粉丝用户ID"`
	Nickname   string `json:"nickname"   dc:"粉丝昵称"`
	Avatar     string `json:"avatar"     dc:"粉丝头像URL(已拼资源域名)"`
	FollowedAt int64  `json:"followedAt" dc:"关注时间(秒)"`
}

type FollowerListRes struct {
	Total    int             `json:"total"    dc:"总数"`
	Page     int             `json:"page"     dc:"当前页码"`
	PageSize int             `json:"pageSize" dc:"每页数量"`
	List     []*FollowerItem `json:"list"     dc:"粉丝列表"`
}
