package domainsitedto

import "github.com/gogf/gf/v2/frame/g"

const (
	CMSSiteKey     = "cms"
	CMSURLPrefix   = "/cms"
	DefaultCMSRoot = "/home/ec2-user/cdn/cms"
)

type GetCMSDomainSiteMappingReq struct {
	g.Meta `path:"/getCMSDomainSiteMapping" method:"post" summary:"查询CMS域名映射" tags:"域名静态目录"`
}

type CMSDomainSiteMappingItem struct {
	ID        string `json:"id"        dc:"域名记录ID"`
	Domain    string `json:"domain"    dc:"CMS域名,只允许一个Host"`
	URLPrefix string `json:"urlPrefix" dc:"CMS静态访问前缀"`
	Root      string `json:"root"      dc:"CMS静态文件物理目录"`
	UpdatedAt string `json:"updatedAt" dc:"最近更新时间"`
}

type GetCMSDomainSiteMappingRes struct {
	Mapping *CMSDomainSiteMappingItem `json:"mapping"`
}

type SaveCMSDomainSiteMappingReq struct {
	g.Meta `path:"/saveCMSDomainSiteMapping" method:"post" summary:"保存CMS域名映射" tags:"域名静态目录"`
	Domain string `json:"domain" v:"required|length:1,253#CMS域名不能为空|CMS域名过长" dc:"CMS域名,只允许一个Host"`
	Root   string `json:"root" v:"required|length:1,1024#CMS目录不能为空|CMS目录过长" dc:"CMS静态文件绝对目录"`
}

type SaveCMSDomainSiteMappingRes struct {
	Success bool                      `json:"success"`
	Mapping *CMSDomainSiteMappingItem `json:"mapping"`
}
