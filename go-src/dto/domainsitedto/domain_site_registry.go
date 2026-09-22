package domainsitedto

import "github.com/gogf/gf/v2/frame/g"

type DomainSiteItem struct {
	Domain   string `json:"domain"   dc:"域名,每项只填写一个Host"`
	Root     string `json:"root"     dc:"静态资源根目录"`
	CertFile string `json:"certFile" dc:"HTTPS证书文件路径,可选"`
	KeyFile  string `json:"keyFile"  dc:"HTTPS私钥文件路径,可选"`
}

type RefreshDomainSiteRegistryReq struct {
	g.Meta `path:"/refreshDomainSiteRegistry" method:"post" summary:"刷新域名静态目录注册表" tags:"域名静态目录"`
	// Sites 省略或传 null 时从 server.staticSites 与运行时注册配置重建；传数组时原子替换当前进程注册表。
	Sites []*DomainSiteItem `json:"sites" dc:"完整域名映射数组;传空数组可清空当前注册表"`
}

type RefreshDomainSiteRegistryRes struct {
	Count  int               `json:"count"  dc:"当前已注册域名数量"`
	Source string            `json:"source" dc:"刷新来源: config/request"`
	Sites  []*DomainSiteItem `json:"sites"  dc:"当前进程生效的域名映射"`
}
