package gameevent

import "xr-game-server/core/event"

const (
	// SystemMessageEvent 发送系统消息事件
	SystemMessageEvent event.Type = "SystemMessageEvent"
)

// SystemMessageEventData 发送系统消息事件数据
type SystemMessageEventData struct {
	ReceiverId uint64 // 接收者用户ID
	Title      string // 标题
	Content    string // 内容
}

func NewSystemMessageEventData(receiverId uint64, title, content string) *SystemMessageEventData {
	return &SystemMessageEventData{
		ReceiverId: receiverId,
		Title:      title,
		Content:    content,
	}
}
