package messagedto

import "github.com/gogf/gf/v2/frame/g"

// AppActivityMessageListReq App端查询活动消息列表
type AppActivityMessageListReq struct {
	g.Meta `path:"/activityMessageList" method:"post" summary:"查询活动消息列表" tags:"活动消息"`
}

type AppActivityMessageItem struct {
	Id          uint64 `json:"id,string" dc:"活动消息ID"`
	IconEn      string `json:"iconEn"`
	IconEs      string `json:"iconEs"`
	IconPt      string `json:"iconPt"`
	IconHi      string `json:"iconHi"`
	BgEn        string `json:"bgEn"`
	BgEs        string `json:"bgEs"`
	BgPt        string `json:"bgPt"`
	BgHi        string `json:"bgHi"`
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
	PublishedAt string `json:"publishedAt"`
}

type AppActivityMessageListRes struct {
	List []*AppActivityMessageItem `json:"list"`
}
