package cfg

import (
	"context"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

// StaticPathCfg URL 前缀与物理目录映射
type StaticPathCfg struct {
	Prefix     string `json:"prefix"     dc:"URL前缀,如 /cms"`
	Path       string `json:"path"       dc:"物理目录绝对路径(GoFrame 原生字段)"`
	Root       string `json:"root"       dc:"物理目录,兼容旧配置"`
	TTLMinutes int    `json:"ttlMinutes" dc:"CMS文件导出过期清理(分钟),仅 /cms-export 使用"`
}

func staticPathPhysicalDir(item *StaticPathCfg) string {
	if item == nil {
		return ""
	}
	if p := strings.TrimSpace(item.Path); p != "" {
		return p
	}
	return strings.TrimSpace(item.Root)
}

var staticPathCfgs []*StaticPathCfg
var imageStaticPrefix string

// GetImageStaticPrefix 图片静态 URL 前缀,如 /images
func GetImageStaticPrefix() string {
	return imageStaticPrefix
}

// GetImageStaticRoot 图片物理目录,来自 staticSites 映射
func GetImageStaticRoot() string {
	if imageStaticPrefix == "" {
		return ""
	}
	return GetStaticPathRoot(imageStaticPrefix)
}

// GetImageStaticPathSegment 图片 URL 路径段,如 images
func GetImageStaticPathSegment() string {
	return strings.Trim(imageStaticPrefix, "/")
}

// GetStaticPathCfgs 获取静态路径映射(按前缀长度降序,优先最长匹配)
func GetStaticPathCfgs() []*StaticPathCfg {
	return staticPathCfgs
}

// GetServerRoot 默认静态资源根目录
func GetServerRoot() string {
	return strings.TrimSpace(g.Cfg().MustGet(gctx.New(), "server.serverRoot").String())
}

func initImageStaticCfg() {
	ctx := gctx.New()
	prefix := resolveImageStaticPrefix(ctx)
	if prefix == "" {
		g.Log().Error(ctx, "未找到图片静态站点,请在 server.staticSites 中配置 images 域名或 prefix:/images")
		return
	}
	root := strings.TrimSpace(GetStaticPathRoot(prefix))
	if root == "" {
		g.Log().Errorf(ctx, "图片静态 prefix=%s 在 server.staticSites 中未找到对应 path", prefix)
		return
	}
	if gfile.RealPath(root) == "" {
		g.Log().Errorf(ctx, "图片静态目录不存在 prefix=%s path=%s", prefix, root)
	}
	imageStaticPrefix = prefix
	g.Log().Warningf(ctx, "图片静态路径已加载 prefix=%s path=%s", prefix, root)
}

func resolveImageStaticPrefix(ctx context.Context) string {
	explicit := normalizeURLPrefix(g.Cfg().MustGet(ctx, "server.imageStaticPrefix").String())
	if explicit != "" && explicit != "/" {
		return explicit
	}
	for _, item := range staticPathCfgs {
		if item == nil {
			continue
		}
		if strings.EqualFold(strings.Trim(item.Prefix, "/"), "images") {
			return item.Prefix
		}
	}
	return ""
}

func normalizeStaticPathCfgs(list []*StaticPathCfg) []*StaticPathCfg {
	ret := make([]*StaticPathCfg, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		item.Prefix = normalizeURLPrefix(item.Prefix)
		dir := staticPathPhysicalDir(item)
		if item.Prefix == "" || dir == "" {
			continue
		}
		item.Path = dir
		item.Root = dir
		ret = append(ret, item)
	}
	sort.Slice(ret, func(i, j int) bool {
		return len(ret[i].Prefix) > len(ret[j].Prefix)
	})
	return ret
}

func normalizeURLPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" || prefix == "/" {
		return "/"
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	return strings.TrimRight(prefix, "/")
}

// MatchStaticPath 按 URL 匹配静态目录,返回物理根目录与相对路径
func MatchStaticPath(urlPath string) (root string, rel string, ok bool) {
	urlPath = normalizeURLPrefix(urlPath)
	if urlPath == "/" {
		urlPath = "/"
	}
	for _, item := range staticPathCfgs {
		if item == nil {
			continue
		}
		if item.Prefix == "/" {
			if urlPath == "/" || strings.HasPrefix(urlPath, "/") {
				return item.Path, strings.TrimPrefix(urlPath, "/"), true
			}
			continue
		}
		if urlPath == item.Prefix || strings.HasPrefix(urlPath, item.Prefix+"/") {
			rel = strings.TrimPrefix(urlPath, item.Prefix)
			rel = strings.TrimPrefix(rel, "/")
			return item.Path, rel, true
		}
	}
	return "", "", false
}

// GetStaticPathRoot 获取指定前缀的物理目录,未配置时返回空
func GetStaticPathRoot(prefix string) string {
	prefix = normalizeURLPrefix(prefix)
	for _, item := range staticPathCfgs {
		if item != nil && item.Prefix == prefix {
			return item.Path
		}
	}
	return ""
}

// ResolvePhysicalDir 解析 URL 前缀或相对路径对应的物理目录
func ResolvePhysicalDir(prefixOrSubDir string) string {
	prefix := normalizeURLPrefix(prefixOrSubDir)
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + strings.Trim(prefix, "/")
	}
	if root := GetStaticPathRoot(prefix); root != "" {
		return root
	}
	base := GetServerRoot()
	if base == "" {
		base = "."
	}
	sub := strings.Trim(prefixOrSubDir, "/")
	if sub == "" {
		return base
	}
	if strings.HasPrefix(prefixOrSubDir, "/") {
		return filepath.Join(base, strings.TrimPrefix(prefix, "/"))
	}
	return filepath.Join(base, sub)
}
