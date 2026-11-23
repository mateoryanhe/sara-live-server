package message

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// privateTextPayload 与 App 发送的文字私信 content JSON 一致（见测试服 user_messages 样例）。
type privateTextPayload struct {
	ID             string  `json:"id"`
	ConversationID string  `json:"conversationId"`
	SenderID       string  `json:"senderId"`
	SenderName     *string `json:"senderName"`
	CreatedAt      string  `json:"createdAt"`
	Type           string  `json:"type"`
	Text           string  `json:"text"`
	Attachment     any     `json:"attachment"`
}

// BuildPrivateTextMessageContent 将纯文字包装为 App 可解析的私信 content JSON。
func BuildPrivateTextMessageContent(senderId, receiverId uint64, text string) string {
	text = strings.TrimSpace(text)
	if text == "" || senderId == 0 || receiverId == 0 {
		return ""
	}
	now := time.Now()
	micro := now.UnixMicro()
	payload := privateTextPayload{
		ID:             fmt.Sprintf("p2p:%d_%d_0", receiverId, micro),
		ConversationID: fmt.Sprintf("p2p:%d", receiverId),
		SenderID:       strconv.FormatUint(senderId, 10),
		SenderName:     nil,
		CreatedAt:      now.Format("2006-01-02T15:04:05.000000"),
		Type:           "text",
		Text:           text,
		Attachment:     nil,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(b)
}

// IsPrivateTextMessageContent 判断 content 是否已是 App 文字私信 JSON。
func IsPrivateTextMessageContent(content string) bool {
	content = strings.TrimSpace(content)
	if content == "" || content[0] != '{' {
		return false
	}
	var probe struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(content), &probe); err != nil {
		return false
	}
	return probe.Type == "text" && probe.Text != ""
}