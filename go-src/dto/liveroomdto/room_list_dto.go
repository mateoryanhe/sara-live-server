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
	AllowCallIcon      bool    `json:"allowCallIcon" dc:"是否允许显示电话图标按钮(仅category=1时按私密邀请类型判断)"`
	CreateAt           int64   `json:"createAt" dc:"创建时间(秒)"`
	AnchorNickname     string  `json:"anchorNickname" dc:"主播昵称"`
	AnchorAvatar       string  `json:"anchorAvatar" dc:"主播头像URL(已拼资源域名)"`
	OnlineCount        int     `json:"onlineCount" dc:"在线观众人数(不含主播)"`
	AgoraToken         string  `json:"agoraToken" dc:"声网RTC Token"`
	AgoraTokenExpireAt int64   `json:"agoraTokenExpireAt" dc:"声网Token过期时间(Unix秒)"`
	IsBotAnchor        bool    `json:"isBotAnchor" dc:"是否机器人主播"`
	CloudPlayerVideo   string  `json:"cloudPlayerVideo" dc:"云播视频地址"`
	UserType           uint8   `json:"userType" dc:"主播用户类型(6=测试型主播)"`
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

const (
	NearbyLiveRoomDirectionDown = 1 // 往下列表序号更大的直播间(在线人数更少)
	NearbyLiveRoomDirectionUp   = 2 // 往下列表序号更小的直播间(在线人数更多)
)

// GetNearbyLiveRoomListReq App 以当前直播间为锚点,按列表顺序获取相邻的直播中直播间
type GetNearbyLiveRoomListReq struct {
	g.Meta    `path:"/nearbyRoomList" method:"post" summary:"查询相邻直播间列表" tags:"直播间"`
	RoomId    uint64 `json:"roomId,string" v:"required#直播间ID不能为空" dc:"当前直播间ID"`
	Direction int    `json:"direction" v:"required|in:1,2#方向不能为空|方向无效" dc:"方向(1=往下,2=往上)"`
	Count     int    `json:"count" dc:"获取数量(默认1,最大20)"`
}

type GetNearbyLiveRoomListRes struct {
	List []*LiveRoomListItem `json:"list" dc:"相邻直播间列表(不含当前房间,仅直播中)"`
}

// GetHotLiveRoomListReq App 分页查询 Hot 分类直播中房间列表(走内存缓存排序)
type GetHotLiveRoomListReq struct {
	g.Meta    `path:"/hotRoomList" method:"post" summary:"查询Hot直播中房间列表" tags:"直播间"`
	PageIndex int `json:"pageIndex" dc:"页码(从1开始,默认1)"`
	PageSize  int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// HotLiveRoomListItem Hot 直播间列表条目(含排名)
type HotLiveRoomListItem struct {
	LiveRoomListItem
	Rank int `json:"rank" dc:"排名(基于列表缓存排序,从1开始)"`
}

type GetHotLiveRoomListRes struct {
	Total     int                    `json:"total" dc:"总条数"`
	PageIndex int                    `json:"pageIndex" dc:"当前页码"`
	PageSize  int                    `json:"pageSize" dc:"每页数量"`
	List      []*HotLiveRoomListItem `json:"list" dc:"Hot直播中房间列表"`
}
