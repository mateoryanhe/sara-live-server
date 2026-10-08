package entity

const (
	VoiceChatMicModeFree     uint8 = 1 // 自由上麦
	VoiceChatMicModeApply    uint8 = 2 // 申请上麦
	VoiceChatMicModeHostOnly uint8 = 3 // 房主单麦
)

// IsValidVoiceChatMicMode 语聊房上麦方式是否合法.
func IsValidVoiceChatMicMode(mode uint8) bool {
	return mode == VoiceChatMicModeFree ||
		mode == VoiceChatMicModeApply ||
		mode == VoiceChatMicModeHostOnly
}
