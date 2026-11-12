package gameplatformdto

// BrowseTitlesReq Sofie 标题列表(对应旧 AppGameList;字段名避免 game).
type BrowseTitlesReq struct {
	PageIndex int `json:"pageIndex" dc:"页码(从1开始,默认1)"`
	PageSize  int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// BrowseTitlesItem Sofie 标题列表项.
type BrowseTitlesItem struct {
	TitleCode string `json:"titleCode" dc:"标题编码"`
	NameEn    string `json:"nameEn" dc:"英文名称"`
	Cover     string `json:"cover" dc:"封面完整URL"`
	Category  string `json:"category" dc:"分类"`
	Platform  string `json:"platform" dc:"平台"`
}

// BrowseTitlesRes Sofie 标题分页列表.
type BrowseTitlesRes struct {
	Total     int                 `json:"total" dc:"总条数"`
	PageIndex int                 `json:"pageIndex" dc:"当前页码"`
	PageSize  int                 `json:"pageSize" dc:"每页数量"`
	List      []*BrowseTitlesItem `json:"list" dc:"标题列表"`
}
