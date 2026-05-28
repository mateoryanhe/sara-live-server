package cfg

import (
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// DomainSiteCfg 域名静态站点配置
type DomainSiteCfg struct {
	Domain   string `json:"domain"   dc:"域名,多个用逗号分隔"`
	Root     string `json:"root"     dc:"静态资源根目录"`
	CertFile string `json:"certFile" dc:"HTTPS证书文件路径"`
	KeyFile  string `json:"keyFile"  dc:"HTTPS私钥文件路径"`
}

var domainSiteCfgs []*DomainSiteCfg

// GetDomainSiteCfgs 获取域名站点配置
func GetDomainSiteCfgs() []*DomainSiteCfg {
	return domainSiteCfgs
}

func initDomainSiteCfg() {
	ctx := gctx.New()
	list := make([]*DomainSiteCfg, 0)
	_ = g.Cfg().MustGet(ctx, "server.domainSites").Scan(&list)
	domainSiteCfgs = normalizeDomainSiteCfgs(list)
	if len(domainSiteCfgs) > 0 {
		cfgJson, _ := json.MarshalIndent(domainSiteCfgs, "", " ")
		g.Log().Warningf(ctx, "成功加载域名站点配置:%s", cfgJson)
	}
}

func normalizeDomainSiteCfgs(list []*DomainSiteCfg) []*DomainSiteCfg {
	ret := make([]*DomainSiteCfg, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		item.Domain = strings.TrimSpace(item.Domain)
		item.Root = strings.TrimSpace(item.Root)
		item.CertFile = strings.TrimSpace(item.CertFile)
		item.KeyFile = strings.TrimSpace(item.KeyFile)
		if item.Domain == "" || item.Root == "" {
			continue
		}
		ret = append(ret, item)
	}
	return ret
}

// SplitDomains 解析逗号分隔域名
func SplitDomains(domain string) []string {
	parts := strings.Split(domain, ",")
	ret := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		ret = append(ret, part)
	}
	return ret
}
