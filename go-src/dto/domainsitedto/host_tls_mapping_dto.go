package domainsitedto

import "github.com/gogf/gf/v2/frame/g"

type HostTlsMappingItem struct {
	ID        string `json:"id"`
	Domain    string `json:"domain"`
	HasTLS    bool   `json:"hasTls"`
	CertFile  string `json:"certFile"`
	KeyFile   string `json:"keyFile"`
	UpdatedAt string `json:"updatedAt"`
}

type ListHostTlsMappingsReq struct {
	g.Meta `path:"/listHostTlsMappings" method:"post" summary:"纯Host证书映射列表" tags:"Host证书映射"`
}

type ListHostTlsMappingsRes struct {
	Items []*HostTlsMappingItem `json:"items"`
}

type SaveHostTlsMappingRes struct {
	Success bool                `json:"success"`
	Item    *HostTlsMappingItem `json:"item"`
}

type DeleteHostTlsMappingReq struct {
	g.Meta `path:"/deleteHostTlsMapping" method:"post" summary:"删除纯Host证书映射" tags:"Host证书映射"`
	ID     string `json:"id" v:"required#记录ID不能为空" dc:"记录ID"`
}

type DeleteHostTlsMappingRes struct {
	Success bool `json:"success"`
}
