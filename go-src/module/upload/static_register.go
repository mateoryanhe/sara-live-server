package upload

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cfg"
	"xr-game-server/core/httpserver"
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

	_, pathPrefix := parseResourceDomainBase(GetResourceDomain())
	if pathPrefix != "" {
		pathPrefix = "/" + pathPrefix
	} else if strings.EqualFold(filepath.Base(root), "images") {
		pathPrefix = "/images"
	}
	if pathPrefix != "" {
		// 路径映射与域名注册相互独立：关闭云桶且配置 IP 时，
		// 请求仍通过 http://IP/<prefix>/<file> 读取本地上传目录。
		cfg.RegisterRuntimeStaticPath(pathPrefix, root, GetCmsExportTtlMinutes())
		cfg.SetImageStaticPrefix(pathPrefix)
		g.Log().Infof(gctx.New(), "已注册 CMS 资源目录映射 prefix=%s root=%s", pathPrefix, root)
	}

	if !IsResourceDomainConfigured() {
		return
	}
	host := resourceDomainHost(GetResourceDomain())
	if host == "" {
		return
	}

	if !shouldRegisterResourceDomain(host, IsS3Enabled()) {
		g.Log().Infof(gctx.New(), "跳过 CMS 资源域名注册 domain=%s s3Enabled=%t pathPrefix=%s；IP 通过路径映射访问本地资源", host, IsS3Enabled(), pathPrefix)
		return
	}

	// 仅本地存储且配置为真实域名时注册；IP/IPv6/localhost 和云桶均不注册。
	cfg.RegisterRuntimeDomainSite(host, root)
	g.Log().Infof(gctx.New(), "已注册 CMS 资源域名静态目录 domain=%s root=%s", host, root)
}

// refreshStaticMappings 保存配置后重建完整域名快照，使目录和域名映射即时生效。
func refreshStaticMappings(ctx context.Context) error {
	registerStaticMappings()
	_, err := httpserver.RefreshDomainSiteRegistryFromConfig(ctx)
	return err
}

func shouldRegisterResourceDomain(host string, s3Enabled bool) bool {
	host = strings.TrimSpace(host)
	if s3Enabled || host == "" || strings.EqualFold(host, "localhost") {
		return false
	}
	return !isIPAddressHost(host)
}

func isIPAddressHost(host string) bool {
	host = strings.Trim(strings.TrimSpace(host), "[]")
	if zoneIndex := strings.LastIndexByte(host, '%'); zoneIndex > 0 {
		host = host[:zoneIndex]
	}
	return net.ParseIP(host) != nil
}
