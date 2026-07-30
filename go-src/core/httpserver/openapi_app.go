package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/text/gstr"
)

const openApiDefaultMethod = "ALL"

var openApiHTTPMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
	http.MethodPatch,
	http.MethodHead,
	http.MethodOptions,
	http.MethodTrace,
	http.MethodConnect,
}

func setupAppOpenApiHook() {
	ctx := context.Background()
	openapiPath := g.Cfg().MustGet(ctx, "server.openapiPath").String()
	if openapiPath == "" {
		return
	}
	// 导出全部时走 GoFrame 默认 handler,不做 CMS 过滤
	if isOpenApiExportAll() {
		return
	}
	httpServer.BindHookHandler(openapiPath, ghttp.HookBeforeServe, hookServeAppOpenApi)
}

func isOpenApiExportAll() bool {
	return g.Cfg().MustGet(context.Background(), "server.openapiExportAll", false).Bool()
}

func hookServeAppOpenApi(r *ghttp.Request) {
	raw, err := buildAppOpenApiJSON()
	if err != nil {
		g.Log().Fatalf(r.Context(), "build app openapi failed: %+v", err)
	}
	r.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.Response.Write(raw)
	r.ExitAll()
}

// buildAppOpenApiJSON 基于 GoFrame 完整 OpenAPI 过滤 CMS 路由,保留 App 接口与 /push 推送文档.
func buildAppOpenApiJSON() ([]byte, error) {
	base := httpServer.GetOpenApi()
	if base == nil || len(base.Paths) == 0 {
		return []byte(`{"openapi":"3.0.0","info":{"title":"","version":""},"paths":{}}`), nil
	}
	cmsSet := buildCMSPathMethodSet()
	return filterAppOpenApiJSON(base, cmsSet)
}

func buildCMSPathMethodSet() map[string]struct{} {
	set := make(map[string]struct{})
	for _, item := range httpServer.GetRoutes() {
		if !isCMSOpenApiRoute(item) {
			continue
		}
		if item.Handler == nil || !item.Handler.Info.IsStrictRoute {
			continue
		}
		methods := []string{item.Method}
		if gstr.Equal(item.Method, openApiDefaultMethod) {
			methods = ghttp.SupportedMethods()
		}
		for _, method := range methods {
			set[openApiPathMethodKey(item.Route, method)] = struct{}{}
		}
	}
	return set
}

func filterAppOpenApiJSON(base *goai.OpenApiV3, cmsSet map[string]struct{}) ([]byte, error) {
	raw := []byte(base.String())
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	paths, _ := doc["paths"].(map[string]any)
	if paths == nil {
		return raw, nil
	}
	for route, pathVal := range paths {
		pathMap, ok := pathVal.(map[string]any)
		if !ok {
			continue
		}
		for _, method := range openApiHTTPMethods {
			if _, isCMS := cmsSet[openApiPathMethodKey(route, method)]; isCMS {
				delete(pathMap, strings.ToLower(method))
			}
		}
		if len(pathMap) == 0 {
			delete(paths, route)
			continue
		}
		paths[route] = pathMap
	}
	return json.Marshal(doc)
}

func openApiPathMethodKey(route, method string) string {
	return gstr.ToUpper(method) + " " + route
}

// isCMSOpenApiRoute 判断是否为 CMS 接口(用于 OpenAPI 导出过滤).
func isCMSOpenApiRoute(item ghttp.RouterItem) bool {
	if strings.Contains(item.Middleware, "MiddlewareCmsAuth") {
		return true
	}
	if item.Handler != nil && item.Handler.Name != "" {
		if idx := strings.LastIndex(item.Handler.Name, "."); idx >= 0 {
			if strings.HasPrefix(item.Handler.Name[idx+1:], "CMS") {
				return true
			}
		}
	}
	return false
}
