package shareappdownload

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cfg"
	"xr-game-server/core/xrlog"
	"xr-game-server/dto/officialsitedeploydto"
)

// RefreshOfficialApkDownloadUrl 扫描官网部署目录 assets 下体积最大的 .apk,写入内存.
// 任意错误只打日志,不向上抛(启动/发版均不可因本逻辑失败).
func RefreshOfficialApkDownloadUrl() {
	defer func() {
		if rec := recover(); rec != nil {
			xrlog.DetailLog.Warningf(gctx.New(), "shareappdownload refresh apk url panic: %v", rec)
		}
	}()
	apkURL, err := resolveLargestOfficialApkURL()
	if err != nil {
		xrlog.DetailLog.Warningf(gctx.New(), "shareappdownload refresh apk url failed: %v", err)
		return
	}
	storeApkDownloadURL(apkURL)
	if apkURL != "" {
		xrlog.DetailLog.Infof(gctx.New(), "shareappdownload apk url refreshed: %s", apkURL)
	}
}

func resolveLargestOfficialApkURL() (string, error) {
	deployDir := strings.TrimSpace(cfg.GetStaticPathRoot(officialsitedeploydto.OfficialSiteStaticPrefix))
	if deployDir == "" {
		deployDir = officialsitedeploydto.DefaultOfficialSiteDeployPath
	}
	assetsDir := filepath.Join(deployDir, "assets")
	entries, err := os.ReadDir(assetsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	var (
		bestName string
		bestSize int64
	)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.EqualFold(filepath.Ext(name), ".apk") {
			continue
		}
		info, infoErr := entry.Info()
		if infoErr != nil {
			continue
		}
		size := info.Size()
		if size > bestSize || (size == bestSize && (bestName == "" || name > bestName)) {
			bestSize = size
			bestName = name
		}
	}
	if bestName == "" {
		return "", nil
	}
	domain := strings.TrimSpace(cfg.GetStaticSiteDomain(officialsitedeploydto.OfficialSiteStaticPrefix))
	if domain == "" {
		return "", nil
	}
	rel := filepath.ToSlash(filepath.Join("assets", bestName))
	return (&url.URL{
		Scheme: "https",
		Host:   domain,
		Path:   "/" + strings.TrimPrefix(rel, "/"),
	}).String(), nil
}
