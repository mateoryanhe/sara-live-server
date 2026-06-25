package liveroomdto

import "github.com/gogf/gf/v2/frame/g"

// GetLiveRoomListReq App 分页查询直播间列表
type GetLiveRoomListReq struct {
	g.Meta       `path:"/roomList" method:"post" summary:"查询直播间列表" tags:"直播间"`
	Page         int    `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize     int    `json:"pageSize" dc:"每页数量(默认20,最大100)"`
	StatusFilter int    `json:"statusFilter" dc:"状态过滤(0=全部,1=仅直播中,2=仅未开播/已下播)"`
	TagId        uint64 `json:"tagId,string" dc:"标签ID(0=全部)"`
	Title        string `json:"title" dc:"标题(模糊匹配)"`
	Notice       string `json:"notice" dc:"公告(模糊匹配)"`
}

// LiveRoomListItem 直播间列表条目
type LiveRoomListItem struct {
	RoomId             string  `json:"roomId" dc:"直播间ID(同主播用户ID)"`
	GuildId            string  `json:"guildId" dc:"所属工会ID"`
	Title              string  `json:"title" dc:"直播间标题"`
	Cover              string  `json:"cover" dc:"封面图URL(已拼资源域名)"`
	Notice             string  `json:"notice" dc:"公告"`
	Status             uint8   `json:"status" dc:"状态(0未开播,1直播中)"`
	Category           uint8   `json:"category" dc:"分类(1=hot,2=game,3=私密)"`
	TagId              string  `json:"tagId" dc:"直播间标签ID"`
	TagName            string  `json:"tagName" dc:"直播间标签名称"`
	Ticket             float64 `json:"ticket" dc:"门票价格(钻石)"`
	Billing            float64 `json:"billing" dc:"计费价格(每分钟钻石)"`
	CreateAt           int64   `json:"createAt" dc:"创建时间(秒)"`
	AnchorNickname     string  `json:"anchorNickname" dc:"主播昵称"`
	AnchorAvatar       string  `json:"anchorAvatar" dc:"主播头像URL(已拼资源域名)"`
	AgoraToken         string  `json:"agoraToken" dc:"声网RTC Token"`
	AgoraTokenExpireAt int64   `json:"agoraTokenExpireAt" dc:"声网Token过期时间(Unix秒)"`
}

type GetLiveRoomListRes struct {
	Total    int                 `json:"total" dc:"总条数"`
	Page     int                 `json:"page" dc:"当前页码"`
	PageSize int                 `json:"pageSize" dc:"每页数量"`
	List     []*LiveRoomListItem `json:"list" dc:"直播间列表"`
}

// GetFollowedLiveRoomListReq App 分页查询当前用户关注的直播间列表
type GetFollowedLiveRoomListReq struct {
	g.Meta       `path:"/followedRoomList" method:"post" summary:"查询我关注的直播间列表" tags:"直播间"`
	Page         int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize     int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
	StatusFilter int `json:"statusFilter" dc:"状态过滤(0=全部,1=仅直播中,2=仅未开播/已下播)"`
}
