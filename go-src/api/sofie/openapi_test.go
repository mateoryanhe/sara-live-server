package sofie

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRewriteOpenAPIRefsRenamesGameCode(t *testing.T) {
	schemas := map[string]any{
		"xr-game-server.dto.gameplatformdto.AppGameStartReq": map[string]any{
			"type": "object",
			"required": []any{
				"gameCode",
			},
			"properties": map[string]any{
				"gameCode": map[string]any{
					"type":        "string",
					"description": "游戏编码(对应第三方 gameId)",
				},
				"roomId": map[string]any{
					"type":        "integer",
					"description": "直播间ID:0=大厅玩游戏(清空关联)",
				},
			},
		},
		"xr-game-server.dto.gameplatformdto.AppGameStartRes": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"link": map[string]any{
					"type":        "string",
					"description": "游戏启动链接(有效期约2分钟,仅可使用一次)",
				},
			},
		},
	}
	op := map[string]any{
		"summary": "启动游戏",
		"requestBody": map[string]any{
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": map[string]any{
						"$ref": "#/components/schemas/xr-game-server.dto.gameplatformdto.AppGameStartReq",
					},
				},
			},
		},
		"responses": map[string]any{
			"200": map[string]any{
				"content": map[string]any{
					"application/json": map[string]any{
						"schema": map[string]any{
							"$ref": "#/components/schemas/xr-game-server.dto.gameplatformdto.AppGameStartRes",
						},
					},
				},
			},
		},
	}
	renames := map[string]string{"gameCode": "titleCode"}
	rewriteOpenAPIRefs(op, schemas, "/sofie/play/openTitle", renames, map[string]string{})

	reqRef := op["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)["$ref"].(string)
	if reqRef != "#/components/schemas/sofie_play_openTitle_AppTitleStartReq" {
		t.Fatalf("request $ref = %q", reqRef)
	}
	reqSch := schemas["sofie_play_openTitle_AppTitleStartReq"].(map[string]any)
	props := reqSch["properties"].(map[string]any)
	if _, ok := props["gameCode"]; ok {
		t.Fatal("gameCode still present")
	}
	title := props["titleCode"].(map[string]any)
	if title["description"] != "标题编码(对应第三方 titleId)" {
		t.Fatalf("titleCode description = %v", title["description"])
	}
	required := reqSch["required"].([]any)
	if len(required) != 1 || required[0] != "titleCode" {
		t.Fatalf("required = %#v", required)
	}
	room := props["roomId"].(map[string]any)
	if room["description"] != "直播间ID:0=大厅(清空关联)" {
		t.Fatalf("roomId description = %v", room["description"])
	}
	resRef := op["responses"].(map[string]any)["200"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)["$ref"].(string)
	if resRef != "#/components/schemas/sofie_play_openTitle_AppTitleStartRes" {
		t.Fatalf("response $ref = %q", resRef)
	}
	link := schemas["sofie_play_openTitle_AppTitleStartRes"].(map[string]any)["properties"].(map[string]any)["link"].(map[string]any)
	if link["description"] != "启动链接(有效期约2分钟,仅可使用一次)" {
		t.Fatalf("link description = %v", link["description"])
	}

	doc := map[string]any{
		"paths": map[string]any{
			"/sofie/play/openTitle": map[string]any{"post": op},
		},
		"components": map[string]any{"schemas": schemas},
	}
	pruneSofieOpenAPIComponents(doc)
	kept := doc["components"].(map[string]any)["schemas"].(map[string]any)
	if _, ok := kept["xr-game-server.dto.gameplatformdto.AppGameStartReq"]; ok {
		t.Fatal("old AppGameStartReq should be pruned")
	}
	if _, ok := kept["sofie_play_openTitle_AppTitleStartReq"]; !ok {
		t.Fatal("sofie req schema missing after prune")
	}
	raw, _ := json.Marshal(kept)
	text := string(raw)
	if strings.Contains(text, "gameCode") || strings.Contains(text, "游戏") {
		t.Fatalf("pruned schemas still contain game wording: %s", text)
	}
}
