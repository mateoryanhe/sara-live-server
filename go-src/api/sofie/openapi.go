package sofie

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

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
	oldPaths, _ := doc["paths"].(map[string]any)
	if oldPaths == nil {
		oldPaths = map[string]any{}
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
		}
		pathItem, _ := newPaths[route.NewPath].(map[string]any)
		if pathItem == nil {
			pathItem = map[string]any{}
		}
		pathItem["post"] = cloned
		newPaths[route.NewPath] = pathItem
	}
	doc["paths"] = newPaths
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
