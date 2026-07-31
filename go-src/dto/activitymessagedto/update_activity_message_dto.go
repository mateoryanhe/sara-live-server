package activitymessagedto

import "github.com/gogf/gf/v2/frame/g"

type UpdateActivityMessageReq struct {
	g.Meta    `path:"/updateActivityMessage" method:"post" summary:"修改活动消息" tags:"活动消息"`
	ID        uint64 `json:"id" v:"required#活动消息ID不能为空"`
	IconEn    string `json:"iconEn" v:"max-length:255#图标资源名(英文)最长255字符"`
	IconEs    string `json:"iconEs" v:"max-length:255#图标资源名(西班牙语)最长255字符"`
	IconPt    string `json:"iconPt" v:"max-length:255#图标资源名(葡萄牙语)最长255字符"`
	IconHi    string `json:"iconHi" v:"max-length:255#图标资源名(印地语)最长255字符"`
	BgEn      string `json:"bgEn" v:"max-length:255#背景图资源名(英文)最长255字符"`
	BgEs      string `json:"bgEs" v:"max-length:255#背景图资源名(西班牙语)最长255字符"`
	BgPt      string `json:"bgPt" v:"max-length:255#背景图资源名(葡萄牙语)最长255字符"`
	BgHi      string `json:"bgHi" v:"max-length:255#背景图资源名(印地语)最长255字符"`
	TitleEn   string `json:"titleEn" v:"required|max-length:128#标题(英文)不能为空|标题(英文)最长128字符"`
	TitleEs   string `json:"titleEs" v:"max-length:128#标题(西班牙语)最长128字符"`
	TitlePt   string `json:"titlePt" v:"max-length:128#标题(葡萄牙语)最长128字符"`
	TitleHi   string `json:"titleHi" v:"max-length:128#标题(印地语)最长128字符"`
	ContentEn string `json:"contentEn" v:"required#内容(英文)不能为空"`
	ContentEs string `json:"contentEs"`
	ContentPt string `json:"contentPt"`
	ContentHi string `json:"contentHi"`
	UrlEn     string `json:"urlEn" v:"max-length:512#跳转链接(英文)最长512字符"`
	UrlEs     string `json:"urlEs" v:"max-length:512#跳转链接(西班牙语)最长512字符"`
	UrlPt     string `json:"urlPt" v:"max-length:512#跳转链接(葡萄牙语)最长512字符"`
	UrlHi     string `json:"urlHi" v:"max-length:512#跳转链接(印地语)最长512字符"`
}

type UpdateActivityMessageRes struct {
	Success bool `json:"success"`
}
