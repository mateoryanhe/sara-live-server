package gameplatformdto

import "github.com/gogf/gf/v2/frame/g"

// AppGameListReq App 分页查询游戏列表(仅已上架,读内存)
type AppGameListReq struct {
	g.Meta    `path:"/appGameList" method:"post" summary:"App分页查询游戏列表(已上架)" tags:"游戏"`
	PageIndex int `json:"pageIndex" dc:"页码(从1开始,默认1)"`
	PageSize  int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// AppGameListItem App 游戏列表项
type AppGameListItem struct {
	GameCode string `json:"gameCode" dc:"游戏编码"`
	NameEn   string `json:"nameEn" dc:"英文名称"`
	Cover    string `json:"cover" dc:"封面完整URL"`
	Category string `json:"category" dc:"分类"`
	Platform string `json:"platform" dc:"平台"`
}

// AppGameListRes App 游戏分页列表响应
type AppGameListRes struct {
	Total     int                `json:"total" dc:"总条数"`
	PageIndex int                `json:"pageIndex" dc:"当前页码"`
	PageSize  int                `json:"pageSize" dc:"每页数量"`
	List      []*AppGameListItem `json:"list" dc:"游戏列表"`
}
