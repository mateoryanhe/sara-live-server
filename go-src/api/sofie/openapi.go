package sofie

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
)

// Sofie 文档从旧接口 schema 克隆时,部分路径需改 JSON 字段名(避免 game 等敏感词).
var sofieOpenAPIJSONFieldRenames = map[string]map[string]string{
	"/sofie/play/browseTitles": {
		"gameCode": "titleCode",
	},
	"/sofie/play/openTitle": {
		"gameCode": "titleCode",
	},
}

func renameOpenAPIJSONFields(v any, renames map[string]string) {
	if len(renames) == 0 || v == nil {
		return
	}
	switch node := v.(type) {
	case map[string]any:
		for oldKey, newKey := range renames {
			if oldKey == newKey {
				continue
			}
			if val, ok := node[oldKey]; ok {
				delete(node, oldKey)
				node[newKey] = val
			}
		}
		if req, ok := node["required"].([]any); ok {
			for i, item := range req {
				s, ok := item.(string)
				if !ok {
					continue
				}
				if neu, hit := renames[s]; hit {
					req[i] = neu
				}
			}
		}
		for _, child := range node {
			renameOpenAPIJSONFields(child, renames)
		}
	case []any:
		for _, child := range node {
			renameOpenAPIJSONFields(child, renames)
		}
	}
}

func scrubSofieOpenAPIText(s string) string {
	replacer := strings.NewReplacer(
		"玩游戏", "",
		"游戏编码", "标题编码",
		"游戏启动链接", "启动链接",
		"启动游戏", "打开标题",
		"游戏列表", "标题列表",
		"游戏", "标题",
		"gameId", "titleId",
		"GameId", "TitleId",
	)
	return strings.TrimSpace(replacer.Replace(s))
}

func scrubSofieOpenAPIWording(v any) {
	if v == nil {
		return
	}
	switch node := v.(type) {
	case map[string]any:
		for k, child := range node {
			if k == "description" || k == "summary" || k == "title" {
				if s, ok := child.(string); ok && s != "" {
					node[k] = scrubSofieOpenAPIText(s)
				}
				continue
			}
			scrubSofieOpenAPIWording(child)
		}
	case []any:
		for _, child := range node {
			scrubSofieOpenAPIWording(child)
		}
	}
}

func sofieCloneSchemaName(newPath, oldSchemaName string) string {
	base := strings.Trim(strings.ReplaceAll(newPath, "/", "_"), "_")
	short := oldSchemaName
	if i := strings.LastIndex(oldSchemaName, "."); i >= 0 {
		short = oldSchemaName[i+1:]
	}
	short = strings.ReplaceAll(short, "Game", "Title")
	short = strings.ReplaceAll(short, "game", "title")
	return base + "_" + short
}

// rewriteOpenAPIRefs 把 operation 内 $ref 克隆为 Sofie 专用 schema,再按需改字段名/文案.
// GoFrame 文档里 requestBody 通常是 $ref 到 components,只改 operation 树上的 key 不会动到真正字段.
func rewriteOpenAPIRefs(node any, schemas map[string]any, newPath string, renames map[string]string, cloned map[string]string) {
	if node == nil || schemas == nil {
		return
	}
	switch n := node.(type) {
	case map[string]any:
		if ref, ok := n["$ref"].(string); ok {
			const prefix = "#/components/schemas/"
			if !strings.HasPrefix(ref, prefix) {
				return
			}
			oldName := strings.TrimPrefix(ref, prefix)
			if neu, done := cloned[oldName]; done {
				n["$ref"] = prefix + neu
				return
			}
			src, ok := schemas[oldName]
			if !ok {
				return
			}
			newName := sofieCloneSchemaName(newPath, oldName)
			clonedSch := cloneJSONValue(src)
			if len(renames) > 0 {
				renameOpenAPIJSONFields(clonedSch, renames)
				scrubSofieOpenAPIWording(clonedSch)
			}
			schemas[newName] = clonedSch
			cloned[oldName] = newName
			n["$ref"] = prefix + newName
			rewriteOpenAPIRefs(clonedSch, schemas, newPath, renames, cloned)
			return
		}
		for _, child := range n {
			rewriteOpenAPIRefs(child, schemas, newPath, renames, cloned)
		}
	case []any:
		for _, child := range n {
			rewriteOpenAPIRefs(child, schemas, newPath, renames, cloned)
		}
	}
}

func collectOpenAPISchemaRefs(node any, out map[string]struct{}) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case map[string]any:
		if ref, ok := n["$ref"].(string); ok {
			const prefix = "#/components/schemas/"
			if strings.HasPrefix(ref, prefix) {
				out[strings.TrimPrefix(ref, prefix)] = struct{}{}
			}
		}
		for _, child := range n {
			collectOpenAPISchemaRefs(child, out)
		}
	case []any:
		for _, child := range n {
			collectOpenAPISchemaRefs(child, out)
		}
	}
}

