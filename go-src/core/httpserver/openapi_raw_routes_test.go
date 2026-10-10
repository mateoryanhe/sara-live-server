package httpserver

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMergeOpenApiRawHandlerRoutes_createLiveRoom(t *testing.T) {
	doc := map[string]any{"paths": map[string]any{}}
	if err := MergeOpenApiRawHandlerRoutes(doc); err != nil {
		t.Fatal(err)
	}
	paths, _ := doc["paths"].(map[string]any)
	item, ok := paths["/liveRoom/create"].(map[string]any)
	if !ok {
		t.Fatal("missing /liveRoom/create")
	}
	post, ok := item["post"].(map[string]any)
	if !ok {
		t.Fatal("missing post")
	}
	body, _ := post["requestBody"].(map[string]any)
	content, _ := body["content"].(map[string]any)
	multipart, _ := content["multipart/form-data"].(map[string]any)
	schema, _ := multipart["schema"].(map[string]any)
	ref, _ := schema["$ref"].(string)
	if ref == "" {
		t.Fatal("missing request schema ref")
	}
	schemas, _ := doc["components"].(map[string]any)["schemas"].(map[string]any)
	reqSchema, _ := schemas[openAPISchemaLiveRoomDTO+"CreateLiveRoomReq"].(map[string]any)
	props, _ := reqSchema["properties"].(map[string]any)
	if props["voiceChatMicMode"] == nil {
		t.Fatal("CreateLiveRoomReq missing voiceChatMicMode")
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "voiceChatMicMode") {
		t.Fatal("exported doc should contain voiceChatMicMode")
	}
}
