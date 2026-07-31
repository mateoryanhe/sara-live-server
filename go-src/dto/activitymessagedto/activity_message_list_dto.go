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
	ID          string `json:"id"`
	IconEn      string `json:"iconEn" dc:"图标URL(英文)"`
	IconEnName  string `json:"iconEnName" dc:"图标资源名(英文,编辑保存用)"`
	IconEs      string `json:"iconEs" dc:"图标URL(西班牙语)"`
	IconEsName  string `json:"iconEsName" dc:"图标资源名(西班牙语,编辑保存用)"`
	IconPt      string `json:"iconPt" dc:"图标URL(葡萄牙语)"`
	IconPtName  string `json:"iconPtName" dc:"图标资源名(葡萄牙语,编辑保存用)"`
	IconHi      string `json:"iconHi" dc:"图标URL(印地语)"`
	IconHiName  string `json:"iconHiName" dc:"图标资源名(印地语,编辑保存用)"`
	BgEn        string `json:"bgEn" dc:"背景图URL(英文)"`
	BgEnName    string `json:"bgEnName" dc:"背景图资源名(英文,编辑保存用)"`
	BgEs        string `json:"bgEs" dc:"背景图URL(西班牙语)"`
	BgEsName    string `json:"bgEsName" dc:"背景图资源名(西班牙语,编辑保存用)"`
	BgPt        string `json:"bgPt" dc:"背景图URL(葡萄牙语)"`
	BgPtName    string `json:"bgPtName" dc:"背景图资源名(葡萄牙语,编辑保存用)"`
	BgHi        string `json:"bgHi" dc:"背景图URL(印地语)"`
	BgHiName    string `json:"bgHiName" dc:"背景图资源名(印地语,编辑保存用)"`
	TitleEn     string `json:"titleEn"`
	TitleEs     string `json:"titleEs"`
	TitlePt     string `json:"titlePt"`
	TitleHi     string `json:"titleHi"`
	ContentEn   string `json:"contentEn"`
	ContentEs   string `json:"contentEs"`
	ContentPt   string `json:"contentPt"`
	ContentHi   string `json:"contentHi"`
	UrlEn       string `json:"urlEn"`
	UrlEs       string `json:"urlEs"`
	UrlPt       string `json:"urlPt"`
	UrlHi       string `json:"urlHi"`
	Status      uint8  `json:"status"`
	PublishedAt string `json:"publishedAt"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
