package domainsitedto

import "github.com/gogf/gf/v2/frame/g"

type ExtraDomainMappingItem struct {
	ID           string `json:"id"`
	Domain       string `json:"domain"`
	Root         string `json:"root"`
	HasTLS       bool   `json:"hasTls" dc:"是否已配置独立证书(上传的橙云源站证书)"`
	CertFile     string `json:"certFile" dc:"证书绝对路径,未上传则为空"`
	KeyFile      string `json:"keyFile" dc:"私钥绝对路径,未上传则为空"`
	UpdatedAt    string `json:"updatedAt"`
}

type ListExtraDomainMappingsReq struct {
	g.Meta `path:"/listExtraDomainMappings" method:"post" summary:"额外域名目录映射列表" tags:"域名静态目录"`
}

type ListExtraDomainMappingsRes struct {
	Items []*ExtraDomainMappingItem `json:"items"`
}

type SaveExtraDomainMappingReq struct {
	g.Meta `path:"/saveExtraDomainMapping" method:"post" summary:"保存额外域名目录映射" tags:"域名静态目录"`
	Domain string `json:"domain" v:"required|length:1,253#域名不能为空|域名过长" dc:"Host,不含协议和路径"`
	Root   string `json:"root" v:"required|length:1,1024#目录不能为空|目录过长" dc:"静态文件绝对目录"`
}

type SaveExtraDomainMappingRes struct {
	Success bool                    `json:"success"`
	Item    *ExtraDomainMappingItem `json:"item"`
}

type DeleteExtraDomainMappingReq struct {
	g.Meta `path:"/deleteExtraDomainMapping" method:"post" summary:"删除额外域名目录映射" tags:"域名静态目录"`
	ID     string `json:"id" v:"required#记录ID不能为空"`
}

type DeleteExtraDomainMappingRes struct {
	Success bool `json:"success"`
}
