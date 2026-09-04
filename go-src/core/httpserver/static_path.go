package httpserver

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"xr-game-server/core/cfg"
)

func setupStaticPaths() {
	if len(cfg.GetStaticPathCfgs()) == 0 {
		bindCMSStaticFallback(context.Background())
		return
	}
	ctx := context.Background()
	httpServer.BindHookHandler("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
		serveConfiguredStaticPath(r)
	})
	g.Log().Warningf(ctx, "已启用多静态目录映射,共 %d 项,默认目录=%s", len(cfg.GetStaticPathCfgs()), cfg.GetServerRoot())
	bindCMSStaticFallback(ctx)
}

// setupStaticPageNoCache 启用静态入口页不缓存策略(见 writeStaticFile)
func setupStaticPageNoCache() {
	g.Log().Warning(context.Background(), "index.html/version.json 已禁用浏览器缓存")
}

func isNoCacheStaticPage(filePath string) bool {
	name := strings.ToLower(filepath.Base(filePath))
	return name == "index.html" || name == "version.json"
}

func applyNoCacheStaticPageHeaders(r *ghttp.Request) {
	if r == nil || r.Response == nil {
		return
	}
	r.Response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	r.Response.Header().Set("Pragma", "no-cache")
	r.Response.Header().Set("Expires", "0")
}

func serveConfiguredStaticPath(r *ghttp.Request) {
	if r == nil {
		return
	}
	if handleStaticCORS(r) {
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}
	root, rel, ok := cfg.MatchStaticPath(r.URL.Path)
	if !ok {
		return
	}
	realRoot := gfile.RealPath(root)
	if realRoot == "" {
		return
	}
	if strings.TrimSpace(rel) == "" {
		serveStaticIndexFile(r, realRoot)
		return
	}
	if strings.HasPrefix(normalizeURLPath(r.URL.Path), "/cms") {
		serveCMSStaticPath(r, realRoot, rel)
		return
	}
	filePath, exists := buildMappedStaticFilePath(realRoot, rel)
	if !exists || !gfile.Exists(filePath) {
		r.Response.WriteStatus(http.StatusNotFound)
		r.ExitAll()
		return
	}
	writeStaticFile(r, filePath)
	r.ExitAll()
}

func serveStaticIndexFile(r *ghttp.Request, root string) {
	indexPath := filepath.Join(root, "index.html")
	if !gfile.Exists(indexPath) {
		r.Response.WriteStatus(http.StatusNotFound)
		r.ExitAll()
		return
	}
	writeStaticFile(r, indexPath)
	r.ExitAll()
}

func serveCMSStaticPath(r *ghttp.Request, root, rel string) {
	filePath, exists := buildMappedStaticFilePath(root, rel)
	if exists && gfile.Exists(filePath) && !gfile.IsDir(filePath) {
		writeStaticFile(r, filePath)
		r.ExitAll()
		return
	}
	if isCMSAssetRequest(rel) {
		r.Response.WriteStatus(http.StatusNotFound)
		r.ExitAll()
		return
	}
	indexPath := filepath.Join(root, "index.html")
	if !gfile.Exists(indexPath) {
		r.Response.WriteStatus(http.StatusNotFound)
		r.ExitAll()
		return
	}
	writeStaticFile(r, indexPath)
	r.ExitAll()
}

func isCMSAssetRequest(rel string) bool {
	rel = strings.TrimSpace(strings.ReplaceAll(rel, "\\", "/"))
	if rel == "" {
		return false
	}
	if strings.HasPrefix(rel, "log-query-export/") {
		return true
	}
	return strings.Contains(filepath.Base(rel), ".")
}

func writeStaticFile(r *ghttp.Request, filePath string) {
	applyRequestCORS(r)
	if isNoCacheStaticPage(filePath) {
		applyNoCacheStaticPageHeaders(r)
	}
	if r.Method == http.MethodHead {
		r.Response.WriteHeader(http.StatusOK)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		r.Response.WriteStatus(http.StatusNotFound)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || info.IsDir() {
		r.Response.WriteStatus(http.StatusNotFound)
		return
	}
	r.Response.ServeContent(info.Name(), info.ModTime(), file)
}

func buildMappedStaticFilePath(root, rel string) (string, bool) {
	cleanRoot := filepath.Clean(root)
	target := cleanRoot
	if rel != "" {
		target = filepath.Clean(filepath.Join(cleanRoot, filepath.FromSlash(rel)))
	}
	if target != cleanRoot && !strings.HasPrefix(target, cleanRoot+string(filepath.Separator)) {
		return "", false
	}
	if gfile.IsDir(target) {
		target = filepath.Join(target, "index.html")
	}
	return target, true
}

func normalizeURLPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return strings.TrimRight(path, "/")
}
