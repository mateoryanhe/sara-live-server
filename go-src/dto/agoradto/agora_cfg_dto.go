package agoradto

import "github.com/gogf/gf/v2/frame/g"

type GetAgoraCfgReq struct {
	g.Meta `path:"/getAgoraCfg" method:"post" summary:"查询声网配置" tags:"声网配置"`
}

type AgoraCfgItem struct {
	ID                 string `json:"id"`
	AppId              string `json:"appId"`
	AppCertificate     string `json:"appCertificate"`
	TokenExpireSeconds uint32 `json:"tokenExpireSeconds"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type GetAgoraCfgRes struct {
	Cfg *AgoraCfgItem `json:"cfg"`
}

type SaveAgoraCfgReq struct {
	g.Meta             `path:"/saveAgoraCfg" method:"post" summary:"保存声网配置" tags:"声网配置"`
	ID                 uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	AppId              string `json:"appId" v:"required|length:1,64#AppId不能为空|AppId长度需在1到64之间"`
	AppCertificate     string `json:"appCertificate" v:"required|length:1,128#AppCertificate不能为空|AppCertificate长度需在1到128之间"`
	TokenExpireSeconds uint32 `json:"tokenExpireSeconds" v:"required|min:60#TokenExpireSeconds不能为空|TokenExpireSeconds不能小于60"`
}

type SaveAgoraCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
