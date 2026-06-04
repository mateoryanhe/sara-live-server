package aliyuntextmoderationdto

import "github.com/gogf/gf/v2/frame/g"

type GetCfgReq struct {
	g.Meta `path:"/getTextModerationCfg" method:"post" summary:"查询阿里云文本审核配置" tags:"文本审核配置"`
}

type CfgItem struct {
	ID              string `json:"id"`
	Enabled         bool   `json:"enabled"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	RegionId        string `json:"regionId"`
	Endpoint        string `json:"endpoint"`
	ChatService     string `json:"chatService"`
	NicknameService string `json:"nicknameService"`
	CommentService  string `json:"commentService"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

type GetCfgRes struct {
	Cfg *CfgItem `json:"cfg"`
}

type SaveCfgReq struct {
	g.Meta          `path:"/saveTextModerationCfg" method:"post" summary:"保存阿里云文本审核配置" tags:"文本审核配置"`
	ID              uint64 `json:"id"`
	Enabled         bool   `json:"enabled"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	RegionId        string `json:"regionId"`
	Endpoint        string `json:"endpoint"`
	ChatService     string `json:"chatService"`
	NicknameService string `json:"nicknameService"`
	CommentService  string `json:"commentService"`
}

type SaveCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
