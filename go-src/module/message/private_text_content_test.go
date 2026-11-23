package message

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildPrivateTextMessageContent_shape(t *testing.T) {
	const sender, receiver uint64 = 6001000477, 6001000472
	raw := BuildPrivateTextMessageContent(sender, receiver, "Hello")
	if raw == "" {
		t.Fatal("empty content")
	}
	if !strings.HasPrefix(raw, "{") {
		t.Fatalf("not json: %s", raw)
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatal(err)
	}
	if m["type"] != "text" || m["text"] != "Hello" {
		t.Fatalf("unexpected payload: %v", m)
	}
	if m["conversationId"] != "p2p:6001000472" {
		t.Fatalf("conversationId=%v", m["conversationId"])
	}
	if m["senderId"] != "6001000477" {
		t.Fatalf("senderId=%v", m["senderId"])
	}
	id, _ := m["id"].(string)
	if !strings.HasPrefix(id, "p2p:6001000472_") || !strings.HasSuffix(id, "_0") {
		t.Fatalf("id=%s", id)
	}
}

func TestIsPrivateTextMessageContent(t *testing.T) {
	sample := `{"id":"p2p:1_2_0","conversationId":"p2p:1","senderId":"9","senderName":null,"createdAt":"2026-11-21T17:22:11.592553","type":"text","text":"Q","attachment":null}`
	if !IsPrivateTextMessageContent(sample) {
		t.Fatal("expected true")
	}
	if IsPrivateTextMessageContent("plain hello") {
		t.Fatal("expected false for plain text")
	}
}