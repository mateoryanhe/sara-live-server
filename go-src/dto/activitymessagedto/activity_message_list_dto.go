package activitymessagedto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type ActivityMessageListReq struct {
	g.Meta `path:"/activityMessageList" method:"post" summary:"获取活动消息列表" tags:"活动消息"`
	httpserver.CMSQueryReq
	Title        string `json:"title" dc:"标题(模糊匹配,任意语言)"`
	StatusFilter int    `json:"statusFilter" dc:"发布状态过滤(0=全部, 1=只看未发布, 2=只看已发布)"`
}

type ActivityMessageListRes struct {
	ID          string `json:"id" dc:"活动消息ID"`
	IconEn      string `json:"iconEn" dc:"图标URL(英文)"`
	IconEnName  string `json:"iconEnName" dc:"图标资源名(英文,编辑保存用)"`
	IconEs      string `json:"iconEs" dc:"图标URL(西班牙语)"`
	IconEsName  string `json:"iconEsName" dc:"图标资源名(西班牙语,编辑保存用)"`
	IconPt      string `json:"iconPt" dc:"图标URL(葡萄牙语)"`
	IconPtName  string `json:"iconPtName" dc:"图标资源名(葡萄牙语,编辑保存用)"`
	IconHi      string `json:"iconHi" dc:"图标URL(印地语)"`
	IconHiName  string `json:"iconHiName" dc:"图标资源名(印地语,编辑保存用)"`
	IconId      string `json:"iconId" dc:"图标URL(印尼语)"`
	IconIdName  string `json:"iconIdName" dc:"图标资源名(印尼语,编辑保存用)"`
	BgEn        string `json:"bgEn" dc:"背景图URL(英文)"`
	BgEnName    string `json:"bgEnName" dc:"背景图资源名(英文,编辑保存用)"`
	BgEs        string `json:"bgEs" dc:"背景图URL(西班牙语)"`
	BgEsName    string `json:"bgEsName" dc:"背景图资源名(西班牙语,编辑保存用)"`
	BgPt        string `json:"bgPt" dc:"背景图URL(葡萄牙语)"`
	BgPtName    string `json:"bgPtName" dc:"背景图资源名(葡萄牙语,编辑保存用)"`
	BgHi        string `json:"bgHi" dc:"背景图URL(印地语)"`
	BgHiName    string `json:"bgHiName" dc:"背景图资源名(印地语,编辑保存用)"`
	BgId        string `json:"bgId" dc:"背景图URL(印尼语)"`
	BgIdName    string `json:"bgIdName" dc:"背景图资源名(印尼语,编辑保存用)"`
	TitleEn     string `json:"titleEn" dc:"标题(英文)"`
	TitleEs     string `json:"titleEs" dc:"标题(西班牙语)"`
	TitlePt     string `json:"titlePt" dc:"标题(葡萄牙语)"`
	TitleHi     string `json:"titleHi" dc:"标题(印地语)"`
	TitleId     string `json:"titleId" dc:"标题(印尼语)"`
	ContentEn   string `json:"contentEn" dc:"内容(英文)"`
	ContentEs   string `json:"contentEs" dc:"内容(西班牙语)"`
	ContentPt   string `json:"contentPt" dc:"内容(葡萄牙语)"`
	ContentHi   string `json:"contentHi" dc:"内容(印地语)"`
	ContentId   string `json:"contentId" dc:"内容(印尼语)"`
	UrlEn       string `json:"urlEn" dc:"跳转链接(英文)"`
	UrlEs       string `json:"urlEs" dc:"跳转链接(西班牙语)"`
	UrlPt       string `json:"urlPt" dc:"跳转链接(葡萄牙语)"`
	UrlHi       string `json:"urlHi" dc:"跳转链接(印地语)"`
	UrlId       string `json:"urlId" dc:"跳转链接(印尼语)"`
	Status      uint8  `json:"status" dc:"发布状态(0未发布,1已发布)"`
	PublishedAt string `json:"publishedAt" dc:"发布时间"`
	CreatedAt   string `json:"createdAt" dc:"创建时间"`
	UpdatedAt   string `json:"updatedAt" dc:"更新时间"`
}
