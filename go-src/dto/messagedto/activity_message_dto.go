package messagedto

import "github.com/gogf/gf/v2/frame/g"

// AppActivityMessageListReq App端查询活动消息列表
type AppActivityMessageListReq struct {
	g.Meta `path:"/activityMessageList" method:"post" summary:"查询活动消息列表" tags:"活动消息"`
}

type AppActivityMessageItem struct {
	Id          uint64 `json:"id,string" dc:"活动消息ID"`
	IconEn      string `json:"iconEn" dc:"图标URL(英文)"`
	IconEs      string `json:"iconEs" dc:"图标URL(西班牙语)"`
	IconPt      string `json:"iconPt" dc:"图标URL(葡萄牙语)"`
	IconHi      string `json:"iconHi" dc:"图标URL(印地语)"`
	BgEn        string `json:"bgEn" dc:"背景图URL(英文)"`
	BgEs        string `json:"bgEs" dc:"背景图URL(西班牙语)"`
	BgPt        string `json:"bgPt" dc:"背景图URL(葡萄牙语)"`
	BgHi        string `json:"bgHi" dc:"背景图URL(印地语)"`
	TitleEn     string `json:"titleEn" dc:"标题(英文)"`
	TitleEs     string `json:"titleEs" dc:"标题(西班牙语)"`
	TitlePt     string `json:"titlePt" dc:"标题(葡萄牙语)"`
	TitleHi     string `json:"titleHi" dc:"标题(印地语)"`
	ContentEn   string `json:"contentEn" dc:"内容(英文)"`
	ContentEs   string `json:"contentEs" dc:"内容(西班牙语)"`
	ContentPt   string `json:"contentPt" dc:"内容(葡萄牙语)"`
	ContentHi   string `json:"contentHi" dc:"内容(印地语)"`
	UrlEn       string `json:"urlEn" dc:"跳转链接(英文)"`
	UrlEs       string `json:"urlEs" dc:"跳转链接(西班牙语)"`
	UrlPt       string `json:"urlPt" dc:"跳转链接(葡萄牙语)"`
	UrlHi       string `json:"urlHi" dc:"跳转链接(印地语)"`
	PublishedAt string `json:"publishedAt" dc:"发布时间"`
}

type AppActivityMessageListRes struct {
	List []*AppActivityMessageItem `json:"list" dc:"活动消息列表"`
}
