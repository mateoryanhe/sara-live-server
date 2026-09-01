package cfg

import (
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const defaultCMSFileExportTTLMinutes = 30

var (
	cmsFileExportRoot         string
	cmsFileExportTTLMinutes   int
	cmsExportStoragePathOverride func() string
	cmsExportTtlOverride         func() int
	cmsExportURLBuilder          func(fileName string) string
	cmsExportURLPrefixProvider   func() string
)

func RegisterCMSExportStoragePathOverride(fn func() string) {
	cmsExportStoragePathOverride = fn
}

func RegisterCMSExportTtlOverride(fn func() int) {
	cmsExportTtlOverride = fn
}

// RegisterCMSExportURLBuilder 注册导出文件 URL 构建(与上传资源共用,由 upload 模块提供)
func RegisterCMSExportURLBuilder(fn func(fileName string) string) {
	cmsExportURLBuilder = fn
}

// RegisterCMSExportURLPrefixProvider 注册导出 URL 根路径展示(一般为 CMS 资源域名)
func RegisterCMSExportURLPrefixProvider(fn func() string) {
	cmsExportURLPrefixProvider = fn
}

func initCMSFileExportCfg() {
	cmsFileExportRoot = ""
	cmsFileExportTTLMinutes = defaultCMSFileExportTTLMinutes

	ctx := gctx.New()
	if cmsExportStoragePathOverride == nil {
		g.Log().Warningf(ctx, "CMS 文件导出目录尚未配置,等待 upload 模块加载 storagePath")
		return
	}
	root := GetCMSFileExportRoot()
	if root == "" {
		g.Log().Warningf(ctx, "CMS 文件导出目录为空")
		return
	}
	g.Log().Warningf(ctx, "CMS文件导出已加载 root=%s urlPrefix=%s ttlMinutes=%d",
		root, BuildCMSFileExportURLPrefix(), GetCMSFileExportTTLMinutes())
}

// GetCMSFileExportStaticPrefix 兼容旧字段,返回资源访问根 URL
func GetCMSFileExportStaticPrefix() string {
	return BuildCMSFileExportURLPrefix()
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

// BuildCMSFileExportURLPrefix CMS 资源访问根 URL(用于日志页展示)
func BuildCMSFileExportURLPrefix() string {
	if cmsExportURLPrefixProvider != nil {
		return strings.TrimRight(strings.TrimSpace(cmsExportURLPrefixProvider()), "/")
	}
	return ""
}

// BuildCMSFileExportURL 构建导出文件下载 URL(与上传资源 URL 规则一致)
func BuildCMSFileExportURL(fileName string) string {
	fileName = strings.Trim(strings.ReplaceAll(fileName, "\\", "/"), "/")
	if fileName == "" {
		return BuildCMSFileExportURLPrefix()
	}
	if cmsExportURLBuilder != nil {
		return cmsExportURLBuilder(fileName)
	}
	prefix := BuildCMSFileExportURLPrefix()
	if prefix == "" {
		return "/" + fileName
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