func pruneSofieOpenAPIComponents(doc map[string]any) {
	components, _ := doc["components"].(map[string]any)
	if components == nil {
		return
	}
	schemas, _ := components["schemas"].(map[string]any)
	if schemas == nil {
		return
	}
	needed := map[string]struct{}{}
	collectOpenAPISchemaRefs(doc["paths"], needed)
	for changed := true; changed; {
		changed = false
		for name := range needed {
			sch, ok := schemas[name]
			if !ok {
				continue
			}
			before := len(needed)
			collectOpenAPISchemaRefs(sch, needed)
			if len(needed) > before {
				changed = true
			}
		}
	}
	kept := make(map[string]any, len(needed))
	for name := range needed {
		if sch, ok := schemas[name]; ok {
			kept[name] = sch
		}
	}
	components["schemas"] = kept
}

func setupSofieOpenApi() {
	ctx := context.Background()
	openapiPath := strings.TrimSpace(g.Cfg().MustGet(ctx, "server.sofieOpenapiPath").String())
	if openapiPath == "" {
		return
	}
	server := g.Server()
	server.BindHookHandler(openapiPath, ghttp.HookBeforeServe, func(r *ghttp.Request) {
		raw, err := buildSofieOpenApiJSON()
		if err != nil {
			g.Log().Errorf(r.Context(), "build sofie openapi failed: %+v", err)
			r.Response.WriteStatus(500)
			r.ExitAll()
			return
		}
		r.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
		r.Response.Write(raw)
		r.ExitAll()
	})
	swaggerPath := strings.TrimSpace(g.Cfg().MustGet(ctx, "server.sofieSwaggerPath").String())
	if swaggerPath == "" {
		return
	}
	server.BindHookHandler(swaggerPath, ghttp.HookBeforeServe, func(r *ghttp.Request) {
		r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
		r.Response.Write(sofieSwaggerHTML(openapiPath))
		r.ExitAll()
	})
}

func buildSofieOpenApiJSON() ([]byte, error) {
	base := g.Server().GetOpenApi()
	if base == nil {
		return []byte(`{"openapi":"3.0.0","info":{"title":"Sofie App API","version":"1.0.0"},"paths":{}}`), nil
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(base.String()), &doc); err != nil {
		return nil, err
	}
	if err := httpserver.MergeOpenApiRawHandlerRoutes(doc); err != nil {
		return nil, err
	}
	oldPaths, _ := doc["paths"].(map[string]any)
	if oldPaths == nil {
		oldPaths = map[string]any{}
	}
	components, _ := doc["components"].(map[string]any)
	if components == nil {
		components = map[string]any{}
		doc["components"] = components
	}
	schemas, _ := components["schemas"].(map[string]any)
	if schemas == nil {
		schemas = map[string]any{}
		components["schemas"] = schemas
	}
	newPaths := make(map[string]any, len(openAPIDocRoutes))
	for _, route := range openAPIDocRoutes {
		op := findOpenApiOperation(oldPaths, route.OldPath, route.OldMethod)
		if op == nil {
			continue
		}
		cloned := cloneJSONValue(op)
		if opMap, ok := cloned.(map[string]any); ok {
			opMap["tags"] = []any{route.Tag}
			if route.Summary != "" {
				opMap["summary"] = route.Summary
			}
			opMap["operationId"] = strings.Trim(strings.ReplaceAll(route.NewPath, "/", "_"), "_")
			renames := sofieOpenAPIJSONFieldRenames[route.NewPath]
			if len(renames) > 0 {
				rewriteOpenAPIRefs(opMap, schemas, route.NewPath, renames, map[string]string{})
			}
		}
		pathItem, _ := newPaths[route.NewPath].(map[string]any)
		if pathItem == nil {
			pathItem = map[string]any{}
		}
		pathItem["post"] = cloned
		newPaths[route.NewPath] = pathItem
	}
	doc["paths"] = newPaths
	pruneSofieOpenAPIComponents(doc)
	info, _ := doc["info"].(map[string]any)
	if info == nil {
		info = map[string]any{}
		doc["info"] = info
	}
	info["title"] = "Sofie App API"
	info["description"] = "Sofie 独立 App 接口文档，与旧 App URL 并存且不混用。"
	info["version"] = "1.0.0"
	return json.Marshal(doc)
}

func findOpenApiOperation(paths map[string]any, oldPath, oldMethod string) any {
	pathItem, _ := paths[oldPath].(map[string]any)
	if pathItem == nil {
		return nil
	}
	if op := pathItem[strings.ToLower(oldMethod)]; op != nil {
		return op
	}
	if strings.EqualFold(oldMethod, "get") {
		return pathItem["post"]
	}
	if op := pathItem["post"]; op != nil {
		return op
	}
	return pathItem["get"]
}

func cloneJSONValue(v any) any {
	raw, err := json.Marshal(v)
	if err != nil {
		return v
	}
	var cloned any
	if err := json.Unmarshal(raw, &cloned); err != nil {
		return v
	}
	return cloned
}

func sofieSwaggerHTML(openapiPath string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Sofie App API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: %q,
        dom_id: '#swagger-ui'
      })
    }
  </script>
</body>
</html>`, openapiPath)
}
