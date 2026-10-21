package domainsitedto

import "github.com/gogf/gf/v2/frame/g"

type CMSDomainAliasItem struct {
	ID        string `json:"id"        dc:"记录ID"`
	Domain    string `json:"domain"    dc:"CMS备用域名"`
	Root      string `json:"root"      dc:"与主CMS映射相同的静态目录"`
	UpdatedAt string `json:"updatedAt" dc:"最近更新时间"`
}

type ListCMSDomainAliasReq struct {
	g.Meta `path:"/listCMSDomainAlias" method:"post" summary:"CMS备用域名列表" tags:"域名静态目录"`
}

type ListCMSDomainAliasRes struct {
	PrimaryDomain string                `json:"primaryDomain" dc:"主CMS域名(来自CMS域名映射)"`
	Root          string                `json:"root"          dc:"CMS静态目录"`
	Items         []*CMSDomainAliasItem `json:"items"`
}

type AddCMSDomainAliasReq struct {
	g.Meta `path:"/addCMSDomainAlias" method:"post" summary:"新增CMS备用域名" tags:"域名静态目录"`
	Domain string `json:"domain" v:"required|length:1,253#域名不能为空|域名过长" dc:"只允许一个Host"`
}

type AddCMSDomainAliasRes struct {
	Success bool                `json:"success"`
	Item    *CMSDomainAliasItem `json:"item"`
}

type DeleteCMSDomainAliasReq struct {
	g.Meta `path:"/deleteCMSDomainAlias" method:"post" summary:"删除CMS备用域名" tags:"域名静态目录"`
	ID     string `json:"id" v:"required#记录ID不能为空" dc:"listCMSDomainAlias 返回的 id"`
}

type DeleteCMSDomainAliasRes struct {
	Success bool `json:"success"`
}
