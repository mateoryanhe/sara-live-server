package cfg

import (
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	CMSFileExportStaticPrefix      = "/cms-export"
	defaultCMSFileExportTTLMinutes = 30
)

var (
	cmsFileExportStaticPrefix string
	cmsFileExportRoot         string
	cmsFileExportTTLMinutes   int

	cmsExportStoragePathOverride func() string
	cmsExportTtlOverride         func() int
)

func RegisterCMSExportStoragePathOverride(fn func() string) {
	cmsExportStoragePathOverride = fn
}

func RegisterCMSExportTtlOverride(fn func() int) {
	cmsExportTtlOverride = fn
}

func initCMSFileExportCfg() {
	cmsFileExportStaticPrefix = CMSFileExportStaticPrefix
	cmsFileExportRoot = ""
	cmsFileExportTTLMinutes = defaultCMSFileExportTTLMinutes

	for _, item := range GetStaticPathCfgs() {
		if item == nil || item.Prefix != CMSFileExportStaticPrefix {
			continue
		}
		cmsFileExportRoot = strings.TrimSpace(item.Path)
		if item.TTLMinutes > 0 {
			cmsFileExportTTLMinutes = item.TTLMinutes
		}
		break
	}

	ctx := gctx.New()
	if cmsFileExportRoot == "" {
		g.Log().Warningf(ctx, "server.staticSites 未配置 %s,CMS 文件导出将不可用", CMSFileExportStaticPrefix)
		return
	}
	g.Log().Warningf(ctx, "CMS文件导出已加载 staticPrefix=%s root=%s ttlMinutes=%d",
		cmsFileExportStaticPrefix, cmsFileExportRoot, cmsFileExportTTLMinutes)
}

// GetCMSFileExportStaticPrefix CMS 文件导出下载 URL 前缀
func GetCMSFileExportStaticPrefix() string {
	if cmsFileExportStaticPrefix == "" {
		return CMSFileExportStaticPrefix
	}
	return cmsFileExportStaticPrefix
}

// GetCMSFileExportRoot CMS 文件导出物理目录
func GetCMSFileExportRoot() string {
	if cmsExportStoragePathOverride != nil {
		if p := strings.TrimSpace(cmsExportStoragePathOverride()); p != "" {
			return p
		}
	}
	return cmsFileExportRoot
}

// GetCMSFileExportTTLMinutes 导出文件过期清理时间(分钟)
func GetCMSFileExportTTLMinutes() int {
	if cmsExportTtlOverride != nil {
		if ttl := cmsExportTtlOverride(); ttl > 0 {
			return ttl
		}
	}
	if cmsFileExportTTLMinutes <= 0 {
		return defaultCMSFileExportTTLMinutes
	}
	return cmsFileExportTTLMinutes
}

// ResolveCMSFileExportDir CMS 文件导出目录(与上传资源共用 storagePath)
func ResolveCMSFileExportDir() string {
	root := GetCMSFileExportRoot()
	if root == "" {
		return "."
	}
	return root
}

// BuildCMSFileExportURLPrefix CMS 文件导出 URL 前缀
func BuildCMSFileExportURLPrefix() string {
	return GetCMSFileExportStaticPrefix()
}

// BuildCMSFileExportURL 构建文件下载 URL,如 /cms-export/{fileName}
func BuildCMSFileExportURL(fileName string) string {
	prefix := BuildCMSFileExportURLPrefix()
	fileName = strings.Trim(strings.ReplaceAll(fileName, "\\", "/"), "/")
	if fileName == "" {
		return prefix
	}
	return prefix + "/" + fileName
}

// JoinCMSFileExportPath 拼接导出目录下的相对路径(禁止 ..)
func JoinCMSFileExportPath(name string) string {
	name = strings.Trim(strings.ReplaceAll(name, "\\", "/"), "/")
	if name == "" || strings.Contains(name, "..") {
		return ResolveCMSFileExportDir()
	}
	return filepath.Join(ResolveCMSFileExportDir(), name)
}
