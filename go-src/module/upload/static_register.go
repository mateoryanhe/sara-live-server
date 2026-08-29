package upload

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cfg"
)

func registerStaticMappings() {
	root := GetStoragePath()
	cfg.ClearRuntimeStaticMappings()
	cfg.SetImageStaticPrefix("")

	if root == "" {
		return
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		g.Log().Warningf(gctx.New(), "创建 CMS 资源存储目录失败 root=%s err=%v", root, err)
	}

	if !IsResourceDomainConfigured() {
		return
	}
	host := resourceDomainHost(GetResourceDomain())
	if host == "" {
		return
	}

	// 本地 IP(127.0.0.1/localhost): 由 config.yaml serverRoot 提供静态访问
	if isLoopbackHost(host) {
		g.Log().Warningf(gctx.New(), "本地 IP %s 使用 serverRoot 静态目录,跳过资源域名 hook", host)
		return
	}

	// 正式服资源域名: 注册图片/导出等静态文件域名站点(与 URL 生成共用 CMS 资源域名)
	cfg.RegisterRuntimeDomainSite(host, root)
	_, pathPrefix := parseResourceDomainBase(GetResourceDomain())
	if pathPrefix != "" {
		cfg.SetImageStaticPrefix("/" + pathPrefix)
	} else if strings.EqualFold(filepath.Base(root), "images") {
		cfg.SetImageStaticPrefix("/images")
	}
	g.Log().Warningf(gctx.New(), "已注册 CMS 资源域名静态目录 domain=%s root=%s", host, root)
}
