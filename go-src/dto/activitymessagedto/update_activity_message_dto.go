package activitymessagedto

import "github.com/gogf/gf/v2/frame/g"

type UpdateActivityMessageReq struct {
	g.Meta    `path:"/updateActivityMessage" method:"post" summary:"修改活动消息" tags:"活动消息"`
	ID        uint64 `json:"id" v:"required#活动消息ID不能为空" dc:"活动消息ID"`
	IconEn    string `json:"iconEn" v:"max-length:255#图标资源名(英文)最长255字符" dc:"图标资源名(英文)"`
	IconEs    string `json:"iconEs" v:"max-length:255#图标资源名(西班牙语)最长255字符" dc:"图标资源名(西班牙语)"`
	IconPt    string `json:"iconPt" v:"max-length:255#图标资源名(葡萄牙语)最长255字符" dc:"图标资源名(葡萄牙语)"`
	IconHi    string `json:"iconHi" v:"max-length:255#图标资源名(印地语)最长255字符" dc:"图标资源名(印地语)"`
	BgEn      string `json:"bgEn" v:"max-length:255#背景图资源名(英文)最长255字符" dc:"背景图资源名(英文)"`
	BgEs      string `json:"bgEs" v:"max-length:255#背景图资源名(西班牙语)最长255字符" dc:"背景图资源名(西班牙语)"`
	BgPt      string `json:"bgPt" v:"max-length:255#背景图资源名(葡萄牙语)最长255字符" dc:"背景图资源名(葡萄牙语)"`
	BgHi      string `json:"bgHi" v:"max-length:255#背景图资源名(印地语)最长255字符" dc:"背景图资源名(印地语)"`
	TitleEn   string `json:"titleEn" v:"required|max-length:128#标题(英文)不能为空|标题(英文)最长128字符" dc:"标题(英文)"`
	TitleEs   string `json:"titleEs" v:"max-length:128#标题(西班牙语)最长128字符" dc:"标题(西班牙语)"`
	TitlePt   string `json:"titlePt" v:"max-length:128#标题(葡萄牙语)最长128字符" dc:"标题(葡萄牙语)"`
	TitleHi   string `json:"titleHi" v:"max-length:128#标题(印地语)最长128字符" dc:"标题(印地语)"`
	ContentEn string `json:"contentEn" v:"required#内容(英文)不能为空" dc:"内容(英文)"`
	ContentEs string `json:"contentEs" dc:"内容(西班牙语)"`
	ContentPt string `json:"contentPt" dc:"内容(葡萄牙语)"`
	ContentHi string `json:"contentHi" dc:"内容(印地语)"`
	UrlEn     string `json:"urlEn" v:"max-length:512#跳转链接(英文)最长512字符" dc:"跳转链接(英文)"`
	UrlEs     string `json:"urlEs" v:"max-length:512#跳转链接(西班牙语)最长512字符" dc:"跳转链接(西班牙语)"`
	UrlPt     string `json:"urlPt" v:"max-length:512#跳转链接(葡萄牙语)最长512字符" dc:"跳转链接(葡萄牙语)"`
	UrlHi     string `json:"urlHi" v:"max-length:512#跳转链接(印地语)最长512字符" dc:"跳转链接(印地语)"`
}

type UpdateActivityMessageRes struct {
	Success bool `json:"success" dc:"是否成功"`
}
